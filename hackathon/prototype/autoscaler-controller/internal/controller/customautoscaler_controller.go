/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	appsv1 "k8s.io/api/apps/v1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"autoscaler-controller/api/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CustomAutoscalerReconciler reconciles a CustomAutoscaler object
type CustomAutoscalerReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	MetricsClient metricsclient.Interface
	RedisClient   *redis.Client
	KubeClient    *kubernetes.Clientset
}

// +kubebuilder:rbac:groups=autoscaler-group.autoscaler,resources=customautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaler-group.autoscaler,resources=customautoscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autoscaler-group.autoscaler,resources=customautoscalers/finalizers,verbs=update

func (r *CustomAutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("customautoscaler", req.NamespacedName)

	// Получаем экземпляр CustomAutoscaler
	var customAutoscaler v1alpha1.CustomAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &customAutoscaler); err != nil {
		if errors.IsNotFound(err) {
			// Ресурс удален, удаляем информацию из Redis
			if err = r.deleteDeploymentInfoFromRedis(req.Name); err != nil {
				log.Error(err, "Ошибка при удалении информации о деплойментах из Redis")
			}
			return ctrl.Result{}, nil
		}
		log.Error(err, "Не удалось получить CustomAutoscaler")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Обновляем информацию о деплойментах в Redis
	if err := r.updateDeploymentInfoInRedis(&customAutoscaler); err != nil {
		log.Error(err, "Ошибка при обновлении информации о деплойментах в Redis")
	}

	// Логика автоскейлинга для каждого деплоймента
	for _, targetDeployment := range customAutoscaler.Spec.TargetDeployments {
		if err := r.handleDeploymentScaling(ctx, &customAutoscaler, targetDeployment); err != nil {
			log.Error(err, "Ошибка при обработке деплоймента", "deployment", targetDeployment.Name)
			// Продолжаем обработку остальных деплойментов
			continue
		}
	}

	// Реконсиляция через определенный интервал времени
	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// updateDeploymentInfoInRedis обновляет информацию о деплойментах в Redis
func (r *CustomAutoscalerReconciler) updateDeploymentInfoInRedis(customAutoscaler *v1alpha1.CustomAutoscaler) error {
	ctx := context.Background()

	// Преобразуем список деплойментов в JSON
	deploymentsData, err := json.Marshal(customAutoscaler.Spec.TargetDeployments)
	if err != nil {
		return err
	}

	// Записываем данные в Redis с использованием идентификатора CRD
	if err := r.RedisClient.HSet(ctx, "deployments_info", customAutoscaler.Name, deploymentsData).Err(); err != nil {
		return err
	}

	// Публикуем сообщение о том, что данные обновлены
	if err := r.RedisClient.Publish(ctx, "deployments_channel", "updated").Err(); err != nil {
		return err
	}

	return nil
}

// deleteDeploymentInfoFromRedis удаляет информацию из Redis при удалении ресурса
func (r *CustomAutoscalerReconciler) deleteDeploymentInfoFromRedis(crdName string) error {
	ctx := context.Background()

	// Удаляем данные из Redis с использованием идентификатора CRD
	if err := r.RedisClient.HDel(ctx, "deployments_info", crdName).Err(); err != nil {
		return err
	}

	// Публикуем сообщение о том, что данные удалены
	if err := r.RedisClient.Publish(ctx, "deployments_channel", "deleted").Err(); err != nil {
		return err
	}

	return nil
}

func (r *CustomAutoscalerReconciler) handleDeploymentScaling(ctx context.Context,
	customAutoscaler *v1alpha1.CustomAutoscaler,
	targetDeployment v1alpha1.TargetDeployment) error {

	log := r.Log.WithValues("deployment", targetDeployment.Name)

	// Получаем деплоймент
	var deployment appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKey{Name: targetDeployment.Name, Namespace: customAutoscaler.Namespace}, &deployment); err != nil {
		if errors.IsNotFound(err) {
			log.Info(fmt.Sprintf("Деплоймент %s не найден", targetDeployment.Name))
			return nil
		}
		return err
	}

	// Получаем метки селектора деплоймента
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return err
	}

	// Получаем метрики подов
	podMetricsList, err := r.MetricsClient.MetricsV1beta1().PodMetricses(customAutoscaler.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		log.Error(err, "Ошибка при получении метрик подов")
		return err
	}

	var totalCPU int64
	var totalMemory int64
	podCount := int64(len(podMetricsList.Items))

	for _, podMetrics := range podMetricsList.Items {
		for _, container := range podMetrics.Containers {
			totalCPU += container.Usage.Cpu().MilliValue()
			totalMemory += container.Usage.Memory().Value() / (1024 * 1024) // В MiB
		}
	}

	if podCount == 0 {
		log.Info("Нет запущенных подов для деплоймента")
		return nil
	}

	avgCPU := totalCPU / podCount
	avgMemory := totalMemory / podCount

	// Получаем метрики из Redis
	redisCtx := context.Background()
	metrics, err := r.RedisClient.HGetAll(redisCtx, deployment.Name).Result()
	if err != nil {
		log.Error(err, "Ошибка при получении метрик из Redis")
		return err
	}

	// Парсим время последнего запроса
	lastRequestTimeUnix, err := strconv.ParseInt(metrics["lastRequestTime"], 10, 64)
	if err != nil {
		log.Error(err, "Ошибка при парсинге времени последнего запроса")
		lastRequestTimeUnix = time.Now().Unix()
	}
	lastRequestTime := time.Unix(lastRequestTimeUnix, 0)

	// Логика масштабирования вверх
	shouldScaleUp := false
	if avgCPU > customAutoscaler.Spec.CpuThreshold {
		shouldScaleUp = true
		log.Info("Среднее использование CPU превышает порог", "avgCPU", avgCPU)
	}
	if avgMemory > customAutoscaler.Spec.MemoryThreshold {
		shouldScaleUp = true
		log.Info("Среднее использование памяти превышает порог", "avgMemory", avgMemory)
	}

	if shouldScaleUp {
		newReplicas := *deployment.Spec.Replicas + 1
		log.Info("Увеличиваем количество реплик", "новое количество", newReplicas)
		return r.scaleDeployment(ctx, &deployment, newReplicas)
	}

	// Логика масштабирования вниз
	if avgCPU < customAutoscaler.Spec.CpuThreshold/2 &&
		avgMemory < customAutoscaler.Spec.MemoryThreshold/2 &&
		*deployment.Spec.Replicas > 1 {
		newReplicas := *deployment.Spec.Replicas - 1
		log.Info("Уменьшаем количество реплик", "новое количество", newReplicas)
		return r.scaleDeployment(ctx, &deployment, newReplicas)
	}

	// Проверка времени последнего запроса
	idleDuration := time.Since(lastRequestTime)
	idleTime := time.Duration(customAutoscaler.Spec.IdleTime) * time.Second

	if idleDuration >= idleTime && *deployment.Spec.Replicas > 0 {
		// Масштабируем до 0 реплик
		log.Info("Приложение неактивно, масштабируем до 0 реплик")
		return r.scaleDeployment(ctx, &deployment, 0)
	}

	return nil
}

func (r *CustomAutoscalerReconciler) scaleDeployment(ctx context.Context, deployment *appsv1.Deployment, replicas int32) error {
	deployment.Spec.Replicas = &replicas
	if err := r.Update(ctx, deployment); err != nil {
		r.Log.Error(err, "Ошибка при масштабировании деплоймента", "deployment", deployment.Name, "replicas", replicas)
		return err
	}
	r.Log.Info("Деплоймент масштабирован", "deployment", deployment.Name, "replicas", replicas)
	return nil
}

func (r *CustomAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Инициализация MetricsClient
	config := ctrl.GetConfigOrDie()
	metricsClient, err := metricsclient.NewForConfig(config)
	if err != nil {
		return err
	}
	r.MetricsClient = metricsClient

	// Инициализация RedisClient
	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := 0
	if redisDBStr != "" {
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			r.Log.Error(err, "Неверное значение REDIS_DB, используется значение по умолчанию 0")
			redisDB = 0
		}
	}

	r.RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	if r.RedisClient == nil {
		r.Log.Error(nil, "Не удалось создать клиент Redis")
	}

	if err = r.RedisClient.Ping(context.Background()).Err(); err != nil {
		r.Log.Error(err, "Не удалось подключиться к Redis")
	} else {
		r.Log.Info("Подключение к Redis установлено")
	}

	// Инициализация KubeClient
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	r.KubeClient = kubeClient

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.CustomAutoscaler{}).
		Complete(r)
}
