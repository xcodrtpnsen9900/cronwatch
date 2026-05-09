// Package shadow implements a shadow-mode execution runner for cronwatch.
//
// Shadow mode allows a secondary implementation of a job to run concurrently
// alongside the primary implementation. The shadow result is recorded and
// logged but never affects alerts or the caller's return value. This is useful
// for safely validating new job implementations before promoting them to
// production without risking missed-run alerts.
//
// Usage:
//
//	runner := shadow.New(logger)
//	err := runner.Run(ctx, "my-job", primaryFn, candidateFn)
//	// err is always from primaryFn; candidateFn outcome is logged only.
package shadow
