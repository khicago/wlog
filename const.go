package wlog

import (
	"sync/atomic"

	"github.com/khicago/irr"
)

// KeyMethod is the key used to specify the method in the context
const KeyMethod = "method_"

// KeyFingerPrint is the key used to specify the fingerprint in the context
const KeyFingerPrint = "finger_print_"

// keyLocalMethod is the key used for local logging methods
const keyLocalMethod = "local_"

// defaultMethodValue is the default value for the method
const defaultMethodValue = "-"

var (

	// CtxKeyCacheEntry is the key to cache log entry into a context
	CtxKeyCacheEntry = struct{ CtxKeyCacheEntry struct{} }{}

	// CtxKeyCacheMFP is the key to cache method and finger print into a context
	CtxKeyCacheMFP = struct{ CtxKeyCacheMFP struct{} }{}
)

var (
	// ErrLackOfEntryMakerOrLogger is returned when no EntryMaker or Logger is provided
	ErrLackOfEntryMakerOrLogger = irr.Error("invalid arguments, entryMakerOrLogger must be given")

	// ErrArgumentTypeNotMatch is returned when the argument type doesn't match the expected type
	ErrArgumentTypeNotMatch = irr.Error("invalid arguments: type error")
)

// DevEnabled controls whether all dev loggers print to ioutil.Discard, concurrent safe
var DevEnabled atomic.Bool
