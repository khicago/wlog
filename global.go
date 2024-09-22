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

// SetEntryGetter sets the EntryMaker of default wlog instance
func SetEntryGetter(em EntryMaker) {
	getDefaultFactory().SetEntryMaker(em)
}

// Leaf - create a log entry from the given context and fingerprints (using the default wlog instance)
func Leaf(ctx context.Context, fingerPrints ...string) WLog {
	wlog, _ := getDefaultFactory().NewBuilder(ctx).
		WithStrategy(ForkLeaf).
		WithFingerPrints(fingerPrints...).
		Build()
	return wlog
}

// Branch - create a log entry from the given context and fingerprints (using the default wlog instance)
func Branch(ctx context.Context, fingerPrints ...string) (WLog, context.Context) {
	return getDefaultFactory().NewBuilder(ctx).
		WithStrategy(ForkBranch).
		WithFingerPrints(fingerPrints...).
		Build()
}

// DetachNew - create a log entry from the given context and fingerprints (using the default wlog instance)
// method and fingerprint will be transferred to ctx, thus the mfp works in future
func DetachNew(ctx context.Context, fingerPrints ...string) (WLog, context.Context) {
	return getDefaultFactory().NewBuilder(ctx).
		WithStrategy(NewTree).
		WithFingerPrints(fingerPrints...).
		Build()
}

// Common create with given ctx and fingerprints (by default wlog instance)
func Common(fingerPrints ...string) WLog {
	wlog, _ := getDefaultFactory().NewBuilder(context.Background()).
		WithStrategy(ForkLeaf).
		WithFingerPrints(fingerPrints...).
		Build()
	return wlog
}

// B - create a new builder with the given context
// using .Build() method to create a new log entry
func B(ctx context.Context) *Builder {
	return getDefaultFactory().NewBuilder(ctx)
}
