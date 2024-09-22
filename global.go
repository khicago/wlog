package wlog

import (
	"context"
	"sync/atomic"

	"github.com/khicago/irr"
)

var defaultFactory atomic.Value

func init() {
	wlog, err := NewFactory(createTextLogger())
	if err != nil {
		panic(irr.Wrap(err, "failed to create default wlog instance"))
	}
	defaultFactory.Store(wlog)
}

func getDefaultFactory() *Factory {
	return defaultFactory.Load().(*Factory)
}

// SetEntryMaker sets the EntryMaker of default wlog instance
func SetEntryMaker(em EntryMaker) {
	getDefaultFactory().SetEntryMaker(em)
}

// By - create a new builder with the given context
// using .Build() method to create a new log entry
func By(ctx context.Context, fingerPrints ...string) *Builder {
	return getDefaultFactory().NewBuilder(ctx).Name(fingerPrints...)
}

// Leaf - create a log entry from the given context and fingerprints (using the default wlog instance)
func Leaf(ctx context.Context, fingerPrints ...string) WLog {
	return By(ctx, fingerPrints...).Leaf()
}

// Branch - create a log entry from the given context and fingerprints (using the default wlog instance)
func Branch(ctx context.Context, fingerPrints ...string) (WLog, context.Context) {
	return By(ctx, fingerPrints...).Branch()
}

// Detach - create a log entry from the given context and fingerprints (using the default wlog instance)
// method and fingerprint will be transferred to ctx, thus the mfp works in future
func Detach(ctx context.Context, fingerPrints ...string) (WLog, context.Context) {
	return By(ctx, fingerPrints...).Detach()
}

// Common create with given ctx and fingerprints (by default wlog instance)
func Common(fingerPrints ...string) WLog {
	return Leaf(localCtx, fingerPrints...)
}
