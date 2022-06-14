package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		metricReconciliationsStarted,
		metricReconciliationsErrored,
		metricReconciliationsCompleted,
	)
}
