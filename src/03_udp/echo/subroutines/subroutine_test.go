// go test -v -race subroutine_test.go subroutine.go

package subroutines

import (
	"context"
	"testing"
	"time"
)

func TestRoutine(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		t.Log("Main: Cancel")
		cancel()
		t.Log("Main: Cancelled")
	}()

	if err := routine(ctx); err != nil {
		t.Fatal(err)
	}

	t.Log("Main: Sleep 5...")
	time.Sleep(5 * time.Second)
	t.Log("Main: ...End")
}
