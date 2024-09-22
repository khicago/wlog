package wlog

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	// Chain is a slice of strings representing fingerprints
	Chain []string
)

// CtxKeyChain is the key to cache fingerprint into a context
var CtxKeyChain = struct{ CtxKeyChain struct{} }{}

// String returns the string representation of the fingerprints
func (cc Chain) String() string {
	if len(cc) == 0 {
		return "/"
	}
	var builder strings.Builder
	builder.WriteString("/")
	for i, s := range cc {
		if i > 0 {
			builder.WriteString("/")
		}
		builder.WriteString(s)
	}
	return builder.String()
}

// Join with the given fingerprints
func (cc Chain) Join(appends Chain) Chain {
	if nil == cc {
		return appends // might returns nil
	}

	if appends == nil {
		return cc
	}

	nA := len(appends)
	if nA == 0 {
		return cc
	}

	// pre-allocate slice capacity to avoid multiple resizes
	result := make(Chain, 0, len(cc)+nA)
	result = append(result, cc...)
	result = append(result, appends...)
	return result
}

// WriteEntry write fingerprints to entry
func (cc Chain) WriteEntry(entry *logrus.Entry) *logrus.Entry {
	return entry.WithField(KeyFingerPrint, cc)
}

// WriteCtx cache fingerprints to context
func (cc Chain) WriteCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxKeyChain, cc)
}

// ChainFromEntry read chain from entry
func ChainFromEntry(entry *logrus.Entry) (Chain, bool) {
	val, ok := entry.Data[KeyFingerPrint]
	if !ok || val == nil {
		return nil, false
	}
	fp, ok := val.(Chain)
	if !ok {
		return nil, false
	}
	return fp, true
}

// ChainFromCtx get cached chain from context
func ChainFromCtx(ctx context.Context) Chain {
	if fp := ctx.Value(CtxKeyChain); fp != nil {
		return fp.(Chain)
	}
	return nil
}

// DetachChain detach chain from context
func DetachChain(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxKeyChain, nil)
}

// fpEntry2Ctx add fingerprints to context
func fpEntry2Ctx(ctx context.Context, entry *logrus.Entry) context.Context {
	fp, ok := ChainFromEntry(entry)
	if !ok {
		return ctx
	}
	return fp.WriteCtx(ctx)
}

// fpCtx2Entry add fingerprints to entry
func fpCtx2Entry(ctx context.Context, entry *logrus.Entry) *logrus.Entry {
	fp := ChainFromCtx(ctx)
	if fp == nil {
		return entry
	}
	return fp.WriteEntry(entry)
}
