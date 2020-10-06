package poneglyph

const shellName = "/bin/sh"

const (
	encPrefix     = "enc_"
	implantBinary = "poneglyph"
	implantScript = "poneglyph.sh"
)

var oharaHostEnv = []string{
	"OHARA_SERVICE_HOST",
	"KEY_RETRIEVING_SERVICE",
	"KEY_RETRIEVING_SERVICE_HOST",
}

var nsNameEnv = []string{
	"NAMESPACE",
	"POD_NAMESPACE",
	"K8s_POD_NAMESPACE",
}

var srvName = []string{
	"ohara",
}
