package wlog

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Builder is used to construct complex logging options
type Builder struct {
	factory   *Factory
	ctx       context.Context
	chainNode []string
	columns   Columns
	strategy  NodeStrategy
}

// ----- new method -----

// NewBuilder creates a new Builder instance
func (f *Factory) NewBuilder(ctx context.Context) *Builder {
	builder := builderPool.Get().(*Builder)
	builder.factory = f
	builder.ctx = ctx
	builder.columns = builder.columns[:0]
	builder.chainNode = builder.chainNode[:0]
	builder.strategy = ForkLeaf // 默认策略
	return builder
}

// ----- public builder methods -----

// WithFingerPrints adds fingerprints to the builder
func (b *Builder) WithFingerPrints(fingerPrints ...string) *Builder {
	b.chainNode = fingerPrints
	return b
}

// WithField adds a single field to the builder
func (b *Builder) WithField(key string, value any) *Builder {
	b.columns.Set(Column{Key: key, Value: value})
	return b
}

// WithFields adds multiple columns to the builder
func (b *Builder) WithFields(fields map[string]any) *Builder {
	b.columns.Set(ColumnsFromFields(fields)...)
	return b
}

// WithStrategy sets the cache strategy for the builder
func (b *Builder) WithStrategy(strategy NodeStrategy) *Builder {
	b.strategy = strategy
	return b
}

// Build makes a WLog instance from the builder
func (b *Builder) Build() (WLog, context.Context) {

	// read chain and columns from context
	chainInCtx := ChainFromCtx(b.ctx)
	columnsInCtx := ColumnsFromCtx(b.ctx)

	var chainForEntry Chain
	var columnsForEntry Columns
	var newCtx context.Context

	switch b.strategy {
	case ForkLeaf:
		// merge entry and ctx
		chainForEntry = chainInCtx.Join(b.chainNode)
		columnsForEntry = columnsInCtx.Combine(b.columns)
		newCtx = b.ctx
	case ForkBranch:
		// merge entry and ctx
		chainForEntry = chainInCtx.Join(b.chainNode)
		columnsForEntry = columnsInCtx.Combine(b.columns)
		newCtx = chainForEntry.WriteCtx(b.ctx)
		newCtx = columnsForEntry.WriteCtx(newCtx)
	case NewTree:
		// only use new chain and columns
		chainForEntry = b.chainNode
		columnsForEntry = b.columns
		newCtx = chainForEntry.WriteCtx(b.ctx)
		newCtx = columnsForEntry.WriteCtx(newCtx)
	default: // default strategy is ForkLeaf
		chainForEntry = chainInCtx.Join(b.chainNode)
		columnsForEntry = columnsInCtx.Combine(b.columns)
		newCtx = b.ctx
	}

	// make new entry
	entry := b.factory.makeEntry(b.ctx)
	// add fields and fingerprints to entry
	entry = columnsForEntry.WriteEntry(entry)
	entry = chainForEntry.WriteEntry(entry)

	// make WLog instance
	wlog := WLog{
		Entry:   entry,
		factory: b.factory,
	}

	// put builder instance back to pool
	defer builderPool.Put(b)

	return wlog, newCtx
}

// helper function: merge fields
func MergeFields(existing, new logrus.Fields) logrus.Fields {
	merged := make(logrus.Fields, len(existing)+len(new))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range new {
		merged[k] = v
	}
	return merged
}

// builderPool is used to recycle Builder objects
var builderPool = sync.Pool{
	New: func() any {
		return &Builder{
			columns: make(Columns, 0),
		}
	},
}
