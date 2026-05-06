// Package status exposes a lightweight HTTP handler that renders the
// current health of every monitored cron job as a JSON document.
//
// Typical usage:
//
//	http.Handle("/status", status.Handler(jobNames, historyStore))
//	http.ListenAndServe(":8080", nil)
//
// The response body is a JSON object with a "jobs" array and a
// "generated_at" timestamp, suitable for polling by dashboards or
// external health-check systems.
package status
