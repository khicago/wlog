package wlog

import (
	"context"
	"sync"

	"github.com/khicago/irr"
)

// DefaultWLog the default wlog instance
var (
	DefaultWLog      *WLog
	defaultWLogMutex sync.RWMutex
)

func init() {
	wlog, err := NewWLog(createTextLogger())
	if err != nil {
		panic(irr.Wrap(err, "failed to create default wlog instance"))
	}
	DefaultWLog = wlog
	DevEnabled.Store(true)
}

// SetEntryGetter sets the EntryMaker of default wlog instance
func SetEntryGetter(em EntryMaker) {
	defaultWLogMutex.Lock()
	defer defaultWLogMutex.Unlock()
	DefaultWLog.SetEntryMaker(em)
}

// Common create with given ctx and fingerprints (by default wlog instance)
func Common(fingerPrints ...string) Log {
	defaultWLogMutex.RLock()
	defer defaultWLogMutex.RUnlock()
	l := DefaultWLog.Common(fingerPrints...)
	return l
}

// From - create a log entry from the given context and fingerprints (using the default wlog instance)
func From(ctx context.Context, fingerPrints ...string) Log {
	return DefaultWLog.ByCtx(ctx, fingerPrints...)
}

// FromHold - create a log entry from the given context and fingerprints (using the default wlog instance)
func FromHold(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	return DefaultWLog.ByCtxAndCache(ctx, fingerPrints...)
}

// FromRelease - create a log entry from the given context and fingerprints (using the default wlog instance)
// method and finger print will be transferred to ctx, thus the mfp works in future
func FromRelease(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	return DefaultWLog.ByCtxAndRemoveCache(ctx, fingerPrints...)
}

// Builder - create a new builder with the given context
// using .Build() method to create a new log entry
func Builder(ctx context.Context) *WLogBuilder {
	return DefaultWLog.NewBuilder(ctx)
}
