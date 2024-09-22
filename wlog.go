package wlog

import (
	"github.com/sirupsen/logrus"
)

// WLog is a wrap of entry
type WLog struct {
	factory *Factory
	*logrus.Entry
}

func (l WLog) Factory() *Factory {
	return l.factory
}
