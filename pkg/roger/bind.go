package roger

import (
	"encoding/base64"
	"laugh-tale/pkg/common/cli"
	"laugh-tale/pkg/common/http"
	kClient "laugh-tale/pkg/kozuki/client"
	"path/filepath"

	"github.com/pkg/errors"
)

func bind(ctx *cli.Context) error {
	verbose = verboseFlag.GetBool(ctx)
	logInfo("Set to verbose")
	// load all args
	keyID, _ := keyIDFlag.GetString(ctx)
	imgID, _ := imageIDFlag.GetString(ctx)
	keyServ, _ := keyServerFlag.GetString(ctx)
	logInfo("All CLI flags loaded")

	tlsCert := filepath.Join(SecretDirPath, TLSCertName)
	tlsKey := filepath.Join(SecretDirPath, TLSKeyName)
	tlsKeyPassB, _ := base64.RawURLEncoding.DecodeString(TLSKeyPass)
	httpClient, err := http.HTTPClientTLSFromFilePassphrase(CACertPath, tlsCert, tlsKey, tlsKeyPassB, false)
	if err != nil {
		logError("Failed to load CA keys")
		return errors.Wrap(err, "Failed to load CA keys")
	}
	ksClient := kClient.Client{
		URL:        keyServ,
		HTTPClient: httpClient,
	}
	if _, err = ksClient.BindKey(keyID, imgID); err != nil {
		logError("Failed to bind key to image id")
		return errors.Wrap(err, "Failed to bind key to image id")
	}
	return nil
}
