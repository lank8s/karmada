/*
Copyright 2023 The Karmada Authors.

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

package crossclusterstatefulset

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1alpha1 "github.com/karmada-io/karmada/pkg/apis/apps/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/resourceinterpreter"
	"github.com/karmada-io/karmada/pkg/sharedcli/ratelimiterflag"
)

// RBApplicationFailoverControllerName is the controller name that will be used when reporting events.
const CrossClusterStatefulSetControllerName = "cross-cluster-statefulset-controller"

// CrossClusterStatefulSetController is to sync ResourceBinding's application failover behavior.
type CrossClusterStatefulSetController struct {
	client.Client
	EventRecorder      record.EventRecorder
	RateLimiterOptions ratelimiterflag.Options

	ResourceInterpreter resourceinterpreter.ResourceInterpreter
}

func (c *CrossClusterStatefulSetController) Reconcile(ctx context.Context, req controllerruntime.Request) (controllerruntime.Result, error) {
	klog.V(4).Infof("Reconciling CrossClusterStatefulSet %s.", req.NamespacedName.String())

	binding := &workv1alpha2.ResourceBinding{}
	if err := c.Client.Get(ctx, req.NamespacedName, binding); err != nil {
		if apierrors.IsNotFound(err) {
			return controllerruntime.Result{}, nil
		}
		return controllerruntime.Result{Requeue: true}, err
	}
	return controllerruntime.Result{}, nil
}

// SetupWithManager creates a controller and register to controller manager.
func (c *CrossClusterStatefulSetController) SetupWithManager(mgr controllerruntime.Manager) error {
	resourcePredicateFn := predicate.Funcs{
		CreateFunc: func(createEvent event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return true
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool { return false },
	}

	return controllerruntime.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.CrossClusterStatefulSet{}, builder.WithPredicates(resourcePredicateFn)).
		WithOptions(controller.Options{RateLimiter: ratelimiterflag.DefaultControllerRateLimiter(c.RateLimiterOptions)}).
		Complete(c)
}
