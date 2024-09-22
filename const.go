package wlog

import (
	"github.com/khicago/irr"
)

// KeyMethod is the key used to specify the method in the context
const KeyMethod = "method_"

// EntryKeyWLogSrc is the key used to specify the wlog source in the context
const (
	EntryKeyWLogSrc             = "wlog.src"
	EntryKeyWLogSrcValueEM      = "em"
	EntryKeyWLogSrcValueDefault = "default"
	EntryKeyWLogSrcValueLogger  = "logger"
)

// KeyFingerPrint is the key used to specify the fingerprint in the context
const KeyFingerPrint = "wlog.fp"

// keyLocalMethod is the key used for local logging methods
const keyLocalMethod = "wlog.local"

// defaultMethodValue is the default value for the method
const defaultMethodValue = "-"

type (
	// NodeStrategy define how to handle chain and columns when new node is created
	NodeStrategy int
)

const (
	// ForkLeaf mode
	// - chain: Entry Join, ctx not save path change
	// - column: Entry Combine, ctx not save new column
	ForkLeaf NodeStrategy = 0

	// ForkBranch mode
	// - chain: Entry and ctx chain Join
	// - column: Entry and ctx columns combine
	ForkBranch NodeStrategy = 1

	// NewTree mode (ctx will not create new, only chain and columns)
	// - chain: Entry create new chain, only keep new chain
	// - column: Entry create new columns, only keep new columns
	NewTree NodeStrategy = 2
)

var (
	// ErrLackOfEntryMakerOrLogger is returned when no EntryMaker or Logger is provided
	ErrLackOfEntryMakerOrLogger = irr.Error("invalid arguments, entryMakerOrLogger must be given")

	// ErrArgumentTypeNotMatch is returned when the argument type doesn't match the expected type
	ErrArgumentTypeNotMatch = irr.Error("invalid arguments: type error")
)
