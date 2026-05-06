// Package healthcheck exposes a composable HTTP health-check handler for
// cronwatch.
//
// Usage:
//
//	checker := healthcheck.New()
//	checker.Register("webhook", func() error {
//		// return non-nil to mark the check as failing
//		return nil
//	})
//	http.Handle("/healthz", checker.Handler())
//
// The handler returns HTTP 200 with status "ok" when all checks pass, or
// HTTP 503 with status "degraded" when any check returns an error.
// Each check result is included in the JSON response body under "checks".
package healthcheck
