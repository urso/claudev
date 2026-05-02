package ticket

import "iter"

// concatSeq yields all elements of each slice in order, without flattening
// into a fresh slice.
func concatSeq[T any](slices ...[]T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, s := range slices {
			for _, v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}
