package wlog

import (
	"context"
	"io"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

// LocalWLogMethod is used to specify the kinds of logger
type LocalWLogMethod string

// DevEnabled controls whether all dev loggers print to io.Discard, concurrent safe
var DevEnabled atomic.Bool

const (
	// LDev can be used to print debug messages
	LDev LocalWLogMethod = "dev"
	// LInit can be used to print init messages
	LInit LocalWLogMethod = "init"
	// LExit can be used to print exit messages
	LExit LocalWLogMethod = "exit"
)

var (
	localF       *Factory
	localDiscard *Factory
	localCtx     = context.Background()
)

func init() {
	var err error
	localF, err = NewFactory(createStdoutLogger())
	if err != nil {
		panic(err)
	}

	discardLogger := &logrus.Logger{
		Out:       io.Discard,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}
	localDiscard, err = NewFactory(discardLogger)
	if err != nil {
		panic(err)
	}

	DevEnabled.Store(true)
}

func (m LocalWLogMethod) String() string {
	return string(m)
}

// Log are used to print devOnly Logs, all results will be print to stdout
func (m LocalWLogMethod) Log(fingerPrints ...string) WLog {
	var factory *Factory
	if !DevEnabled.Load() {
		factory = localDiscard
	} else {
		factory = localF
	}

	wlog := factory.NewBuilder(localCtx).Name(fingerPrints...).Field(keyLocalMethod, m).Leaf()
	return wlog
}
