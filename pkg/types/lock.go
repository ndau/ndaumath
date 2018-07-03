package types

import (
	"errors"

	"github.com/attic-labs/noms/go/marshal"
	nt "github.com/attic-labs/noms/go/types"
	util "github.com/oneiro-ndev/noms-util"
)

//go:generate msgp

// Lock keeps track of an account's Lock information
type Lock struct {
	NoticePeriod Duration `msg:"notice"`
	// if a lock has not been notified, this is nil
	UnlocksOn *Timestamp `msg:"unlock"`
}

var _ marshal.Marshaler = (*Lock)(nil)
var _ marshal.Unmarshaler = (*Lock)(nil)

// MarshalNoms implements Marshaler for lock
func (l Lock) MarshalNoms(vrw nt.ValueReadWriter) (val nt.Value, err error) {
	return marshal.Marshal(vrw, l.toNomsLock())
}

// UnmarshalNoms implements Unmarshaler for lock
func (l *Lock) UnmarshalNoms(v nt.Value) error {
	nl := nomsLock{}
	err := marshal.Unmarshal(v, &nl)
	if err != nil {
		return err
	}
	l.fromNomsLock(nl)
	return nil
}

type nomsLock struct {
	Duration   util.Int
	IsNotified bool
	UnlocksOn  util.Int
}

func (l Lock) toNomsLock() nomsLock {
	nl := nomsLock{
		Duration:   util.Int(l.NoticePeriod),
		IsNotified: l.UnlocksOn != nil,
	}
	if l.UnlocksOn != nil {
		nl.UnlocksOn = util.Int(*l.UnlocksOn)
	}
	return nl
}

func (l *Lock) fromNomsLock(nl nomsLock) {
	l.NoticePeriod = Duration(nl.Duration)
	if nl.IsNotified {
		ts := Timestamp(nl.UnlocksOn)
		l.UnlocksOn = &ts
	} else {
		l.UnlocksOn = nil
	}
}

// Notify updates this lock with notification of intent to unlock
func (l *Lock) Notify(blockTime Timestamp, weightedAverageAge Duration) error {
	if l.UnlocksOn != nil {
		return errors.New("already notified")
	}
	uo := blockTime.Add(l.NoticePeriod)
	l.UnlocksOn = &uo
	return nil
}
