package announce

import (
	"sync"
	"testing"
	"time"
)

func TestPublish_DeliveresToSubscriber(t *testing.T) {
	b := New()
	var got Event
	b.Subscribe("s1", func(e Event) { got = e })

	b.Publish(Event{Job: "backup", Kind: KindSucceeded, Message: "ok"})

	if got.Job != "backup" {
		t.Fatalf("expected job=backup, got %q", got.Job)
	}
	if got.Kind != KindSucceeded {
		t.Fatalf("expected kind=succeeded, got %q", got.Kind)
	}
}

func TestPublish_SetsTimestampIfZero(t *testing.T) {
	b := New()
	var got Event
	b.Subscribe("ts", func(e Event) { got = e })

	before := time.Now().UTC()
	b.Publish(Event{Job: "j", Kind: KindStarted})
	after := time.Now().UTC()

	if got.OccurredAt.Before(before) || got.OccurredAt.After(after) {
		t.Fatalf("timestamp %v outside [%v, %v]", got.OccurredAt, before, after)
	}
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	b := New()
	var mu sync.Mutex
	var calls []string

	for _, name := range []string{"a", "b", "c"} {
		n := name
		b.Subscribe(n, func(Event) {
			mu.Lock()
			calls = append(calls, n)
			mu.Unlock()
		})
	}

	b.Publish(Event{Job: "x", Kind: KindFailed})

	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(calls))
	}
}

func TestUnsubscribe_StopsDelivery(t *testing.T) {
	b := New()
	count := 0
	b.Subscribe("s", func(Event) { count++ })
	b.Unsubscribe("s")

	b.Publish(Event{Job: "j", Kind: KindFailed})

	if count != 0 {
		t.Fatalf("expected 0 calls after unsubscribe, got %d", count)
	}
}

func TestSubscribe_ReplacesExisting(t *testing.T) {
	b := New()
	calls := 0
	b.Subscribe("s", func(Event) { calls++ })
	b.Subscribe("s", func(Event) { calls += 10 })

	b.Publish(Event{Job: "j", Kind: KindStarted})

	if calls != 10 {
		t.Fatalf("expected 10, got %d", calls)
	}
}

func TestLen_ReflectsSubscriberCount(t *testing.T) {
	b := New()
	if b.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	b.Subscribe("a", func(Event) {})
	b.Subscribe("b", func(Event) {})
	if b.Len() != 2 {
		t.Fatalf("expected 2, got %d", b.Len())
	}
	b.Unsubscribe("a")
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
}
