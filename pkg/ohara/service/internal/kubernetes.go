package internal

type ClusterProvider interface {
	SetTargetContainer(string) error
	GetImageID() (string, error)
	GetMountedVolume() (string, error)
}

type KubernetesProvider struct {
}

func NewKubernetesProvider(k8sAPISrv string) (*KubernetesProvider, error) {

	return nil, nil
}

func (bp *KubernetesProvider) SetTargetContainer(string) error {

	return nil
}

func (bp *KubernetesProvider) GetImageID() (string, error) {

	return "", nil
}

func (bp *KubernetesProvider) GetMountedVolume() (string, error) {

	return "", nil
}
