package utils

import (
	"database/sql"
	"fmt"
)

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

func RowsToSlice(rows *sql.Rows) ([]any, error) {
    columns, err := rows.ColumnTypes()
    if (err != nil) {
        return nil, err
    }


        values := make([]interface{}, len(columns))
        rows.Scan(values...)
        fmt.Print(values)
    // for rows.Next() {
    //     values := make([]interface{}, len(columns))
    //     rows.Scan(values...)
    //     fmt.Print(values)
    // }
    return []any{}, nil;
}
