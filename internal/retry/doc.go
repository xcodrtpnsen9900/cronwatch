// Package retry provides a simple, context-aware retry mechanism with
// configurable backoff for use across cronwatch components.
//
// Usage:
//
//	policy := retry.DefaultPolicy() // 3 attempts, 2s base delay, 2x multiplier
//	err := retry.Do(ctx, policy, func() error {
//		return webhook.Send(payload)
//	})
//	if errors.Is(err, retry.ErrMaxAttempts) {
//		log.Println("webhook delivery failed after all retries")
//	}
package retry
