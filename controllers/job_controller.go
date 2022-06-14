package controllers

import (
	"context"

	toolsv1alpha1 "github.com/opdev/preflight-controller/api/v1alpha1"
	subrec "github.com/opdev/subreconciler"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PreflightCheckJobReconciler reconciles the deployment resource.
type PreflightCheckJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=update

// Reconcile will ensure that the Kubernetes Job for PreflightCheck
// reaches the desired state.
func (r *PreflightCheckJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	metricReconciliationsStarted.Inc()
	l := log.FromContext(ctx).WithName("job-reconciler")
	l.Info("job reconciliation initiated.")
	defer l.Info("job reconciliation complete.")
	instanceKey := req.NamespacedName

	// Get the PreflightCheck instance to make sure it still exists.
	var instance toolsv1alpha1.PreflightCheck
	err := r.Client.Get(ctx, instanceKey, &instance)

	if apierrors.IsNotFound(err) {
		metricReconciliationsCompleted.Inc()
		return subrec.Evaluate(subrec.DoNotRequeue())
	}

	if err != nil {
		metricReconciliationsErrored.Inc()
		return subrec.Evaluate(subrec.RequeueWithError(err))
	}

	new := instance.GenerateJob()

	err = ctrl.SetControllerReference(&instance, &new, r.Scheme)
	if err != nil {
		metricReconciliationsErrored.Inc()
		return subrec.Evaluate(subrec.RequeueWithError(err))
	}

	// If the job exists, get it and patch it
	// var existing batchv1.Job
	// err = r.Client.Get(ctx, client.ObjectKeyFromObject(&new), &existing)

	var existing batchv1.JobList
	err = r.Client.List(
		ctx,
		&existing,
		client.MatchingLabels{"spawning-resource-name": instanceKey.Name},
	)

	// if apierrors.IsNotFound(err) {
	if len(existing.Items) == 0 {
		// create the resource because it does not exist.
		l.Info("creating resource", new.Kind, new.Name)
		if err := r.Client.Create(ctx, &new); err != nil {
			metricReconciliationsErrored.Inc()
			return subrec.Evaluate(subrec.RequeueWithError(err))
		}

		if err != nil {
			metricReconciliationsErrored.Inc()
			return subrec.Evaluate(subrec.RequeueWithError(err))
		}
	}

	// l.Info("updating resources if necessary", existing.Kind, existing.GetName())
	// patchDiff := client.MergeFrom(&existing)
	// if err = mergo.Merge(&existing, new, mergo.WithOverride); err != nil {
	// 	return subrec.Evaluate(subrec.RequeueWithError(err))
	// }

	// if err = r.Patch(ctx, &existing, patchDiff); err != nil {
	// 	return subrec.Evaluate(subrec.RequeueWithError(err))
	// }

	metricReconciliationsCompleted.Inc()
	return subrec.Evaluate(subrec.DoNotRequeue()) // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *PreflightCheckJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.PreflightCheck{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
