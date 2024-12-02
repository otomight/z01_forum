package database

import (
	"database/sql"
	"forum/internal/utils"
	"log"
)

func countRows(tableKey string) int {
	var	query	string
	var	count	int
	var	err		error

	query = `
		SELECT COUNT(*) FROM `+tableKey+`;
	`
	err = DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error at SQLCountRowInTable: %v\n", err)
		return 0
	}
	return count
}

func	generatePlaceHolders(nb int) string {
	var i int
	var result string

	i = 0
	for i < nb {
		result += "?"
		i++
		if i != nb {
			result += ", "
		}
	}
	return result
}

func	insertInto(data InsertIntoQuery) (sql.Result, error) {
	var	query			string
	var	values			[]any
	var	placeHolders	[]string
	var	i				int
	var	j				int

	for i = 0; i < len(data.Values); i++ {
		if len(data.Values[i]) != len(data.Keys) {
			continue
		}
		for j = 0; j < len(data.Values[i]); j++ {
			values = append(values, data.Values[i][j])
		}
		placeHolders = append(
			placeHolders, `(`+generatePlaceHolders(len(data.Values[i]))+`)`,
		)
	}
	query = `
		INSERT INTO `+data.Table+` (`+utils.ToCSV(data.Keys)+`)
		VALUES `+utils.ToCSV(placeHolders)+`
		`+data.Ending+`;
	`
	return DB.Exec(query, values...)
}
