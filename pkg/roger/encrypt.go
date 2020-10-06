package roger

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"laugh-tale/pkg/common/cli"
	"laugh-tale/pkg/common/http"
	"laugh-tale/pkg/common/task"
	kClient "laugh-tale/pkg/kozuki/client"
	kTypes "laugh-tale/pkg/kozuki/types"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

func encrypt(ctx *cli.Context) error {
	verbose = verboseFlag.GetBool(ctx)
	logInfo("Set to verbose")
	// load all args
	keyServ, _ := keyServerFlag.GetString(ctx)
	img, err := imageFlag.GetString(ctx)
	if err != nil {
		logInfo("Using default base image " + defaultBaseImage)
		img = defaultBaseImage
	}
	vols, err := volumeFlag.GetStringSlice(ctx)
	if err != nil {
		logWarning("No volume is declared")
	}
	logInfo("All CLI flags loaded")
	var bobPwd, payloadPwd string
	if len(keyServ) > 0 {
		bobPwd, payloadPwd, err = getKey(keyServ)
		if err != nil {
			return err
		}
	}
	logInfo("All password generated")
	// cleanup output directory
	dir, err := ioutil.ReadDir(OutputDirPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{OutputDirPath, d.Name()}...))
	}
	// tar and encrypt payload
	taet := &tarAndEncTask{PayloadPassword: payloadPwd}
	// copy other files
	cpt := &cpFilesTask{}
	// run tasks asynchronously
	tr := task.AsyncTaskRunner{
		Tasks: map[string]task.Task{
			"tar-and-enc": taet,
			"cp-files":    cpt,
		},
	}
	// run all tasks
	tr.InitAll()
	tr.RunAll()
	errs := tr.WaitAll()
	if len(errs) != 0 {
		ret := errors.New("")
		for _, err := range errs {
			ret = errors.Wrap(ret, err.Error())
		}
		return errors.Wrap(ret, "One or more error(s) encountered preparing files")
	}
	// generate dockerfile
	dockerfileOut := filepath.Join(OutputDirPath, "Dockerfile")
	// bobPwd := crypto.RandB64String(bobPWDLen)
	f, err := os.OpenFile(dockerfileOut, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to create Dockerfile")
	}
	if _, err := fmt.Fprintln(f, dockerfileTemplate,
		img, bobPwd, img, prepVolumeStr(vols)); err != nil {
		return errors.Wrap(err, "Failed to write Dockerfile")
	}
	logInfo("Dockerfile generated")
	fmt.Fprintln(ConsoleWriter, "Files for building encrypted image are ready")
	fmt.Fprintln(ConsoleWriter, "Run the following command in output directory to build your image")
	fmt.Fprintln(ConsoleWriter, "")
	fmt.Fprintln(ConsoleWriter, "cat Dockerfile | docker build -t <your image tag> ./")
	fmt.Fprintln(ConsoleWriter, "")
	fmt.Fprintln(ConsoleWriter, "After build and push the image to docker registry,")
	fmt.Fprintln(ConsoleWriter, "pull the image with command:")
	fmt.Fprintln(ConsoleWriter, "docker pull <your image tag>")
	fmt.Fprintln(ConsoleWriter, "")
	fmt.Fprintln(ConsoleWriter, "And get the image digest:")
	fmt.Fprintln(ConsoleWriter, "docker images --digests")
	return nil
}

func getKey(keyServ string) (string, string, error) {
	tlsCert := filepath.Join(SecretDirPath, TLSCertName)
	tlsKey := filepath.Join(SecretDirPath, TLSKeyName)
	tlsKeyPassB, _ := base64.RawURLEncoding.DecodeString(TLSKeyPass)
	httpClient, err := http.HTTPClientTLSFromFilePassphrase(CACertPath, tlsCert, tlsKey, tlsKeyPassB, false)
	if err != nil {
		logError("Failed to load CA keys")
		return "", "", errors.Wrap(err, "Failed to load CA keys")
	}
	ksClient := kClient.Client{
		URL:        keyServ,
		HTTPClient: httpClient,
	}
	var key kTypes.Key
	if key, err = ksClient.CreateKey(); err != nil {
		logError("Failed to create key from key manager")
		return "", "", errors.Wrap(err, "Failed to create key from key manager")
	}
	return key.ImplantKey, key.DecryptKey, nil
}

// volume must not be mounted right under root directory
var goodVolRegex, _ = regexp.Compile(`^\/(?!\.\.)[0-9a-zA-Z_\-.]+\/(?!\.\.)[0-9a-zA-Z_\-.]+`)

func prepVolumeStr(vols []string) string {
	ret := ""
	for _, v := range vols {
		if !goodVolRegex.MatchString(v) {
			logWarning("Insecured volume \"" + v + "\" detected, edit the Dockerfile to fix it")
			logWarning("Please do not mount volume under root directory")
		}
		ret = fmt.Sprintf(",\"%s\"%s", v, ret)
	}
	if len(ret) > 0 {
		ret = ret[1:]
		ret = "VOLUME [" + ret + "]"
	}
	return ret
}
