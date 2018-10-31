package tools

import (
	"math/rand"
	"testing"
)

func TestLogSortKeys(t *testing.T) {
	sortedKeys := [...]string{
		"time",
		"level",
		"_id",
		"subscription",
		"resource_group",
		"account",
		"pool_name",
		"job_name",
		"msg",
	}

	shuffledKeys := make([]string, len(sortedKeys))
	perm := rand.Perm(len(sortedKeys))

	for i, v := range perm {
		shuffledKeys[v] = sortedKeys[i]
	}

	SortLogKeys(shuffledKeys)

	for i := range sortedKeys {
		if sortedKeys[i] != shuffledKeys[i] {
			t.Fatalf("Expected %s but got %s", sortedKeys, shuffledKeys)
		}
	}
}
