package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricReconciliationsStarted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "preflightcheck_job_reconciliations_started",
		Help: "The total number of job reconciliations started.",
	})

	metricReconciliationsCompleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "preflightcheck_job_reconciliations_completed",
		Help: "The total number of job reconciliations that reach are ended without needing a requeue.",
	})

	metricReconciliationsErrored = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "preflightcheck_job_reconciliations_errored",
		Help: "The total number of job reconciliations that reach an error before completion.",
	})
)
