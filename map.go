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

type mapReconciler[K comparable, V any] struct {
	eqFn CompareFunc[V]
}

func (m *mapReconciler[K, V]) Diff(left map[K]V, right map[K]V) ([]K, []K, []K, []K) {
	var (
		same      []K
		changed   []K
		onlyLeft  []K
		onlyRight []K
	)

	for k, v := range left {
		// try to find item in right map, based on value of K
		if x, exists := right[k]; exists {
			if m.eqFn(v, x) {
				same = append(same, k)
			} else {
				changed = append(changed, k)
			}
		} else {
			onlyLeft = append(onlyLeft, k)
		}
	}
	for k := range right {
		if _, exists := left[k]; !exists {
			onlyRight = append(onlyRight, k)
		}
	}

	return same, changed, onlyLeft, onlyRight
}

func (m *mapReconciler[K, V]) WithEqualityFunc(fn CompareFunc[V]) MapReconciler[K, V] {
	m.eqFn = fn
	return m
}

func ForMap[K comparable, V any]() MapReconciler[K, V] {
	return &mapReconciler[K, V]{
		eqFn: DefaultEqualityFunc[V](),
	}
}
