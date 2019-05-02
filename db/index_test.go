package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertWithIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.True(t, exists, "Index not stored in database at the expected key")
}

func TestUpdateWithIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	oldKeyExists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.False(t, oldKeyExists, "Old index was still stored after update")
	updatedKeyExists, err := db.ldb.Has([]byte("index:people:age:43:foo"), nil)
	require.NoError(t, err)
	assert.True(t, updatedKeyExists, "Index not stored in database at the updated key")
}

func TestFindWithValue(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	// expected is a set of testModels with Age = 42
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "ExpectedPerson_" + strconv.Itoa(i),
			Age:  42,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}

	// We also insert some other models with a different age.
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "OtherPerson_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithValue(ageIndex, []byte("42"), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithRange(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	all := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		all = append(all, model)
	}
	// expected is the set of people with 1 <= age < 4
	expected := all[1:4]

	var actual []*testModel
	require.NoError(t, col.FindWithRange(ageIndex, []byte("1"), []byte("4"), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithPrefix(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	// expected is a set of testModels with an age that starts with "2"
	expected := []*testModel{
		{
			Name: "ExpectedPerson_0",
			Age:  2021,
		},
		{
			Name: "ExpectedPerson_1",
			Age:  22,
		},
		{
			Name: "ExpectedPerson_2",
			Age:  250,
		},
	}
	for _, model := range expected {
		require.NoError(t, col.Insert(model))
	}

	// We also insert some other models with different ages.
	excluded := []*testModel{
		{
			Name: "ExcludedPerson_0",
			Age:  40,
		},
		{
			Name: "ExcludedPerson_1",
			Age:  41,
		},
		{
			Name: "ExcludedPerson_2",
			Age:  42,
		},
	}
	for _, model := range excluded {
		require.NoError(t, col.Insert(model))
	}

	{
		var actual []*testModel
		require.NoError(t, col.FindWithPrefix(ageIndex, []byte("2"), &actual))
		assert.Equal(t, expected, actual)
	}
	{
		// An empty prefix should return all models.
		all := append(expected, excluded...)
		var actual []*testModel
		require.NoError(t, col.FindWithPrefix(ageIndex, []byte{}, &actual))
		assert.Equal(t, all, actual)
	}
}

func TestFindWithValueWithMultiIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	nicknameIndex := col.AddMultiIndex("nicknames", func(m Model) [][]byte {
		person := m.(*testModel)
		indexValues := make([][]byte, len(person.Nicknames))
		for i, nickname := range person.Nicknames {
			indexValues[i] = []byte(nickname)
		}
		return indexValues
	})

	// expected is a set of testModels that include the nickname "Bob"
	expected := []*testModel{
		{
			Name:      "ExpectedPerson_0",
			Age:       42,
			Nicknames: []string{"Bob", "Jim", "John"},
		},
		{
			Name:      "ExpectedPerson_1",
			Age:       43,
			Nicknames: []string{"Alice", "Bob", "Emily"},
		},
		{
			Name:      "ExpectedPerson_2",
			Age:       44,
			Nicknames: []string{"Bob", "No one"},
		},
	}
	for _, model := range expected {
		require.NoError(t, col.Insert(model))
	}

	// We also insert some other models with different nicknames.
	excluded := []*testModel{
		{
			Name:      "ExcludedPerson_0",
			Age:       42,
			Nicknames: []string{"Bill", "Jim", "John"},
		},
		{
			Name:      "ExcludedPerson_1",
			Age:       43,
			Nicknames: []string{"Alice", "Jane", "Emily"},
		},
		{
			Name:      "ExcludedPerson_2",
			Age:       44,
			Nicknames: []string{"Nemo", "No one"},
		},
	}
	for _, model := range excluded {
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithValue(nicknameIndex, []byte("Bob"), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithRangeWithMultiIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	nicknameIndex := col.AddMultiIndex("nicknames", func(m Model) [][]byte {
		person := m.(*testModel)
		indexValues := make([][]byte, len(person.Nicknames))
		for i, nickname := range person.Nicknames {
			indexValues[i] = []byte(nickname)
		}
		return indexValues
	})

	// expected is a set of testModels that include at least one nickname that
	// satisfies "B" <= nickname < "E"
	expected := []*testModel{
		{
			Name:      "ExpectedPerson_0",
			Age:       42,
			Nicknames: []string{"Alice", "Beth", "Emily"},
		},
		{
			Name:      "ExpectedPerson_1",
			Age:       43,
			Nicknames: []string{"Bob", "Charles", "Dan"},
		},
		{
			Name:      "ExpectedPerson_2",
			Age:       44,
			Nicknames: []string{"James", "Darell"},
		},
	}
	for _, model := range expected {
		require.NoError(t, col.Insert(model))
	}

	// We also insert some other models with different nicknames.
	excluded := []*testModel{
		{
			Name:      "ExcludedPerson_0",
			Age:       42,
			Nicknames: []string{"Allen", "Jim", "John"},
		},
		{
			Name:      "ExcludedPerson_1",
			Age:       43,
			Nicknames: []string{"Sophia", "Jane", "Emily"},
		},
		{
			Name:      "ExcludedPerson_2",
			Age:       44,
			Nicknames: []string{"Nemo", "No one"},
		},
	}
	for _, model := range excluded {
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithRange(nicknameIndex, []byte("B"), []byte("E"), &actual))
	assert.Equal(t, expected, actual)
}
