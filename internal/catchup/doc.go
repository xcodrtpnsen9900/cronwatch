// Package catchup provides catch-up detection for cronwatch.
//
// When the cronwatch process restarts after a period of downtime, some
// scheduled cron windows may have been missed entirely. The Detector in
// this package compares the last recorded checkpoint time for each job
// against the current wall-clock time, walking the expected schedule
// forward to enumerate every fire time that was skipped.
//
// Typical usage:
//
//	det := catchup.New()
//	missed := det.Scan(job.Name, lastCheckpoint, scheduler.Next)
//	for _, m := range missed {
//		// send alert or record event
//	}
package catchup
