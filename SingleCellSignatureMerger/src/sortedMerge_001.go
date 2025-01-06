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

package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strings"
)

func sortHeader(header []string) ([]string, []int) {
	headerDic := make(map[string]int) // dic of colname => colindex
	indexes := make([]int, 1)         // append a zero in the slice = the first column containing cells names

	// leave column 0 containing cell names at position 0
	header2 := make([]string, len(header))
	copy(header2, header)
	header2 = header2[1:]

	for i, col := range header2 {
		headerDic[col] = i + 1
	}

	//indexes = append(indexes, 0) // append the first column containing cells names
	sort.Strings(header2) // sort columns names

	for _, v := range header2 {
		value, exist := headerDic[v]
		if exist {
			indexes = append(indexes, value)
		}
	}

	header2 = append([]string{"CellsId"}, header2...) // push front the colname for cells

	return header2, indexes

}

// ConcatFiles merge 2 tables according to the first column
func sortedMerge(pathway string, tsneDict map[string][]string, tsneHeader []string, ch1 chan<- string) {
	path := "scores/" + pathway

	// open result file for write
	baseNamePW, _ := remExt(pathway)
	fout := "results/" + baseNamePW + ".tsv"
	out, err1 := os.Create(fout)
	check(err1)
	defer out.Close()

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	// process tables headers
	pathwayHeader, err := reader.Read() //read first line of pathway table
	checkScoreHeader(pathwayHeader, pathway)
	sortedHeader, colIndex := sortHeader(pathwayHeader)
	mergedHeader := append(sortedHeader, tsneHeader[1:]...)

	writeOneLine(out, strings.Join(mergedHeader, "\t")+"\n") // write header in result file

	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		value, ok := tsneDict[record[0]] // test if record  exist in tsne
		if ok {
			record = selByIndex(record, colIndex) // extract records corresponding to colIndexe
			line := append(record, value...)
			row := strings.Join(line, "\t") + "\n"
			writeOneLine(out, row)
		}
	}
	ch1 <- pathway
}
