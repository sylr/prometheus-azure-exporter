package tools

import (
	"sort"
)

// SortLogKeys sort log keys with fixed order for know keys
func SortLogKeys(keys []string) {
	sortingOrder := []string{
		"time",
		"level",
		"_id",
		"_interval",
		"_func",
		"subscription",
		"rg",
		"resource_group",
		"account",
		"metric",
		"pool",
		"pool_name",
		"job",
		"job_name",
	}

	sortedKeys := make([]string, 0, len(keys))
	unknownKeys := make([]string, 0)

ORDER:
	for _, o := range sortingOrder {
		for _, k := range keys {
			if o == k {
				sortedKeys = append(sortedKeys, k)
				continue ORDER
			}
		}
	}

KEYS:
	for _, k := range keys {
		for _, s := range sortedKeys {
			if k == s {
				continue KEYS
			}
		}

		unknownKeys = append(unknownKeys, k)
	}

	if len(unknownKeys) > 0 {
		sort.Strings(unknownKeys)
	}

	for i, k := range sortedKeys {
		keys[i] = k
	}

	for i, k := range unknownKeys {
		keys[len(sortedKeys)+i] = k
	}
}
