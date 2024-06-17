package utils

import tea "github.com/charmbracelet/bubbletea"

// import (
// 	"database/sql"
// 	"fmt"
// )

func MakeCustomCommand(data interface{}) tea.Cmd {
	return func() tea.Msg { return data } 
}

/*
GetSupportedVariants returns a map containing supported SQL database variants as keys.

This method initializes a map with the supported variant types and returns it for further use in determining if a given string
(representing the type of SQL database) is present within the map.

Returns:

	map[string]struct{}: A map with supported SQL database variant types as keys.

Example:

	supportedVariants := GetSupportedVariants()
	_, isSupported := supportedVariants["postgres"]
*/
func GetSupportedVariants() map[string]struct{} {
	variant_types := [...]string{"postgres", "mysql"}
	variants := make(map[string]struct{}, len(variant_types))

	for _, value := range variant_types {
		variants[value] = struct{}{}
	}

	return variants
}

// from https://stackoverflow.com/a/50025091/14759055
func Map(vs []map[string]string, f func(map[string]string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
