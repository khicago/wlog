package wlog

import (
	"context"
	"sync"
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

// builderPool is used to recycle Builder objects
var builderPool = sync.Pool{
	New: func() any {
		return &Builder{
			columns: make(Columns, 0),
		}
	},
}

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

// Name adds fingerprints to the builder
func (b *Builder) Name(chainNodes ...string) *Builder {
	b.chainNode = chainNodes
	return b
}

// Field adds a single field to the builder
func (b *Builder) Field(key string, value any) *Builder {
	b.columns.Set(Column{Key: key, Value: value})
	return b
}

// Fields adds multiple columns to the builder
func (b *Builder) Fields(fields Fields) *Builder {
	b.columns.Set(ColumnsFromFields(fields)...)
	return b
}

// Strategy sets the cache strategy for the builder
func (b *Builder) Strategy(strategy NodeStrategy) *Builder {
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

func (b *Builder) Branch() (WLog, context.Context) {
	return b.Strategy(ForkBranch).Build()
}

func (b *Builder) Leaf() WLog {
	l, _ := b.Strategy(ForkLeaf).Build()
	return l
}

func (b *Builder) Detach() (WLog, context.Context) {
	return b.Strategy(NewTree).Build()
}
