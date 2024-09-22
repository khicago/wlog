package wlog

import (
	"context"
	"sort"

	"github.com/sirupsen/logrus"
)

type (
	// Column represents a key-value pair
	Column struct {
		Key   string
		Value any
	}

	// Columns is a slice of Column
	Columns []Column

	Fields = map[string]any
)

// CtxKeyColumns is the key to cache Columns in context
var CtxKeyColumns = struct{ CtxKeyColumns struct{} }{}

// WriteEntry write columns to entry
func (c Columns) WriteEntry(entry *logrus.Entry) *logrus.Entry {
	fields := c.ToFields()
	return entry.WithFields(fields)
}

// WriteCtx cache columns to context
func (c Columns) WriteCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxKeyColumns, c)
}

// Combine merge current columns with new columns
func (c Columns) Combine(other Columns) Columns {
	return c.Set(other...)
}

// ToFields convert columns to fields
func (c Columns) ToFields() Fields {
	fields := make(Fields, len(c))
	for _, col := range c {
		fields[col.Key] = col.Value
	}
	return fields
}

// Sorted keep order, and assume it's mostly ordered
func (c Columns) Sorted() Columns {
	if len(c) <= 1 {
		return c
	}

	// find the first unsorted position
	unsortedIndex := 1
	for unsortedIndex < len(c) && c[unsortedIndex-1].Key <= c[unsortedIndex].Key {
		unsortedIndex++
	}

	// if already sorted, return directly
	if unsortedIndex == len(c) {
		return c
	}

	// sort the unsorted part
	sort.Slice(c[unsortedIndex:], func(i, j int) bool {
		return c[unsortedIndex+i].Key < c[unsortedIndex+j].Key
	})

	// merge sorted part and new sorted part
	return merge(c[:unsortedIndex], c[unsortedIndex:])
}

// Set add or update columns, keep order and unique, assume it's mostly ordered
func (c Columns) Set(cols ...Column) Columns {
	// if original slice is empty, return new columns directly
	if len(c) == 0 {
		return cols
	}

	// mostly ordered and less items, try to modify in place
	c = c.Sorted()

	for _, col := range cols {
		index := sort.Search(len(c), func(i int) bool {
			return c[i].Key >= col.Key
		})

		if index < len(c) && c[index].Key == col.Key {
			// update existing key
			c[index] = col
		} else {
			// insert new key
			c = append(c, Column{})
			copy(c[index+1:], c[index:])
			c[index] = col
		}
	}

	return c
}

// ColumnsFromFields create columns from map
func ColumnsFromFields(fields Fields) Columns {
	columns := make([]Column, 0, len(fields))
	for k, v := range fields {
		columns = append(columns, Column{Key: k, Value: v})
	}
	return columns
}

// ColumnsFromCtx get cached columns from context
func ColumnsFromCtx(ctx context.Context) Columns {
	if cols := ctx.Value(CtxKeyColumns); cols != nil {
		return cols.(Columns)
	}
	return nil
}

// ---- private ----

// merge merge two sorted columns
func merge(left, right Columns) Columns {
	result := make(Columns, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Key <= right[j].Key {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
