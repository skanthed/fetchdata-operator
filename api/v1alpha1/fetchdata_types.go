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

package v1alpha1

import (
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FetchdataSpec defines the desired state of Fetchdata
type FetchdataSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Namespace in which cron job is created
	CronjobNamespace string `json:"cronjobNamespace,omitempty"`

	//Schedule period for the CronJob
	Schedule string `json:"schedule,omitempty"`

	//Koku metrics pvc zipped files storage path
	BackupSrc string `json:"backupSrc,omitempty"`

	//Koku-metrics-pvc path to unzip files
	UnzipDir string `json:"unzipDir,omitempty"`
}

// FetchdataStatus defines the observed state of Fetchdata
type FetchdataStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Name of the CronJob object created and managed by it
	CronJobName string `json:"cronJobName"`

	//CronJobStatus represents the current state of a cronjob
	CronJobStatus v1beta1.CronJobStatus `json:"cronJobStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Fetchdata is the Schema for the fetchdata API
type Fetchdata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FetchdataSpec   `json:"spec,omitempty"`
	Status FetchdataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FetchdataList contains a list of Fetchdata
type FetchdataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fetchdata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fetchdata{}, &FetchdataList{})
}
