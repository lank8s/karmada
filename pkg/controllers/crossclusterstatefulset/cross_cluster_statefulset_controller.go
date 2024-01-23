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
	"math"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
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
	"github.com/karmada-io/karmada/pkg/features"
	"github.com/karmada-io/karmada/pkg/resourceinterpreter"
	"github.com/karmada-io/karmada/pkg/sharedcli/ratelimiterflag"
	"github.com/karmada-io/karmada/pkg/util/helper"
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

// Reconcile performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (c *CrossClusterStatefulSetController) Reconcile(ctx context.Context, req controllerruntime.Request) (controllerruntime.Result, error) {
	klog.V(4).Infof("Reconciling ResourceBinding %s.", req.NamespacedName.String())

	binding := &workv1alpha2.ResourceBinding{}
	if err := c.Client.Get(ctx, req.NamespacedName, binding); err != nil {
		if apierrors.IsNotFound(err) {
			return controllerruntime.Result{}, nil
		}
		return controllerruntime.Result{Requeue: true}, err
	}
	return controllerruntime.Result{}, nil
}

func (c *CrossClusterStatefulSetController) detectFailure(clusters []string, tolerationSeconds *int32, key types.NamespacedName) (int32, []string) {
	var needEvictClusters []string
	duration := int32(math.MaxInt32)

	if duration == int32(math.MaxInt32) {
		duration = 0
	}
	return duration, needEvictClusters
}

func (c *CrossClusterStatefulSetController) evictBinding(binding *workv1alpha2.ResourceBinding, clusters []string) error {
	return nil
}

func (c *CrossClusterStatefulSetController) updateBinding(binding *workv1alpha2.ResourceBinding, allClusters sets.Set[string], needEvictClusters []string) error {
	if err := c.Update(context.TODO(), binding); err != nil {
		for _, cluster := range needEvictClusters {
			helper.EmitClusterEvictionEventForResourceBinding(binding, cluster, c.EventRecorder, err)
		}
		klog.ErrorS(err, "Failed to update binding", "binding", klog.KObj(binding))
		return err
	}
	for _, cluster := range needEvictClusters {
		allClusters.Delete(cluster)
	}
	if !features.FeatureGate.Enabled(features.GracefulEviction) {
		for _, cluster := range needEvictClusters {
			helper.EmitClusterEvictionEventForResourceBinding(binding, cluster, c.EventRecorder, nil)
		}
	}

	return nil
}

// SetupWithManager creates a controller and register to controller manager.
func (c *CrossClusterStatefulSetController) SetupWithManager(mgr controllerruntime.Manager) error {
	resourceBindingPredicateFn := predicate.Funcs{
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
		For(&appsv1alpha1.CrossClusterStatefulSet{}, builder.WithPredicates(resourceBindingPredicateFn)).
		WithOptions(controller.Options{RateLimiter: ratelimiterflag.DefaultControllerRateLimiter(c.RateLimiterOptions)}).
		Complete(c)
}
