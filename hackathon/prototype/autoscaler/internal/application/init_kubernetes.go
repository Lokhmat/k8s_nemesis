package application

import (
	"context"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (app *Application) initKubernetesClient(_ context.Context) error {
	var kubeConfig *rest.Config
	var err error
	if app.config.KubeConfig != "" {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", app.config.KubeConfig)
	} else {
		kubeConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return err
	}

	clientSet, err := kube.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	app.clientSet = clientSet
	return nil
}
