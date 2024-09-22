package wlog

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// EntryMaker is the function which specifies how wlog creates a logger associated with a context
type EntryMaker func(ctx context.Context) *logrus.Entry

// Factory is a sandbox for the logger
type Factory struct {
	entryMaker   EntryMaker
	defaultEntry *logrus.Entry
	mu           sync.RWMutex
}

// SetEntryMaker updates the EntryMaker of the Factory instance
// It returns the updated Factory instance
func (f *Factory) SetEntryMaker(em EntryMaker) *Factory {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entryMaker = em
	return f
}

// Logger returns the underlying logrus.Logger
func (f *Factory) Logger() *logrus.Logger {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.defaultEntry.Logger
}

// SetLevel sets the logging level for the Factory instance
func (f *Factory) SetLevel(level logrus.Level) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	f.Logger().SetLevel(level)
}

// LoggerSource defines the types that can be used to create a Factory
type LoggerSource interface {
	EntryMaker | *logrus.Entry | *logrus.Logger
}

// NewFactory create a new Factory instance
// T can be EntryMaker, *logrus.Entry, or *logrus.Logger
func NewFactory[T LoggerSource](source T) (*Factory, error) {
	switch v := any(source).(type) {
	case nil:
		return nil, ErrLackOfEntryMakerOrLogger
	case EntryMaker:
		return newFactoryWithEntryMaker(v), nil
	case *logrus.Entry:
		return newFactoryWithEntry(v), nil
	case *logrus.Logger:
		return newFactoryWithLogger(v), nil
	default:
		return nil, ErrArgumentTypeNotMatch
	}
}

// newFactoryWithEntryMaker create a new Factory with EntryMaker
func newFactoryWithEntryMaker(em EntryMaker) *Factory {
	return &Factory{entryMaker: em}
}

// newFactoryWithEntry create a new Factory with *logrus.Entry
func newFactoryWithEntry(entry *logrus.Entry) *Factory {
	return &Factory{defaultEntry: entry}
}

// newFactoryWithLogger create a new Factory with *logrus.Logger
func newFactoryWithLogger(logger *logrus.Logger) *Factory {
	return &Factory{
		defaultEntry: logger.WithField(KeyMethod, defaultMethodValue),
	}
}

// makeEntry create a new logrus.Entry
// this method will create Entry based on the way the Factory is initialized
func (f *Factory) makeEntry(ctx context.Context) *logrus.Entry {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.entryMaker != nil {
		// if Factory is initialized with EntryMaker, use EntryMaker to create Entry
		return f.entryMaker(ctx).WithField(EntryKeyWLogSrc, EntryKeyWLogSrcValueEM)
	}

	if f.defaultEntry != nil {
		// if Factory is initialized with *logrus.Entry, create new Entry based on default Entry
		return f.defaultEntry.WithContext(ctx).WithField(EntryKeyWLogSrc, EntryKeyWLogSrcValueDefault)
	}

	// if Factory is neither initialized with EntryMaker nor defaultEntry,
	// it is likely initialized with *logrus.Logger (or not initialized correctly)
	// in this case, we create a basic Entry
	logger := f.Logger()
	if logger == nil {
		// if even Logger is not set, create a new one
		logger = logrus.New()
	}

	// create a new Entry, with default Method field and context
	return logger.WithField(EntryKeyWLogSrc, EntryKeyWLogSrcValueLogger).WithContext(ctx)
}
