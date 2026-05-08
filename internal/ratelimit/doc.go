// Package ratelimit provides per-job alert rate limiting to prevent
// notification storms when a cron job repeatedly fails. Each job key
// is tracked independently with a configurable cooldown window.
package ratelimit
