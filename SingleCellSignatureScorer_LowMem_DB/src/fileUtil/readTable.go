// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Written by Frederic PONT.
//(c) Frederic Pont 2018

package fileUtil

import (
	"ScorerLowMem/src/userTypes"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/akrylysov/pogreb"
	"github.com/schollz/progressbar"
)

func ReadTable(fileDB *pogreb.DB, path string, database map[string][]string, nbCPU int, conf userTypes.CONF) {

	sortedPWnames := SortKeysMap(database)

	fileLength, err := LineCounter(path)
	check(err)

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// Assume we don't know the number of fields per line.  By setting
	// FieldsPerRecord negative, each row may have a variable number of fields.
	reader.FieldsPerRecord = -1
	reader.Comma = '\t'

	// read column names = gene names
	colNames, error := reader.Read()
	if error != nil {
		log.Fatal(error)
	}
	colNames = colNames[1:] // empty first cell in table

	var wg sync.WaitGroup

	bar := progressbar.New(fileLength - 1) // Add a new progress bar
	for count := 1; count < fileLength; count++ {
		if conf.Server == 0 {
			bar.Add(1) // show progress bar
		}
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		cellName := line[0]
		geneValues := line[1:]
		wg.Add(1)
		go ProcessLine(fileDB, cellName, colNames, sortedPWnames, geneValues, database, conf, &wg) // process one line = one cell against all pathways
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println()
}
