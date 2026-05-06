// Package eventlog implements a thread-safe, bounded ring buffer for recording
// cronwatch operational events such as alert dispatches, job completions, and
// internal errors.
//
// Events are stored in chronological order and are evicted oldest-first once
// the configured capacity is reached. Callers may retrieve all events or filter
// by severity level and/or job name.
//
// Example usage:
//
//	log := eventlog.New(500)
//	log.Add(eventlog.LevelInfo, "backup", "job started", nil)
//	log.Add(eventlog.LevelError, "backup", "exit code 1", map[string]string{"code": "1"})
//	events := log.Filter(eventlog.LevelError, "")
package eventlog
