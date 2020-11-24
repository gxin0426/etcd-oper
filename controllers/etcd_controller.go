/*


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
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"my.domain/etcd/resources/service"
	"my.domain/etcd/resources/statefulset"
	"reflect"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	appsv1 "k8s.io/api/apps/v1"
	devopsv1 "my.domain/etcd/api/v1"
)

// EtcdReconciler reconciles a Etcd object
type EtcdReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=devops.my.domain,resources=etcds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.my.domain,resources=etcds/status,verbs=get;update;patch

func (r *EtcdReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("etcd", req.NamespacedName)

	etcd := &devopsv1.Etcd{}
	err := r.Client.Get(ctx, req.NamespacedName, etcd)
	if err != nil {
		if errors.IsNotFound(err) {
			return  ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if etcd.DeletionTimestamp != nil {
		return  ctrl.Result{}, nil
	}

	found := &appsv1.StatefulSet{}

	err = r.Client.Get(ctx, types.NamespacedName{Namespace: etcd.Namespace, Name: etcd.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		headlessService := service.New(etcd)
		err = r.Client.Create(ctx, headlessService)
		if err != nil {
			return ctrl.Result{}, err
		}
		ss := statefulset.New(etcd)
		err = r.Client.Create(ctx, ss)
		if err != nil {
			go r.Client.Delete(ctx, headlessService)
			return ctrl.Result{}, err
		}

		etcd.Annotations = map[string]string{"etcd.app.example.com/spec": toString(etcd)}
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Client.Update(ctx, etcd)
		})
		if retryErr != nil {
			fmt.Println(retryErr.Error())
		}
		return ctrl.Result{}, nil
	}else if err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(etcd.Spec, toSpec(etcd.Annotations["etcd.app.example.com/sepc"])) {
		ss := statefulset.New(etcd)
		found.Spec = ss.Spec
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Client.Update(ctx, ss)
		})
		if retryErr != nil {
			return ctrl.Result{}, nil
		}
	}
	return ctrl.Result{}, nil
}

func (r *EtcdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsv1.Etcd{}).
		Complete(r)
}


func toString(etcd *devopsv1.Etcd) string {
	data, _ := json.Marshal(etcd.Spec)
	return  string(data)
}

func toSpec(data string) devopsv1.EtcdSpec {
	spec := devopsv1.EtcdSpec{}
	json.Unmarshal([]byte(data), &spec)
	return  spec
}
