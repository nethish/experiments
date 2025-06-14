package main

import (
	"context"
	"fmt"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// HttpbinApp represents our custom resource
type HttpbinApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HttpbinAppSpec   `json:"spec,omitempty"`
	Status            HttpbinAppStatus `json:"status,omitempty"`
}

type HttpbinAppSpec struct {
	Replicas int32 `json:"replicas"`
}

type HttpbinAppStatus struct {
	ReadyReplicas int32                 `json:"readyReplicas,omitempty"`
	Conditions    []HttpbinAppCondition `json:"conditions,omitempty"`
}

type HttpbinAppCondition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// HttpbinAppList contains a list of HttpbinApp
type HttpbinAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpbinApp `json:"items"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpbinApp) DeepCopyInto(out *HttpbinApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpbinApp.
func (in *HttpbinApp) DeepCopy() *HttpbinApp {
	if in == nil {
		return nil
	}
	out := new(HttpbinApp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HttpbinApp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpbinAppList) DeepCopyInto(out *HttpbinAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HttpbinApp, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpbinAppList.
func (in *HttpbinAppList) DeepCopy() *HttpbinAppList {
	if in == nil {
		return nil
	}
	out := new(HttpbinAppList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HttpbinAppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpbinAppSpec) DeepCopyInto(out *HttpbinAppSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpbinAppSpec.
func (in *HttpbinAppSpec) DeepCopy() *HttpbinAppSpec {
	if in == nil {
		return nil
	}
	out := new(HttpbinAppSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpbinAppStatus) DeepCopyInto(out *HttpbinAppStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]HttpbinAppCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpbinAppStatus.
func (in *HttpbinAppStatus) DeepCopy() *HttpbinAppStatus {
	if in == nil {
		return nil
	}
	out := new(HttpbinAppStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpbinAppCondition) DeepCopyInto(out *HttpbinAppCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpbinAppCondition.
func (in *HttpbinAppCondition) DeepCopy() *HttpbinAppCondition {
	if in == nil {
		return nil
	}
	out := new(HttpbinAppCondition)
	in.DeepCopyInto(out)
	return out
}

// HttpbinAppReconciler reconciles a HttpbinApp object
type HttpbinAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HttpbinAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling HttpbinApp", "name", req.Name, "namespace", req.Namespace)
	logger.Info("What is namespaced name", "namespacedname", req.NamespacedName)

	// Fetch the HttpbinApp instance
	var httpbinApp HttpbinApp
	if err := r.Get(ctx, req.NamespacedName, &httpbinApp); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("HttpbinApp resource not found. Ignoring since object must be deleted", "name", req.Name, "namespace", req.Namespace)
			// Check if there's an orphaned deployment that needs cleanup
			deployment := &appsv1.Deployment{}
			deploymentName := req.Name + "-httpbin"
			deploymentKey := types.NamespacedName{Name: deploymentName, Namespace: req.Namespace}

			if err := r.Get(ctx, deploymentKey, deployment); err == nil {
				// Deployment exists but HttpbinApp is gone - check if it's managed by us
				if deployment.Labels != nil && deployment.Labels["managed-by"] == "httpbin-controller" {
					logger.Info("Cleaning up orphaned deployment", "deployment", deploymentName, "namespace", req.Namespace)
					if err := r.Delete(ctx, deployment); err != nil {
						logger.Error(err, "Failed to delete orphaned deployment", "deployment", deploymentName)
						return ctrl.Result{RequeueAfter: time.Second * 30}, err
					}
				}
			}
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get HttpbinApp", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{RequeueAfter: time.Second * 30}, err
	}

	// Define the desired deployment
	deployment := r.deploymentForHttpbinApp(&httpbinApp)

	// Check if the deployment already exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{RequeueAfter: time.Second * 30}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Update the deployment if replicas are different
	if *found.Spec.Replicas != httpbinApp.Spec.Replicas {
		logger.Info("Updating Deployment replicas", "current", *found.Spec.Replicas, "desired", httpbinApp.Spec.Replicas)
		found.Spec.Replicas = &httpbinApp.Spec.Replicas
		err = r.Update(ctx, found)
		if err != nil {
			logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Second * 30}, nil
	}

	// Update status
	httpbinApp.Status.ReadyReplicas = found.Status.ReadyReplicas

	// Update conditions
	condition := HttpbinAppCondition{
		Type:               "Ready",
		Status:             "True",
		LastTransitionTime: metav1.Now(),
		Reason:             "DeploymentReady",
		Message:            fmt.Sprintf("Deployment has %d ready replicas", found.Status.ReadyReplicas),
	}

	httpbinApp.Status.Conditions = []HttpbinAppCondition{condition}

	err = r.Status().Update(ctx, &httpbinApp)
	if err != nil {
		logger.Error(err, "Failed to update HttpbinApp status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

// deploymentForHttpbinApp returns a httpbin Deployment object
func (r *HttpbinAppReconciler) deploymentForHttpbinApp(m *HttpbinApp) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "httpbin",
		"managed-by": "httpbin-controller",
		"instance":   m.Name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-httpbin",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &m.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "kennethreitz/httpbin:latest",
						Name:  "httpbin",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/get",
									Port: intstr.FromInt(80),
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       10,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/get",
									Port: intstr.FromInt(80),
								},
							},
							InitialDelaySeconds: 15,
							PeriodSeconds:       20,
						},
					}},
				},
			},
		},
	}

	// Set HttpbinApp instance as the owner and controller
	ctrl.SetControllerReference(m, deployment, r.Scheme)
	return deployment
}

// SetupWithManager sets up the controller with the Manager
func (r *HttpbinAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&HttpbinApp{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func main() {
	var (
		scheme   = runtime.NewScheme()
		setupLog = ctrl.Log.WithName("setup")
	)

	// Setup logging
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Add the core types to the scheme
	if err := corev1.AddToScheme(scheme); err != nil {
		setupLog.Error(err, "unable to add core types to scheme")
		os.Exit(1)
	}

	if err := appsv1.AddToScheme(scheme); err != nil {
		setupLog.Error(err, "unable to add apps types to scheme")
		os.Exit(1)
	}

	// Add our custom resource to the scheme with proper GroupVersion
	gv := schema.GroupVersion{Group: "example.com", Version: "v1"}
	scheme.AddKnownTypes(gv,
		&HttpbinApp{},
		&HttpbinAppList{},
	)
	metav1.AddToGroupVersion(scheme, gv)

	// Create manager with proper configuration
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Logger: ctrl.Log.WithName("manager"),
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup the controller
	if err = (&HttpbinAppReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HttpbinApp")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

