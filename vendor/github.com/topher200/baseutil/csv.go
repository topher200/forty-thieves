package baseutil

import (
	"encoding/csv"
	"os"
	"strings"
)

// MapReader reads from a csv, creates a list of maps from the information.
//
// The first line of the file (the header line) provides the keys for each of
// the columns in the csv. Those keys are used, in order, when creating the map
// for each line.
//
// Strips leading and trailing whitespace from each value.
func MapReader(inputFilename string) []map[string]string {
	// Open the input file
	file, err := os.Open(inputFilename)
	Check(err)
	defer file.Close()

	// Read in the header line
	csvFile := csv.NewReader(file)
	columnNames, err := csvFile.Read()
	Check(err)

	// Read in all the rows, and populate a map of the <header column names>:value
	rowLines, err := csvFile.ReadAll()
	Check(err)
	rows := make([]map[string]string, len(rowLines))
	for rowNum, row := range rowLines {
		rows[rowNum] = make(map[string]string)
		for columnNum, value := range row {
			rows[rowNum][columnNames[columnNum]] = strings.TrimSpace(value)
		}
	}
	return rows
}
