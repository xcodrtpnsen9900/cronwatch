// Package checkpoint records the last successful completion timestamp for
// each monitored cron job.
//
// cronwatch uses checkpoint data to determine whether a job has missed its
// expected execution window after a process restart: if the stored LastOK
// time is older than the job's schedule interval, an alert is raised.
//
// Usage:
//
//	store := checkpoint.New()
//	store.Record("nightly-backup", time.Now())
//
//	if entry, ok := store.Get("nightly-backup"); ok {
//		fmt.Println(entry.LastOK, entry.RunCount)
//	}
package checkpoint
