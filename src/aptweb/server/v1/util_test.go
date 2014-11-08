package v1

import (
	"sort"
	"testing"
)

func assert(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected: %v, but received: %v", expected, actual)
	}
}

func assertArrayEqual(t *testing.T, expected []string, actual []string) {
	assert(t, len(expected), len(actual))
	for index, _ := range expected {
		assert(t, expected[index], actual[index])
	}
}

func TestUnique(t *testing.T) {
	input := []string{"one", "two", "one", "three"}
	actual := unique(input)

	expected := []string{"one", "two", "three"}

	sort.Strings(expected)
	sort.Strings(actual)
	assertArrayEqual(t, expected, actual)
}

func TestSplitPackages(t *testing.T) {
	packages := "one two three two three"
	actual := splitPackages(packages)

	expected := []string{"one", "two", "three"}

	sort.Strings(expected)
	sort.Strings(actual)
	assertArrayEqual(t, expected, actual)
}
