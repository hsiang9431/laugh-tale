package console

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"laugh-tale/pkg/common/cli"
	"laugh-tale/pkg/common/http"
	"laugh-tale/pkg/kozuki/client"
	"laugh-tale/pkg/kozuki/types"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var verbose bool

func create(ctx *cli.Context) error {
	inputKey, kozukiClient, err := loadFileAndCreateClient(ctx)
	if err != nil {
		return err
	}
	retKey, err := kozukiClient.Create(inputKey)
	if err != nil {
		return errors.Wrap(err, "failed to create key")
	}
	printToConsole("http request finished")
	return outputToConsole(jsonFlag.GetBool(ctx), retKey)
}

func retrieve(ctx *cli.Context) error {
	inputKey, kozukiClient, err := loadFileAndCreateClient(ctx)
	if err != nil {
		return err
	}
	retKey, err := kozukiClient.Create(inputKey)
	if err != nil {
		return errors.Wrap(err, "failed to create key")
	}
	printToConsole("http request finished")
	return outputToConsole(jsonFlag.GetBool(ctx), retKey)
}

func update(ctx *cli.Context) error {
	inputKey, kozukiClient, err := loadFileAndCreateClient(ctx)
	if err != nil {
		return err
	}
	retKey, err := kozukiClient.Create(inputKey)
	if err != nil {
		return errors.Wrap(err, "failed to create key")
	}
	printToConsole("http request finished")
	return outputToConsole(jsonFlag.GetBool(ctx), retKey)
}

func delete(ctx *cli.Context) error {
	inputKey, kozukiClient, err := loadFileAndCreateClient(ctx)
	if err != nil {
		return err
	}
	err = kozukiClient.Delete(inputKey)
	if err != nil {
		return errors.Wrap(err, "failed to create key")
	}
	printToConsole("http request finished")
	return nil
}

func loadFileAndCreateClient(ctx *cli.Context) (types.Key, *client.CURLClient, error) {
	verbose = verboseFlag.GetBool(ctx)
	printToConsole("Set to verbose")
	inputFile, _ := filenameFlag.GetString(ctx)
	inputB, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return types.Key{}, nil, errors.Wrap(err, "failed to load input file")
	}
	printToConsole("Input file loaded: " + inputFile)
	inputKey := types.Key{}
	err = yaml.Unmarshal(inputB, &inputKey)
	if err != nil {
		return types.Key{}, nil, errors.Wrap(err, "failed to unmarshal input file")
	}
	printToConsole("Input file unmarshaled")
	httpsClient, err := http.HTTPClientTLSFromFile(CACertName, TLSCertName, TLSKeyName, true)
	if err != nil {
		return types.Key{}, nil, errors.Wrap(err, "failed to create https client")
	}
	printToConsole("HTTPS client created")
	kozukiClient := &client.CURLClient{
		URL:        ServerBaseURL,
		HTTPClient: httpsClient,
	}
	return inputKey, kozukiClient, nil
}

func outputToConsole(outputJson bool, k types.Key) error {
	if outputJson {
		retKeyB, err := json.Marshal(&k)
		if err != nil {
			return errors.Wrap(err, "failed to marshal returned key to json")
		}
		fmt.Println(string(retKeyB))
	} else {
		retKeyB, err := yaml.Marshal(&k)
		if err != nil {
			return errors.Wrap(err, "failed to marshal returned key to yaml")
		}
		fmt.Println(string(retKeyB))
	}
	return nil
}

func printToConsole(m string) {
	if verbose {
		fmt.Println("[KOZUKI_CLI] " + m)
	}
}
