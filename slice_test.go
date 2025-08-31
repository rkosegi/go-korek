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

func TestReconcileSlice(t *testing.T) {
	type emp struct {
		Name       string
		Department string
		Salary     int
	}
	sr := ForSlice[emp]().
		WithEqualityFunc(DefaultEqualityFunc[emp]()).
		WithIdentityFunc(func(left, right emp) bool {
			return left.Name == right.Name
		})

	left := []emp{
		{Name: "Alice", Department: "HR", Salary: 10},
		{Name: "Bob", Department: "IT", Salary: 10},
		{Name: "Charlie", Department: "Toilets", Salary: 99},
	}
	right := []emp{
		{Name: "Alice", Department: "HR", Salary: 10},
		{Name: "Bob", Department: "IT", Salary: 20},
		{Name: "Cyril", Department: "Sales", Salary: 10},
		{Name: "Dave", Department: "Management", Salary: 1},
	}
	same, changed, onlyLeft, onlyRight := sr.Diff(left, right)

	assert.Len(t, same, 1)
	assert.Equal(t, same[0].Name, "Alice")

	assert.Len(t, changed, 1)
	assert.Equal(t, changed[0].Name, "Bob")

	assert.Len(t, onlyRight, 2)
	assert.Equal(t, onlyRight[0].Name, "Cyril")
	assert.Equal(t, onlyRight[1].Name, "Dave")

	assert.Len(t, onlyLeft, 1)
	assert.Equal(t, onlyLeft[0].Name, "Charlie")
}
