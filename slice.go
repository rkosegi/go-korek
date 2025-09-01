/*
Copyright 2025 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconciler

import "github.com/samber/lo"

type sliceReconciler[T any] struct {
	idFn CompareFunc[T]
	eqFn CompareFunc[T]
}

func (s *sliceReconciler[T]) WithEqualityFunc(fn CompareFunc[T]) SliceReconciler[T] {
	s.eqFn = fn
	return s
}

func (s *sliceReconciler[T]) Diff(left []T, right []T) ([]T, []T, []T, []T) {
	hr := hybridSliceReconciler[T, T]{
		idFn: CompareDifferentFunc[T, T](s.idFn),
		eqFn: CompareDifferentFunc[T, T](s.eqFn),
	}
	return hr.Diff(left, right)
}

// ForSlice creates SliceReconciler.
// idFn is a function that is used to compare identity of 2 T items.
func ForSlice[T any](idFn CompareFunc[T]) SliceReconciler[T] {
	return &sliceReconciler[T]{
		eqFn: DefaultEqualityFunc[T](),
		idFn: idFn,
	}
}

type hybridSliceReconciler[T1, T2 any] struct {
	idFn CompareDifferentFunc[T1, T2]
	eqFn CompareDifferentFunc[T1, T2]
}

func (h *hybridSliceReconciler[T1, T2]) Diff(left []T1, right []T2) ([]T1, []T1, []T1, []T2) {
	var (
		same      []T1
		changed   []T1
		onlyLeft  []T1
		onlyRight []T2
	)

	for _, leftItem := range left {
		// try to find item in right slice based on identity predicate
		if rightItem, exists := lo.Find(right, func(other T2) bool {
			return h.idFn(leftItem, other)
		}); exists {
			// now if it exists, it could be equal by value as well
			if h.eqFn(leftItem, rightItem) {
				// so it's "same"
				same = append(same, leftItem)
			} else {
				// or it could be different
				changed = append(changed, leftItem)
			}
		} else {
			// nope, it only exists in left slice
			onlyLeft = append(onlyLeft, leftItem)
		}
	}
	for _, rightItem := range right {
		// try to find item in left slice based on identity predicate
		if _, exists := lo.Find(left, func(other T1) bool {
			return h.idFn(other, rightItem)
		}); !exists {
			// nope, it only exists in right slice
			onlyRight = append(onlyRight, rightItem)
		}
	}
	return same, changed, onlyLeft, onlyRight
}

// ForHybridSlices creates HybridSliceReconciler.
// Following required arguments must be provided:
//
//	idFn is a function that is used to compare an identity of an instance of T1 and T2.
//	eqFn is a function that is used to compare a "value" of an instance of T1 and T2.
func ForHybridSlices[T1 any, T2 any](idFn CompareDifferentFunc[T1, T2],
	eqFn CompareDifferentFunc[T1, T2]) HybridSliceReconciler[T1, T2] {
	return &hybridSliceReconciler[T1, T2]{
		idFn: idFn,
		eqFn: eqFn,
	}
}
