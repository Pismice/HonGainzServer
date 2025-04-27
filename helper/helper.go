package helper

import "database/sql"

// NullStringToPointer converts sql.NullString to a pointer or nil
func NullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
