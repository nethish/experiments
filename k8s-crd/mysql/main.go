package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	// ✅ typed clientset
	//
	"k8s.io/client-go/kubernetes" // ✅ typed clientset
)

type MySQLInstanceSpec struct {
	Version  string `json:"version"`
	Storage  string `json:"storage"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MySQLInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLInstanceSpec `json:"spec"`
}

func getKubeConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	// config, err := rest.InClusterConfig() // or rest.InClusterConfig()

	// if err != nil {
	// 	panic(err)
	// }
	config := getKubeConfig()

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "mysql.my.domain",
		Version:  "v1alpha1",
		Resource: "mysqlinstances",
	}

	resClient := dynClient.Resource(gvr).Namespace("default")

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return resClient.List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return resClient.Watch(context.TODO(), opts)
		},
	}

	_, controller := cache.NewInformer(
		lw,
		&unstructured.Unstructured{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)
				name := u.GetName()
				ns := u.GetNamespace()
				spec := u.Object["spec"].(map[string]interface{})
				version := spec["version"].(string)

				// Create Deployment for MySQL
				dep := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mysql-" + name,
						Namespace: ns,
						Labels: map[string]string{
							"app": "mysql",
							"cr":  name,
						},
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: int32Ptr(1),
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "mysql", "cr": name},
						},
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{"app": "mysql", "cr": name},
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "mysql",
										Image: "mysql:" + version,
										Ports: []corev1.ContainerPort{{ContainerPort: 3306}},
										Env: []corev1.EnvVar{
											{Name: "MYSQL_ROOT_PASSWORD", Value: "password"},
										},
									},
								},
							},
						},
					},
				}

				// Create it using typed client
				clientset, err := kubernetes.NewForConfig(config)
				if err != nil {
					panic(err)
				}
				_, err = clientset.AppsV1().Deployments(ns).Create(context.TODO(), dep, metav1.CreateOptions{})
				if err != nil {
					fmt.Println("Error creating deployment:", err)
				} else {
					fmt.Println("Deployment created for MySQLInstance:", name)
				}
			},
			UpdateFunc: func(_, newObj interface{}) {
				raw := newObj.(*unstructured.Unstructured).Object
				time.Sleep(5 * time.Second)
				handleMySQL(raw)
			},
			DeleteFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)
				name := u.GetName()
				ns := u.GetNamespace()

				clientset, err := kubernetes.NewForConfig(config)
				if err != nil {
					panic(err)
				}
				err = clientset.AppsV1().Deployments(ns).Delete(context.TODO(), "mysql-"+name, metav1.DeleteOptions{})
				if err != nil {
					fmt.Println("Error deleting deployment:", err)
				} else {
					fmt.Println("Deleted deployment for:", name)
				}
			},
		},
	)

	stop := make(chan struct{})
	fmt.Println("Starting controller...")
	controller.Run(stop)
}

func handleMySQL(obj map[string]interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error marshaling obj:", err)
		return
	}

	var mysql MySQLInstance
	if err := json.Unmarshal(data, &mysql); err != nil {
		fmt.Println("error unmarshaling:", err)
		return
	}

	fmt.Printf("Reconciling MySQLInstance %s: version=%s storage=%s\n",
		mysql.Name, mysql.Spec.Version, mysql.Spec.Storage)
}

func int32Ptr(i int32) *int32 { return &i }
