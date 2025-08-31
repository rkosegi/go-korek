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

func (s *sliceReconciler[T]) WithIdentityFunc(fn CompareFunc[T]) SliceReconciler[T] {
	s.idFn = fn
	return s
}

func (s *sliceReconciler[T]) Diff(left []T, right []T) ([]T, []T, []T, []T) {
	var (
		same      []T
		changed   []T
		onlyLeft  []T
		onlyRight []T
	)
	for _, leftItem := range left {
		// try to find item in right slice based on identity predicate
		if rightItem, exists := lo.Find(right, func(other T) bool {
			return s.idFn(leftItem, other)
		}); exists {
			// now if it exists, it could be equal by value as well
			if s.eqFn(leftItem, rightItem) {
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
		if _, exists := lo.Find(left, func(other T) bool {
			return s.idFn(rightItem, other)
		}); !exists {
			// nope, it only exists in right slice
			onlyRight = append(onlyRight, rightItem)
		}
	}
	return same, changed, onlyLeft, onlyRight
}

func ForSlice[T any]() SliceReconciler[T] {
	return &sliceReconciler[T]{
		eqFn: DefaultEqualityFunc[T](),
	}
}
