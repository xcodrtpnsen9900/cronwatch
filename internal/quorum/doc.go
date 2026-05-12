// Package quorum provides a distributed-quorum tracker for cron jobs.
//
// In environments where multiple worker processes execute the same
// scheduled job, a single heartbeat may not be sufficient to confirm
// a successful run. The quorum package allows operators to declare
// that at least N distinct instances must report success within a
// rolling time window before the job is considered healthy.
//
// Usage:
//
//	q := quorum.New(5 * time.Minute)
//	q.Require("nightly-backup", 3)
//
//	// called by each worker after a successful run:
//	q.Report("nightly-backup", workerID)
//
//	// checked by the monitor:
//	status := q.Check("nightly-backup")
//	if !status.Met {
//		// fire alert
//	}
package quorum
