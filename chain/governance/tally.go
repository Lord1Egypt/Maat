// Package governance implements the deterministic proposal-tally rules from
// PLANNED_ECONOMY.md: quorum, approval, and the emergency fast-track threshold.
// Integer-only, consensus-safe.
package governance

const BpsDenom int64 = 10_000

// Params are the governance thresholds (basis points of the relevant base).
// Defaults match PLANNED_ECONOMY.md.
type Params struct {
	QuorumBps            int64 // min turnout vs total staked (default 3000 = 30%)
	ApprovalBps          int64 // min yes/(yes+no) (default 6000 = 60%)
	EmergencyApprovalBps int64 // emergency fast-track (default 8000 = 80%)
}

// DefaultParams returns the documented thresholds.
func DefaultParams() Params {
	return Params{QuorumBps: 3000, ApprovalBps: 6000, EmergencyApprovalBps: 8000}
}

// Tally is the vote weights (in staked MAAT). Abstain counts toward quorum but
// not toward the approval ratio.
type Tally struct {
	Yes     int64
	No      int64
	Abstain int64
}

// Turnout is all votes cast (the quorum base).
func (t Tally) Turnout() int64 { return t.Yes + t.No + t.Abstain }

// Outcome is the result of a tally.
type Outcome int

const (
	OutcomeRejectedNoQuorum Outcome = iota
	OutcomeRejected                 // quorum met, approval not reached
	OutcomePassed
)

func (o Outcome) String() string {
	switch o {
	case OutcomePassed:
		return "passed"
	case OutcomeRejected:
		return "rejected"
	default:
		return "rejected: no quorum"
	}
}

// Decide tallies a proposal against total staked supply. If emergency is true,
// the higher emergency approval threshold is applied.
func Decide(t Tally, totalStaked int64, p Params, emergency bool) Outcome {
	if totalStaked <= 0 {
		return OutcomeRejectedNoQuorum
	}
	// quorum: turnout must be at least QuorumBps of total staked
	if t.Turnout()*BpsDenom < totalStaked*p.QuorumBps {
		return OutcomeRejectedNoQuorum
	}
	// approval: yes / (yes+no), abstain excluded
	decisive := t.Yes + t.No
	if decisive == 0 {
		return OutcomeRejected
	}
	threshold := p.ApprovalBps
	if emergency {
		threshold = p.EmergencyApprovalBps
	}
	if t.Yes*BpsDenom >= decisive*threshold {
		return OutcomePassed
	}
	return OutcomeRejected
}
