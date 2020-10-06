package poneglyph

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"laugh-tale/pkg/common/cli"
	"laugh-tale/pkg/common/crypto"
	"laugh-tale/pkg/common/http"
	oClient "laugh-tale/pkg/ohara/client"

	"github.com/pkg/errors"
)

func start(ctx *cli.Context) error {
	if os.Getpid() != 1 {
		return errors.New("This command must run with pid 1")
	}
	template := crypto.GetTemplateFromCSR(nil, nil)
	tlsCert, tlsKey, err := crypto.GenSelfSignCert(template, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to generate TLS key and certificate")
	}
	httpClient, err := http.HTTPClientTLSFromPem(ServerCACert, tlsCert, tlsKey, true)
	if err != nil {
		return errors.Wrap(err, "Cannot load CA certificate of key retriever service")
	}
	krClient := oClient.Client{
		HTTPClient: httpClient,
	}
	if len(ServerURL) > 0 {
		krClient.URL = ServerURL
	} else if err = serviceDiscovery(&krClient); err != nil {
		return err
	}
	key, err := krClient.GetKey()
	if err != nil {
		return errors.Wrap(err, "Failed to get key from key retriever service")
	}
	// execute privilege escalation script
	sh := filepath.Join(WorkDir, implantScript)
	cmd := exec.Command(shellName, sh,
		key.ImplantKey, key.DecryptKey)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Phase 2 returned with error:\n"+out.String())
	}
	fmt.Println("Phase 2 returned successfully")
	fmt.Println(out.String())
	return nil
}

func serviceDiscovery(c *oClient.Client) error {
	for _, e := range oharaHostEnv {
		host := os.Getenv(e)
		if len(host) > 0 {
			c.URL = host
			if err := c.Discover(); err == nil {
				return nil
			}
		}
	}
	for _, e := range nsNameEnv {
		ns := os.Getenv(e)
		if len(ns) > 0 {
			for _, s := range srvName {
				dnsName := s + "." + ns
				c.URL = dnsName
				if err := c.Discover(); err == nil {
					return nil
				}
			}
		}
	}
	return errors.New("Cannot discover key retriever service")
}
