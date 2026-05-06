// Package circuitbreaker provides a thread-safe circuit breaker for
// protecting outbound webhook calls in cronwatch.
//
// Usage:
//
//	br := circuitbreaker.New(5, 30*time.Second)
//
//	if err := br.Allow(); err != nil {
//		// circuit is open — skip the call
//		return err
//	}
//	if err := sendWebhook(payload); err != nil {
//		br.RecordFailure()
//		return err
//	}
//	br.RecordSuccess()
//
// The breaker moves through three states:
//   - Closed  — calls are allowed normally.
//   - Open    — calls are blocked until the cooldown elapses.
//   - HalfOpen — one probe call is allowed; success closes, failure reopens.
package circuitbreaker
