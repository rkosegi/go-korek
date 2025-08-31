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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReconcileMap(t *testing.T) {
	type department struct {
		Name    string
		Manager string
	}
	mr := ForMap[string, department]().WithEqualityFunc(ReflectEqualityFunc[department]())

	org1 := map[string]department{
		"HR": {
			Name:    "HR",
			Manager: "Alice",
		},
		"IT": {
			Name:    "IT",
			Manager: "Bob",
		},
		"Sales": {
			Name:    "Sales",
			Manager: "Carl",
		},
	}

	org2 := map[string]department{
		"HR": {
			Name:    "HR",
			Manager: "Alice",
		},
		"IT": {
			Name:    "IT",
			Manager: "Charlie",
		},
		"Management": {
			Name:    "Management",
			Manager: "Dave",
		},
	}

	same, changed, onlyLeft, onlyRight := mr.Diff(org1, org2)
	assert.Len(t, same, 1)
	assert.Len(t, changed, 1)
	assert.Len(t, onlyLeft, 1)
	assert.Len(t, onlyRight, 1)
}
