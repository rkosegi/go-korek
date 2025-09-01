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

// this represents data structure from external source
type empExt struct {
	Name       string
	Department string
	Salary     int
}

// this represents data structure from internal source
type empInt struct {
	FullName   string
	DeptName   string
	Salary     int
	internalId int
}

func TestReconcileSameSlices(t *testing.T) {
	sr := ForSlice[empExt](func(left, right empExt) bool {
		return left.Name == right.Name
	}).WithEqualityFunc(DefaultEqualityFunc[empExt]())

	left := []empExt{
		{Name: "Alice", Department: "HR", Salary: 10},
		{Name: "Bob", Department: "IT", Salary: 10},
		{Name: "Charlie", Department: "Toilets", Salary: 99},
	}
	right := []empExt{
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

func TestReconcileHybridSlices(t *testing.T) {
	hr := ForHybridSlices[empExt, empInt](func(left empExt, right empInt) bool {
		return left.Name == right.FullName
	}, func(left empExt, right empInt) bool {
		return left.Salary == right.Salary && left.Name == right.FullName && left.Department == right.DeptName
	})

	left := []empExt{
		{Name: "Alice", Department: "HR", Salary: 10},
		{Name: "Bob", Department: "IT", Salary: 10},
		{Name: "Charlie", Department: "Toilets", Salary: 99},
	}
	right := []empInt{
		{FullName: "Alice", DeptName: "HR", Salary: 10},
		{FullName: "Bob", DeptName: "IT", Salary: 20},
		{FullName: "Cyril", DeptName: "Sales", Salary: 10},
		{FullName: "Dave", DeptName: "Management", Salary: 1},
	}
	same, changed, onlyLeft, onlyRight := hr.Diff(left, right)

	assert.Len(t, same, 1)
	assert.Equal(t, same[0].Name, "Alice")

	assert.Len(t, changed, 1)
	assert.Equal(t, changed[0].Name, "Bob")

	assert.Len(t, onlyRight, 2)
	assert.Equal(t, onlyRight[0].FullName, "Cyril")
	assert.Equal(t, onlyRight[1].FullName, "Dave")

	assert.Len(t, onlyLeft, 1)
	assert.Equal(t, onlyLeft[0].Name, "Charlie")
}
