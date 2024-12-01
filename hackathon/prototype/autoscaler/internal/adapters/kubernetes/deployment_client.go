package kubernetes

import (
	"autoscaler/internal/domain/entities"
	"autoscaler/internal/port/services"
	"context"
	"strings"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type deploymentClient struct {
	clientSet     *kube.Clientset
	namespace     string
	podReadyCache map[string]bool
	podCacheMutex sync.RWMutex
}

func NewDeploymentClient(clientSet *kube.Clientset, namespace string) services.DeploymentClient {
	dc := &deploymentClient{
		clientSet:     clientSet,
		namespace:     namespace,
		podReadyCache: make(map[string]bool),
	}
	dc.startPodInformer()
	return dc
}

func (dc *deploymentClient) ScaleDeployment(ctx context.Context, deployment *entities.Deployment, replicas int32) error {
	deploymentsClient := dc.clientSet.AppsV1().Deployments(dc.namespace)
	scale, err := deploymentsClient.GetScale(ctx, deployment.Name, v1.GetOptions{})
	if err != nil {
		return err
	}

	scale.Spec.Replicas = replicas

	_, err = deploymentsClient.UpdateScale(ctx, deployment.Name, scale, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (dc *deploymentClient) ArePodsAvailable(ctx context.Context, deployment *entities.Deployment) (bool, error) {
	dc.podCacheMutex.RLock()
	defer dc.podCacheMutex.RUnlock()
	for podName, ready := range dc.podReadyCache {
		if strings.Contains(podName, deployment.Name) && ready {
			return true, nil
		}
	}
	return false, nil
}

func (dc *deploymentClient) startPodInformer() {
	factory := informers.NewSharedInformerFactoryWithOptions(
		dc.clientSet,
		0,
		informers.WithNamespace(dc.namespace),
	)

	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dc.onPodAdd,
		UpdateFunc: dc.onPodUpdate,
		DeleteFunc: dc.onPodDelete,
	})

	stopCh := make(chan struct{})
	go podInformer.Run(stopCh)
}

func (dc *deploymentClient) onPodAdd(obj interface{}) {
	pod := obj.(*corev1.Pod)
	dc.updatePodCache(pod)
}

func (dc *deploymentClient) onPodUpdate(oldObj, newObj interface{}) {
	pod := newObj.(*corev1.Pod)
	dc.updatePodCache(pod)
}

func (dc *deploymentClient) onPodDelete(obj interface{}) {
	pod := obj.(*corev1.Pod)
	dc.podCacheMutex.Lock()
	delete(dc.podReadyCache, pod.Name)
	dc.podCacheMutex.Unlock()
}

func (dc *deploymentClient) updatePodCache(pod *corev1.Pod) {
	dc.podCacheMutex.Lock()
	dc.podReadyCache[pod.Name] = pod.Status.Phase == corev1.PodRunning
	dc.podCacheMutex.Unlock()
}
