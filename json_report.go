package main

import (
	"encoding/json"
	"slices"
	"maps"
	"os"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	sortedKeys := slices.Sorted(maps.Keys(pages))
	
	sortedPages := make(map[string]PageData)
	for _, key := range sortedKeys {
		sortedPages[key] = pages[key]
	}

	data, err := json.MarshalIndent(sortedPages, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}