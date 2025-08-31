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

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// CompareFunc compares 2 values of the same type and return true if they are "same".
type CompareFunc[T any] func(left, right T) bool

// DefaultEqualityFunc uses github.com/google/go-cmp's Equal to perform equality check.
func DefaultEqualityFunc[T any]() CompareFunc[T] {
	return func(left, right T) bool {
		return cmp.Equal(left, right)
	}
}

// ReflectEqualityFunc uses reflect.DeepEqual to perform equality check.
func ReflectEqualityFunc[T any]() CompareFunc[T] {
	return func(left, right T) bool {
		return reflect.DeepEqual(left, right)
	}
}

// SliceReconciler can be used to reconcile state from 2 slices by computing difference
// between them based on identity and equality.
// Common use case is when you have 2 slices of identical type from different sources,
// and you need to take actions, based on their actual difference.
type SliceReconciler[T any] interface {
	// WithEqualityFunc sets CompareFunc that is used to compare equality of 2 items.
	// When omitted, then DefaultEqualityFunc is used to perform these equality checks.
	WithEqualityFunc(fn CompareFunc[T]) SliceReconciler[T]

	// WithIdentityFunc sets functions that is used to compare identity of 2 items.
	// When omitted, identity check is delegated to equality check.
	WithIdentityFunc(fn CompareFunc[T]) SliceReconciler[T]

	// Diff takes 2 input slices of T and return 4 slices:
	//  1, items that are "same" in both slices (according to equality func).
	//  2, items that exists in both input slices, but they are not *same* (according to equality func).
	//  3, items that only exists in left input slice (according to identity func).
	//  4, items that only exists in right input slice (according to identity func).
	Diff(left []T, right []T) (same []T, changed []T, onlyLeft []T, onlyRight []T)
}

type MapReconciler[K comparable, V any] interface {
	// WithEqualityFunc sets CompareFunc that is used to compare equality of 2 objects.
	// When omitted, then DefaultEqualityFunc is used to perform these equality checks.
	WithEqualityFunc(fn CompareFunc[V]) MapReconciler[K, V]

	// Diff takes 2 input maps (with same key type and value type) and produces 4 slices of key type:
	// 1, keys that maps to "same" value in both input maps (according to equality func).
	// 2, keys that maps to different value (according to equality func).
	// 3, keys that only exists in left input map.
	// 4, keys that only exists in right input map.
	Diff(left map[K]V, right map[K]V) (same []K, changed []K, onlyLeft []K, onlyRight []K)
}
