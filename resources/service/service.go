package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	devopsv1 "my.domain/etcd/api/v1"
)

func New(etcd *devopsv1.Etcd) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      etcd.Name,
			Namespace: etcd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(etcd, schema.GroupVersionKind{
					Group:   devopsv1.GroupVersion.Group,
					Version: devopsv1.GroupVersion.Version,
					Kind:    "Etcd",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: 2379,
					Name: "etcd-client",
				},
				corev1.ServicePort{
					Port: 2380,
					Name: "etcd-server",
				},
			},
			ClusterIP: corev1.ClusterIPNone,
			Selector: map[string]string{
				"app.example.com/v1alpha1": etcd.Name,
			},
		},
	}
}
