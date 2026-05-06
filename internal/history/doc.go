// Package history provides an in-memory ring-buffer store for cron job
// execution events. Each job retains up to a configurable maximum number
// of entries, with older entries automatically pruned.
//
// Typical usage:
//
//	store := history.New(50)
//	store.Record(history.Entry{
//		JobName:   "nightly-backup",
//		Timestamp: time.Now(),
//		Success:   true,
//		Message:   "completed in 4.2s",
//	})
//
//	if last, ok := store.Latest("nightly-backup"); ok {
//		fmt.Println(last.Success)
//	}
package history
