package utils

import (
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// returns table and column deignated by the sqlite3 error if they exists
func GetSqlite3UniqueErrorInfos(err sqlite3.Error) (string, string) {
	var	parts		[]string
	var	tableAndCol	[]string

	if err.Code == sqlite3.ErrConstraint &&
	err.ExtendedCode == sqlite3.ErrConstraintUnique {
		log.Println(err.Error())
		parts = strings.Split(err.Error(), ": ")
		if len(parts) <= 1 {
			return "", ""
		}
		tableAndCol = strings.Split(parts[1], ".")
		if len(tableAndCol) <= 1 {
			return "", ""
		} else {
			return tableAndCol[0], tableAndCol[1]
		}
	}
	return "", ""
}
