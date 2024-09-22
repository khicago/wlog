package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

// WLog is a wrap of entry
type WLog struct {
	factory *Factory
	*logrus.Entry
}

// WithField adds a single field to the logger (Leaf strategy)
func (l WLog) WithField(key string, value any) WLog {
	builder := l.factory.NewBuilder(l.Context).
		WithStrategy(ForkLeaf).
		WithField(key, value)

	wlog, _ := builder.Build()
	return wlog
}

// WithFields adds multiple fields to the logger (Leaf strategy)
func (l WLog) WithFields(fields Fields) WLog {
	builder := l.factory.NewBuilder(l.Context).
		WithStrategy(ForkLeaf).
		WithFields(fields)

	wlog, _ := builder.Build()
	return wlog
}

// WithBranchField adds a single field to the logger and updates the context (Branch strategy)
func (l WLog) BranchField(ctx context.Context, key string, value any) (WLog, context.Context) {
	builder := l.factory.NewBuilder(ctx).
		WithStrategy(ForkBranch).
		WithField(key, value)

	return builder.Build()
}

// WithBranchFields adds multiple fields to the logger and updates the context (Branch strategy)
func (l WLog) BranchFields(ctx context.Context, fields Fields) (WLog, context.Context) {
	builder := l.factory.NewBuilder(ctx).
		WithStrategy(ForkBranch).
		WithFields(fields)

	return builder.Build()
}

// Branch returns a new logger append chainNode to the origin logger
func (l WLog) Branch(ctx context.Context, fingerPrints ...string) (WLog, context.Context) {
	builder := l.factory.NewBuilder(ctx).
		WithStrategy(ForkBranch).
		WithFingerPrints(fingerPrints...)

	return builder.Build()
}
