// Package bridge implements Ma'at's bridge-out safety controls — the mitigations
// for the #1 attack vector in crypto (bridge hacks). Deterministic, height-based,
// integer-only (consensus-safe).
//
// Two layers:
//   1. A rolling per-asset withdrawal CAP (max units out per WindowBlocks).
//   2. A DELAY queue: any withdrawal >= LargeTxUnits is held for DelayBlocks,
//      during which governance / the security council can CANCEL it.
package bridge

import "errors"

var (
	ErrBadAmount    = errors.New("bridge: amount must be positive")
	ErrCapExceeded  = errors.New("bridge: withdrawal cap exceeded for window")
	ErrNotFound     = errors.New("bridge: pending withdrawal not found")
	ErrNotMatured   = errors.New("bridge: delay not elapsed")
	ErrNotPending   = errors.New("bridge: withdrawal not in pending state")
)

type Status uint8

const (
	StatusImmediate Status = iota // below large-tx threshold: execute now
	StatusPending                 // queued, waiting out the delay
	StatusExecuted
	StatusCancelled
)

// Withdrawal is a recorded bridge-out request.
type Withdrawal struct {
	ID          uint64
	Amount      int64
	RequestedAt int64 // block height
	UnlockAt    int64 // block height it may be executed
	Status      Status
}

// Limiter holds the per-asset safety state. Construct via NewLimiter.
type Limiter struct {
	DailyCapUnits int64 // max units out per window
	WindowBlocks  int64 // window length in blocks
	LargeTxUnits  int64 // at/above this, a withdrawal is delayed
	DelayBlocks   int64 // delay applied to large withdrawals

	windowStart  int64
	usedInWindow int64
	nextID       uint64
	pending      map[uint64]*Withdrawal
}

func NewLimiter(dailyCap, windowBlocks, largeTx, delayBlocks int64) *Limiter {
	return &Limiter{
		DailyCapUnits: dailyCap,
		WindowBlocks:  windowBlocks,
		LargeTxUnits:  largeTx,
		DelayBlocks:   delayBlocks,
		pending:       make(map[uint64]*Withdrawal),
	}
}

// rollWindow resets the used counter when the current window has elapsed.
func (l *Limiter) rollWindow(height int64) {
	if l.windowStart == 0 {
		l.windowStart = height
		return
	}
	for height >= l.windowStart+l.WindowBlocks {
		l.windowStart += l.WindowBlocks
		l.usedInWindow = 0
	}
}

// Remaining returns the cap left in the current window.
func (l *Limiter) Remaining(height int64) int64 {
	l.rollWindow(height)
	return l.DailyCapUnits - l.usedInWindow
}

// RequestOut reserves cap and either authorizes immediate execution (small tx)
// or queues a delayed, cancellable withdrawal (large tx). Cap is consumed at
// request time so queued large txs can't be used to bypass the cap.
func (l *Limiter) RequestOut(amount, height int64) (*Withdrawal, error) {
	if amount <= 0 {
		return nil, ErrBadAmount
	}
	l.rollWindow(height)
	if l.usedInWindow+amount > l.DailyCapUnits {
		return nil, ErrCapExceeded
	}
	l.usedInWindow += amount
	l.nextID++
	w := &Withdrawal{ID: l.nextID, Amount: amount, RequestedAt: height}
	if amount >= l.LargeTxUnits {
		w.UnlockAt = height + l.DelayBlocks
		w.Status = StatusPending
	} else {
		w.UnlockAt = height
		w.Status = StatusImmediate
	}
	l.pending[w.ID] = w
	return w, nil
}

// Execute finalizes a withdrawal once its delay has elapsed.
func (l *Limiter) Execute(id, height int64) error {
	w, ok := l.pending[uint64(id)]
	if !ok {
		return ErrNotFound
	}
	if w.Status != StatusPending && w.Status != StatusImmediate {
		return ErrNotPending
	}
	if height < w.UnlockAt {
		return ErrNotMatured
	}
	w.Status = StatusExecuted
	return nil
}

// Cancel aborts a pending withdrawal and refunds its reserved cap. Only large,
// still-pending withdrawals can be cancelled (the delay window is what makes a
// bridge exploit catchable).
func (l *Limiter) Cancel(id int64) error {
	w, ok := l.pending[uint64(id)]
	if !ok {
		return ErrNotFound
	}
	if w.Status != StatusPending {
		return ErrNotPending
	}
	w.Status = StatusCancelled
	l.usedInWindow -= w.Amount // refund cap
	if l.usedInWindow < 0 {
		l.usedInWindow = 0
	}
	return nil
}

// Get returns a recorded withdrawal.
func (l *Limiter) Get(id int64) (*Withdrawal, bool) {
	w, ok := l.pending[uint64(id)]
	return w, ok
}

func (l *Limiter) GetWindowStart() int64 {
	return l.windowStart
}

func (l *Limiter) SetWindowStart(val int64) {
	l.windowStart = val
}

func (l *Limiter) GetUsedInWindow() int64 {
	return l.usedInWindow
}

func (l *Limiter) SetUsedInWindow(val int64) {
	l.usedInWindow = val
}

func (l *Limiter) GetNextID() uint64 {
	return l.nextID
}

func (l *Limiter) SetNextID(val uint64) {
	l.nextID = val
}

func (l *Limiter) GetPending() map[uint64]*Withdrawal {
	return l.pending
}

func (l *Limiter) SetPending(val map[uint64]*Withdrawal) {
	l.pending = val
}
