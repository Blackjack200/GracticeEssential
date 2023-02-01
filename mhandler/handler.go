package mhandler

import "golang.org/x/exp/slices"

// Index returns the index of the first occurrence of v in s, or -1 if not
// present. Index accepts any type, as opposed to slices.Index, but might panic
// if E is not comparable.
func index[E any](s []E, v E) int {
	for i, vs := range s {
		if (any)(v) == (any)(vs) {
			return i
		}
	}
	return -1
}

// deleteVal deletes the first occurrence of a value in a slice of the type E
// and returns a new slice without the value.
func deleteVal[E any](s []E, v E) []E {
	if i := index(s, v); i != -1 {
		return slices.Clone(slices.Delete(s, i, i+1))
	}
	return s
}

func New() *MultipleHandler {
	return &MultipleHandler{}
}
