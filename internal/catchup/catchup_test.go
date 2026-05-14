package catchup

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// everyMinute advances t by one minute — simulates a "* * * * *" schedule.
func everyMinute(t time.Time) time.Time {
	return t.Add(time.Minute)
}

func TestScan_NoMissedRuns(t *testing.T) {
	// now == lastRun, so nextFn(lastRun) is in the future
	d := newWithClock(func() time.Time { return epoch })
	missed := d.Scan("job1", epoch, everyMinute)
	if len(missed) != 0 {
		t.Fatalf("expected 0 missed runs, got %d", len(missed))
	}
}

func TestScan_OneMissedRun(t *testing.T) {
	now := epoch.Add(90 * time.Second) // 1.5 minutes after epoch
	d := newWithClock(func() time.Time { return now })
	missed := d.Scan("job1", epoch, everyMinute)
	if len(missed) != 1 {
		t.Fatalf("expected 1 missed run, got %d", len(missed))
	}
	if !missed[0].Scheduled.Equal(epoch.Add(time.Minute)) {
		t.Errorf("unexpected scheduled time: %v", missed[0].Scheduled)
	}
	if missed[0].Job != "job1" {
		t.Errorf("unexpected job name: %s", missed[0].Job)
	}
}

func TestScan_MultipleMissedRuns(t *testing.T) {
	now := epoch.Add(5 * time.Minute)
	d := newWithClock(func() time.Time { return now })
	missed := d.Scan("job2", epoch, everyMinute)
	if len(missed) != 5 {
		t.Fatalf("expected 5 missed runs, got %d", len(missed))
	}
}

func TestAll_AccumulatesAcrossScans(t *testing.T) {
	now := epoch.Add(3 * time.Minute)
	d := newWithClock(func() time.Time { return now })
	d.Scan("jobA", epoch, everyMinute)
	d.Scan("jobB", epoch, everyMinute)
	all := d.All()
	if len(all) != 6 {
		t.Fatalf("expected 6 total missed runs, got %d", len(all))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	now := epoch.Add(2 * time.Minute)
	d := newWithClock(func() time.Time { return now })
	d.Scan("jobA", epoch, everyMinute)
	a := d.All()
	a[0].Job = "mutated"
	b := d.All()
	if b[0].Job == "mutated" {
		t.Error("All() should return an independent copy")
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	now := epoch.Add(2 * time.Minute)
	d := newWithClock(func() time.Time { return now })
	d.Scan("jobA", epoch, everyMinute)
	d.Clear()
	if len(d.All()) != 0 {
		t.Error("expected empty store after Clear")
	}
}

func TestScan_DetectedAtIsNow(t *testing.T) {
	now := epoch.Add(2 * time.Minute)
	d := newWithClock(func() time.Time { return now })
	missed := d.Scan("jobA", epoch, everyMinute)
	for _, m := range missed {
		if !m.DetectedAt.Equal(now) {
			t.Errorf("DetectedAt %v != now %v", m.DetectedAt, now)
		}
	}
}
