package kubecli

import (
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewClientSetFromDefaultConfig() (*kubernetes.Clientset, error) {
	return NewClientSetFromKubeconfig(path.Join(homedir.HomeDir(), ".kube", "config"))
}

func NewClientSetFromKubeconfig(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}