// Package slo provides Service Level Objective tracking for cronwatch jobs.
//
// A Tracker records whether each cron run met its expected deadline and
// computes a rolling compliance percentage over a configurable time window.
//
// Usage:
//
//	tr := slo.New(24*time.Hour, 500)
//	tr.Record("backup", time.Now(), true)  // run met SLO
//	tr.Record("backup", time.Now(), false) // run missed SLO
//
//	snap := tr.Snapshot("backup")
//	fmt.Printf("compliance: %.1f%%\n", snap.Compliance)
//
// The HTTP handler exposes snapshots at /slo and supports filtering
// by job name via the ?job= query parameter.
package slo
