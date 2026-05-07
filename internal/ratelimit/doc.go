// Package ratelimit provides per-job alert rate limiting to prevent
// alert storms when a cron job repeatedly fails. Each job has an
// independent cooldown window; once an alert is sent, subsequent
// alerts for the same job are suppressed until the cooldown expires.
package ratelimit
