// Package window implements a sliding-window counter with configurable
// bucket granularity.
//
// Use New to create a Counter, then Add to record events and Total to
// query how many events occurred within the most recent window duration.
//
// Example:
//
//	c := window.New(time.Minute, 12) // 12 × 5-second buckets
//	c.Add(1)
//	fmt.Println(c.Total()) // events in the last minute
//
// The counter is safe for concurrent use.
package window
