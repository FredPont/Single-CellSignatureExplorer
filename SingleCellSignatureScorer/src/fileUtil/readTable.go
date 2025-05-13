/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2019
*/

package fileUtil

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/schollz/progressbar"
)

func ReadTable(path string, database map[string][]string, nbCPU int, conf CONF) (map[idPW]float64, []string) {
	var allCellNames []string
	fileLength, err := lineCounter(path)
	check(err)
	allPW := make(map[idPW]float64, fileLength) // cellName + pathway -> []scores

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

	ch1 := make(chan map[idPW]float64, nbCPU)

	for count := 1; count < fileLength; count++ {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		cellName := line[0]
		geneValues := line[1:]
		allCellNames = append(allCellNames, cellName)
		go processLine(cellName, colNames, geneValues, database, ch1, conf) // process one line = one cell against all pathways
	}

	bar := progressbar.New(fileLength - 1) // Add a new progress bar

	for i := 1; i < fileLength; i++ {
		if conf.Server == 0 {
			bar.Add(1) // show progress bar
		}
		msg := <-ch1
		// merge all maps, each map = one pwname + one cellname -> score
		for k, v := range msg {
			allPW[k] = v
		}
	}
	close(ch1) // close channel
	return allPW, allCellNames
}
