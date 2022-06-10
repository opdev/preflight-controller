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

// PreflightCheckStatusReconciler reconciles the deployment resource.
type PreflightCheckStatusReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.opdev.io,resources=preflightchecks/finalizers,verbs=update
// Reconcile will ensure that the Kubernetes Job for PreflightCheck
// reaches the desired state.
func (r *PreflightCheckStatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithName("status-reconciler")
	l.Info("status reconciliation initiated.")
	defer l.Info("status reconciliation complete.")
	instanceKey := req.NamespacedName

	// Get the PreflightCheck instance to make sure it still exists.
	var instance toolsv1alpha1.PreflightCheck
	err := r.Client.Get(ctx, instanceKey, &instance)
	if apierrors.IsNotFound(err) {
		return subrec.Evaluate(subrec.DoNotRequeue())
	}

	if err != nil {
		return subrec.Evaluate(subrec.RequeueWithError(err))
	}

	var existing batchv1.JobList
	if err = r.Client.List(
		ctx,
		&existing,
		client.MatchingLabels{"spawning-resource-name": instanceKey.Name},
	); err != nil {
		return subrec.Evaluate(subrec.RequeueWithError(err))
	}

	newStatus := toolsv1alpha1.PreflightCheckStatus{
		Type: instance.CheckType(),
	}

	if len(existing.Items) > 0 {
		jobList := make([]string, len(existing.Items))
		var completed bool
		for i, job := range existing.Items {
			jobList[i] = job.GetName()
			if job.Status.Succeeded > 0 {
				completed = true
			}
		}

		newStatus.Completed = completed
		newStatus.Jobs = jobList
	}

	instance.Status = newStatus
	if err := r.Status().Update(ctx, &instance); err != nil {
		return subrec.Evaluate(subrec.RequeueWithError(err))
	}

	return subrec.Evaluate(subrec.DoNotRequeue()) // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *PreflightCheckStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.PreflightCheck{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
