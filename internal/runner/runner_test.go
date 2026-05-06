package runner_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/runner"
)

func newRunner() *runner.Runner {
	return runner.New("")
}

func TestRun_Success(t *testing.T) {
	r := newRunner()
	res := r.Run(context.Background(), "echo-job", "echo hello")

	if !res.Succeeded() {
		t.Fatalf("expected success, got err=%v exitCode=%d", res.Err, res.ExitCode)
	}
	if !strings.Contains(string(res.Output), "hello") {
		t.Errorf("expected output to contain 'hello', got %q", res.Output)
	}
	if res.JobName != "echo-job" {
		t.Errorf("expected JobName 'echo-job', got %q", res.JobName)
	}
}

func TestRun_Failure(t *testing.T) {
	r := newRunner()
	res := r.Run(context.Background(), "fail-job", "exit 2")

	if res.Succeeded() {
		t.Fatal("expected failure but got success")
	}
	if res.ExitCode != 2 {
		t.Errorf("expected exit code 2, got %d", res.ExitCode)
	}
}

func TestRun_CapturesOutput(t *testing.T) {
	r := newRunner()
	res := r.Run(context.Background(), "out-job", "echo cronwatch")

	if !strings.Contains(string(res.Output), "cronwatch") {
		t.Errorf("expected output to contain 'cronwatch', got %q", res.Output)
	}
}

func TestRun_ContextCancel(t *testing.T) {
	r := newRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res := r.Run(ctx, "sleep-job", "sleep 10")

	if res.Succeeded() {
		t.Fatal("expected context cancellation to cause failure")
	}
}

func TestRun_RecordsDuration(t *testing.T) {
	r := newRunner()
	res := r.Run(context.Background(), "dur-job", "echo ok")

	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}

func TestNew_DefaultShell(t *testing.T) {
	r := runner.New("")
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}
