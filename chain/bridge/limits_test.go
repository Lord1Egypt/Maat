package bridge

import "testing"

func TestSmallTxIsImmediate(t *testing.T) {
	l := NewLimiter(100, 100, 10, 20) // cap 100/window, large>=10, delay 20
	w, err := l.RequestOut(5, 1)
	if err != nil || w.Status != StatusImmediate {
		t.Fatalf("status=%v err=%v, want immediate", w.Status, err)
	}
	if err := l.Execute(int64(w.ID), 1); err != nil {
		t.Fatalf("execute immediate: %v", err)
	}
}

func TestLargeTxIsDelayedThenExecutes(t *testing.T) {
	l := NewLimiter(100, 100, 10, 20)
	w, err := l.RequestOut(10, 1)
	if err != nil || w.Status != StatusPending {
		t.Fatalf("status=%v err=%v, want pending", w.Status, err)
	}
	if w.UnlockAt != 21 {
		t.Fatalf("unlockAt=%d, want 21", w.UnlockAt)
	}
	if err := l.Execute(int64(w.ID), 10); err != ErrNotMatured {
		t.Fatalf("early execute err=%v, want ErrNotMatured", err)
	}
	if err := l.Execute(int64(w.ID), 21); err != nil {
		t.Fatalf("matured execute: %v", err)
	}
}

func TestCapExceeded(t *testing.T) {
	l := NewLimiter(100, 100, 1000, 20)
	if _, err := l.RequestOut(60, 1); err != nil {
		t.Fatalf("first out: %v", err)
	}
	if _, err := l.RequestOut(50, 2); err != ErrCapExceeded {
		t.Fatalf("err=%v, want ErrCapExceeded", err)
	}
}

func TestWindowResets(t *testing.T) {
	l := NewLimiter(100, 100, 1000, 20)
	if _, err := l.RequestOut(100, 1); err != nil {
		t.Fatalf("fill window: %v", err)
	}
	if _, err := l.RequestOut(1, 50); err != ErrCapExceeded {
		t.Fatalf("mid-window err=%v, want ErrCapExceeded", err)
	}
	if _, err := l.RequestOut(1, 101); err != nil {
		t.Fatalf("after reset: %v", err)
	}
}

func TestCancelRefundsCap(t *testing.T) {
	l := NewLimiter(100, 100, 10, 20)
	w, _ := l.RequestOut(80, 1) // large -> pending, reserves 80 of cap
	if rem := l.Remaining(1); rem != 20 {
		t.Fatalf("remaining=%d, want 20", rem)
	}
	if err := l.Cancel(int64(w.ID)); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if rem := l.Remaining(1); rem != 100 {
		t.Fatalf("remaining after cancel=%d, want 100", rem)
	}
	if err := l.Execute(int64(w.ID), 100); err != ErrNotPending {
		t.Fatalf("execute cancelled err=%v, want ErrNotPending", err)
	}
}

func TestBadAmount(t *testing.T) {
	l := NewLimiter(100, 100, 10, 20)
	if _, err := l.RequestOut(0, 1); err != ErrBadAmount {
		t.Fatalf("err=%v, want ErrBadAmount", err)
	}
}
