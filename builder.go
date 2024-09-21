package wlog

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// WLogBuilder is used to construct complex logging options
type WLogBuilder struct {
	w            *WLog
	ctx          context.Context
	fingerPrints []string
	fields       logrus.Fields
	cache        bool
	removeCache  bool
}

// ----- builder instanciate method -----

// NewBuilder creates a new WLogBuilder instance
func (w *WLog) NewBuilder(ctx context.Context) *WLogBuilder {
	builder := builderPool.Get().(*WLogBuilder)
	builder.w = w
	builder.ctx = ctx
	// clear the previous data
	for k := range builder.fields {
		delete(builder.fields, k)
	}
	builder.fingerPrints = builder.fingerPrints[:0]
	builder.cache = false
	builder.removeCache = false
	return builder
}

// ----- public builder methods -----

// WithFingerPrints adds fingerprints to the builder
func (b *WLogBuilder) WithFingerPrints(fingerPrints ...string) *WLogBuilder {
	b.fingerPrints = fingerPrints
	return b
}

// WithField adds a single field to the builder
func (b *WLogBuilder) WithField(key string, value interface{}) *WLogBuilder {
	b.fields[key] = value
	return b
}

// WithFields adds multiple fields to the builder
func (b *WLogBuilder) WithFields(fields logrus.Fields) *WLogBuilder {
	for k, v := range fields {
		b.fields[k] = v
	}
	return b
}

// Cache sets whether to cache the log entry
func (b *WLogBuilder) Cache() *WLogBuilder {
	b.cache = true
	return b
}

// RemoveCache sets whether to remove the cached log entry
func (b *WLogBuilder) RemoveCache() *WLogBuilder {
	b.removeCache = true
	return b
}

// makeEntry creates a new log entry
func (b *WLogBuilder) makeEntry(ctx context.Context) *logrus.Entry {
	if ctx != nil {
		if l := ctx.Value(CtxKeyCacheEntry); l != nil {
			return l.(*logrus.Entry)
		}
	}

	if b.w.entryMaker != nil {
		return b.w.entryMaker(ctx)
	}

	// create a new Entry object, not using pool
	entry := logrus.NewEntry(b.w.defaultEntry.Logger)
	*entry = *b.w.defaultEntry
	entry = unboxMFPFromCtx(ctx, entry)

	return entry
}

// Build constructs and returns the log entry and context
func (b *WLogBuilder) Build() (Log, context.Context) {
	// add method and fingerprint to the entry, unless the cache is discarded, the fingerprint will follow the entry
	// and when the entry is discarded, the fingerprint should be extracted from the entry, and then add to the context
	entry := b.makeEntry(b.ctx)
	if entry == nil {
		return Log{}, b.ctx
	}

	entry = entry.WithFields(b.fields)
	// add fingerprints to the entry
	entry = insertFingerPrintToEntry(entry, b.fingerPrints)

	if b.cache {
		b.ctx = context.WithValue(b.ctx, CtxKeyCacheEntry, entry)
	}

	if b.removeCache {
		b.ctx = context.WithValue(b.ctx, CtxKeyCacheEntry, nil)

		// add method and fingerprint to the context, unless the cache is discarded, the fingerprint will follow the entry
		// and when the entry is discarded, the fingerprint should be extracted from the entry, and then add to the context
		b.ctx = boxMFPToCtx(b.ctx, entry)
	}

	defer builderPool.Put(b)
	return Log{entry}, b.ctx
}

// builderPool is used to recycle WLogBuilder objects
var builderPool = sync.Pool{
	New: func() interface{} {
		return &WLogBuilder{
			fields: make(logrus.Fields),
		}
	},
}

// ----- private methods -----

// insertFingerPrintToEntry adds fingerprints to the log entry
// If the method doesn't exist, it treats fingerPrints[0] as the method
func insertFingerPrintToEntry(entry *logrus.Entry, fingerPrints FingerPrints) *logrus.Entry {
	n := len(fingerPrints)
	if n == 0 {
		return entry
	}

	if v, exist := entry.Data[KeyMethod]; exist && v != nil && v != defaultMethodValue {
		fp := entry.Data[KeyFingerPrint]
		return entry.WithField(KeyFingerPrint, mustCombineFingerPrint(fp, fingerPrints))
	}

	return entry.WithField(KeyMethod, fingerPrints[0]).WithField(KeyFingerPrint, fingerPrints[1:])
}

// unboxMFPFromCtx retrieves method and fingerprint from context and adds them to the log entry
func unboxMFPFromCtx(ctx context.Context, entry *logrus.Entry) *logrus.Entry {
	if ctx == nil {
		return entry
	}

	cacheMFP := mustCombineFingerPrint(ctx.Value(CtxKeyCacheMFP), nil)
	if cacheMFP == nil {
		return entry
	}

	if n := len(cacheMFP); n != 0 {
		return insertFingerPrintToEntry(entry, cacheMFP)
	}

	return entry
}

// boxMFPToCtx adds method and fingerprint to the context
func boxMFPToCtx(ctx context.Context, entry *logrus.Entry) context.Context {
	method, ok := entry.Data[KeyMethod]
	if !ok || method == nil || method == "-" {
		return ctx
	}

	methodStr, ok := method.(string)
	if !ok {
		return ctx
	}

	head := FingerPrints{methodStr}

	fingerPrint, ok := entry.Data[KeyFingerPrint]
	if !ok || fingerPrint == nil {
		return context.WithValue(ctx, CtxKeyCacheMFP, head)
	}

	fingerPrintArr, ok := fingerPrint.(FingerPrints)
	if !ok {
		return context.WithValue(ctx, CtxKeyCacheMFP, head)
	}

	return context.WithValue(ctx, CtxKeyCacheMFP, mustCombineFingerPrint(head, fingerPrintArr))
}
