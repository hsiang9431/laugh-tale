package roger

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// the complexity of bob's password is
// pwd_len * 3 / 4 bytes
// pwd_len = 60 provides the complexity of 360 bits
const bobPWDLen = 60

const errStrBadFlag = "Failed to load required flag"

const (
	encPrefix     = "enc_"
	implantBinary = "poneglyph"
	implantScript = "poneglyph.sh"
)

var signTemplate = &x509.Certificate{
	SerialNumber: big.NewInt(1),
	Issuer: pkix.Name{
		Organization: []string{"secret-container"},
	},
	Subject: pkix.Name{
		Organization: []string{"default"},
	},
	NotBefore:             time.Now(),
	NotAfter:              time.Now().Add(time.Hour * 24 * 180),
	KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth | x509.ExtKeyUsageClientAuth},
	BasicConstraintsValid: true,
}

const defaultBaseImage = "ubuntu:latest"

// CentOS and Ubuntu
// strings:
// 1, 3 base image tag
// 2. password in plain text
// 4. comma separated list of volume names with double quot: "/mnt/vol-1","/mnt/vol-2"
const dockerfileTemplate = `
FROM %s
RUN apt-get update && apt-get install -y sudo openssl
RUN useradd -p $(openssl passwd -6 %s) -g sudo poneglyph

FROM %s
RUN apt-get update && apt-get install -y sudo
COPY --from=0 ["/etc/passwd", "/etc/shadow", "/etc/group", "/etc/gshadow", "/etc/"]
RUN (chmod o-rwx /* /bin/* || :) && \
    chmod o+x /bin /bin/sh /lib /lib64 /usr && \
    echo 'set +o history' >> /etc/profile
%s
RUN mkdir /secret-container
COPY . /secret-container/
RUN chmod -R o-w /secret-container
WORKDIR /secret-container
USER poneglyph
ENTRYPOINT ["/secret-container/poneglyph", "start"]
`
