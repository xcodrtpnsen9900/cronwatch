// Package sampling provides adaptive event sampling for cronwatch alerts.
//
// A Sampler tracks per-key event rates and suppresses events that exceed
// a configurable burst limit within a rolling time window. This prevents
// alert storms when a cron job fails repeatedly in quick succession.
//
// Example usage:
//
//	sampler := sampling.New(sampling.DefaultPolicy())
//	if sampler.Allow(jobName) {
//		// forward the alert
//	}
package sampling
