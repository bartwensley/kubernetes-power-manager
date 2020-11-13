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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	powerv1alpha1 "gitlab.devtools.intel.com/OrchSW/CNO/power-operator.git/api/v1alpha1"
)

// ProfileReconciler reconciles a Profile object
type ProfileReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=power.intel.com,resources=profiles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=power.intel.com,resources=profiles/status,verbs=get;update;patch

func (r *ProfileReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("profile", req.NamespacedName)
	logger.Info("Reconciling Profile")

	profile := &powerv1alpha1.Profile{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, profile)
	if err != nil {
		if errors.IsNotFound(err) {
			// TODO: everything (will elaborate on this)
			logger.Info(fmt.Sprintf("Profile %v has been deleted, cleaning up...", req.NamespacedName))
			// Returning a nil error will stop the requeue
			return ctrl.Result{}, nil
		}
		// Requeue the request
		return ctrl.Result{}, err
	}

	logger.Info("Specs for Profile:")
	logger.Info(fmt.Sprintf("Name: %s", profile.Spec.Name))
	logger.Info(fmt.Sprintf("Max: %d", profile.Spec.Max))
	logger.Info(fmt.Sprintf("Min: %d", profile.Spec.Min))

	if profile.Spec.Name == "Shared" {
		logger.Info("This Profile has been designated as the Shared Profile...")
		logger.Info("Creating Config for Shared Profile...")
	} else {
		logger.Info("This Profile is exclusive, skipping Config creation...")
	}

	return ctrl.Result{}, nil
}

func (r *ProfileReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&powerv1alpha1.Profile{}).
		Complete(r)
}