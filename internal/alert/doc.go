// Package alert provides types and a builder for constructing structured alert
// payloads that describe cron job lifecycle events.
//
// Three event kinds are supported:
//
//   - missed: the job did not produce a heartbeat within its expected window.
//   - failed: the job explicitly reported an error or non-zero exit.
//   - recovered: a previously alerting job has returned to a healthy state.
//
// Payloads produced by Builder are intended to be serialised as JSON and
// forwarded to a webhook via the webhook package.
//
// Example:
//
//	b := alert.NewBuilder()
//	p := b.Missed("daily-backup", expectedAt)
//	// pass p to webhook.Sender.Send(ctx, p)
package alert
