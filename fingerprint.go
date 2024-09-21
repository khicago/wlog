package wlog

import "strings"

// FingerPrints is a slice of strings representing fingerprints
type FingerPrints []string

// String returns the string representation of the fingerprints
func (fp FingerPrints) String() string {
	if len(fp) == 0 {
		return "/"
	}
	var builder strings.Builder
	builder.WriteString("/")
	for i, s := range fp {
		if i > 0 {
			builder.WriteString("/")
		}
		builder.WriteString(s)
	}
	return builder.String()
}

// mustCombineFingerPrint combines the given fingerprints
func mustCombineFingerPrint(fp any, appends FingerPrints) FingerPrints {
	if nil == fp {
		return appends // might returns nil
	}

	fpArr, ok := fp.(FingerPrints)
	if !ok {
		return appends
	}

	if appends == nil {
		return fpArr
	}

	nA := len(appends)
	if nA == 0 {
		return fpArr
	}

	// 预分配切片容量，避免多次扩容
	result := make(FingerPrints, 0, len(fpArr)+nA)
	result = append(result, fpArr...)
	result = append(result, appends...)
	return result
}
