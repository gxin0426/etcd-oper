

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type EtcdRestorePhase string

var (
	EtcdRestoreRunning   EtcdRestorePhase = "Running"
	EtcdRestoreComplated EtcdRestorePhase = "Complated"
	EtcdRestoreFailed    EtcdRestorePhase = "Failed"
)


// EtcdRestoreSpec defines the desired state of EtcdRestore
type EtcdRestoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ClusterReference string `json:"clusterReference"`
	DataURl          string `json:"dataURL"`
}

// EtcdRestoreStatus defines the observed state of EtcdRestore
type EtcdRestoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []EtcdRestoreCondition `json:"conditions"`
	Phase      EtcdRestorePhase       `json:"phase"`
}

type EtcdRestoreCondition struct {
	Ready                 bool        `json:"ready"`
	Reason                string      `json:"reason,omitempty"`
	Message               string      `json:"message,omitempty"`
	LastedTranslationTime metav1.Time `json:"lastedTranslationTime"`
}

// +kubebuilder:object:root=true

// EtcdRestore is the Schema for the etcdrestores API
type EtcdRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdRestoreSpec   `json:"spec,omitempty"`
	Status EtcdRestoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdRestoreList contains a list of EtcdRestore
type EtcdRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdRestore{}, &EtcdRestoreList{})
}
