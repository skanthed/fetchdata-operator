/*
Copyright 2022.

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

package controllers

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	//"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mydomainv1alpha1 "operators/api/v1alpha1"
)

// FetchdataReconciler reconciles a Fetchdata object
type FetchdataReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=fetchdata,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=fetchdata/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=fetchdata/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Fetchdata object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *FetchdataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Enter reconcile", "req", req)
	Fetchdata := &mydomainv1alpha1.Fetchdata{}
	//err := r.Get(ctx, types.NamespacedName, Fetchdata)

	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, Fetchdata)

	//l.Info("Inside reconcile", "error", err)

	//l.Info("Specification from crds", "Specification", Fetchdata.Spec)

	createjob, err := r.createCronJob(Fetchdata)

	if err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Inside job", "error", createjob)

	return ctrl.Result{}, nil
}

func (r *FetchdataReconciler) createCronJob(m *mydomainv1alpha1.Fetchdata) (string, error) {
	if _, err := FetchCronJob(m.Name, m.Namespace, r.Client); err != nil {
		//fmt.Println("Inside fetch cronjob ", err)
		if err := r.Client.Create(context.TODO(), NewCronJob(m, r.Scheme)); err != nil {
			//fmt.Println("Create method ", err)
			return m.Spec.Schedule, nil
		}
	}

	return m.Spec.Schedule, nil
}

func FetchCronJob(name, namespace string, client client.Client) (*batchv1.CronJob, error) {
	cronJob := &batchv1.CronJob{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	//fmt.Println("Cronjob value ", cronJob, err)
	return cronJob, err
}

func NewCronJob(m *mydomainv1alpha1.Fetchdata, scheme *runtime.Scheme) *batchv1.CronJob {
	//fmt.Println("Name and Namespace", m.Namespace, m.Name)
	cronjob := &batchv1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind: "Cronjob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: "fetchdata-operator-system",
		},
		Spec: batchv1.CronJobSpec{
			ConcurrencyPolicy: "Forbid",
			Schedule:          m.Spec.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								corev1.Volume{
									Name: "koku-metrics-operator-data",
									VolumeSource: corev1.VolumeSource{
										PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
											ClaimName: "koku-metrics-operator-data",
										},
									},
								},
							},
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "s3sync",
									Image: "quay.io/operate-first/curator-s3-sync:latest",
									VolumeMounts: []corev1.VolumeMount{
										corev1.VolumeMount{
											Name:      "koku-metrics-operator-data",
											MountPath: "/tmp/koku-metrics-operator-data",
										},
									},
								},
							},
							RestartPolicy: "Never",
						},
					},
				},
			},
		},
	}

	return cronjob

}

// SetupWithManager sets up the controller with the Manager.
func (r *FetchdataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1alpha1.Fetchdata{}).
		Complete(r)
}
