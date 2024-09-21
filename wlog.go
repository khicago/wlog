package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

// EntryMaker is the function which specifies how wlog creates a logger associated with a context
type EntryMaker func(ctx context.Context) *logrus.Entry

// WLog is a sandbox for the logger
type WLog struct {
	entryMaker   EntryMaker
	defaultEntry *logrus.Entry
}

// SetEntryMaker updates the EntryMaker of the WLog instance
// It returns the updated WLog instance
func (w *WLog) SetEntryMaker(em EntryMaker) *WLog {
	w.entryMaker = em
	return w
}

// Logger returns the underlying logrus.Logger
func (w *WLog) Logger() *logrus.Logger {
	return w.defaultEntry.Logger
}

// SetLevel sets the logging level for the WLog instance
func (w *WLog) SetLevel(level logrus.Level) {
	w.Logger().SetLevel(level)
}

// NewWLog creates a new WLog instance
// The argument can be EntryMaker, *logrus.Logger or nil
func NewWLog(entryMakerOrLogger any) (*WLog, error) {
	if entryMakerOrLogger == nil {
		return nil, ErrLackOfEntryMakerOrLogger
	}

	if em, ok := entryMakerOrLogger.(EntryMaker); ok {
		return &WLog{
			entryMaker: em,
		}, nil
	}

	if entry, ok := entryMakerOrLogger.(*logrus.Entry); ok {
		return &WLog{
			defaultEntry: entry,
		}, nil
	}

	if logger, ok := entryMakerOrLogger.(*logrus.Logger); ok {
		return &WLog{
			defaultEntry: logger.WithField(KeyMethod, defaultMethodValue),
		}, nil
	}

	return nil, ErrArgumentTypeNotMatch
}

// --------------------------------------------------- wlog - logger create methods ---------------------------------------------------

// ByCtx creates a new log entry associated with the given context
func (w *WLog) ByCtx(ctx context.Context, fingerPrints ...string) Log {
	log, _ := w.NewBuilder(ctx).WithFingerPrints(fingerPrints...).Build()
	return Log{Entry: log.Entry}
}

// Common creates a new log entry with an empty context
func (w *WLog) Common(fingerPrints ...string) Log {
	log, _ := w.NewBuilder(context.Background()).WithFingerPrints(fingerPrints...).Build()
	return Log{Entry: log.Entry}
}

// ByCtxAndCache is a convenience method that returns a log entry and caches it in the context
// It's recommended to use the NewBuilder method for more flexibility
func (w *WLog) ByCtxAndCache(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	log, newCtx := w.NewBuilder(ctx).WithFingerPrints(fingerPrints...).Cache().Build()
	return Log{Entry: log.Entry}, newCtx
}

// ByCtxAndRemoveCache is a convenience method that returns a log entry and removes the cache from the context
// It's recommended to use the NewBuilder method for more flexibility
func (w *WLog) ByCtxAndRemoveCache(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	log, newCtx := w.NewBuilder(ctx).WithFingerPrints(fingerPrints...).RemoveCache().Build()
	return Log{Entry: log.Entry}, newCtx
}
