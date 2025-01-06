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
	"strings"
)

// detect columns to select
// test :
// header := []string{"a", "b", "c", "d", "e", "f"}
// list := []string{"e", "b", "z", "c"}
// result : [4 1 2]

func getColIndex(header, list []string) []int {
	var indexes []int
	indDic := make(map[string]int) // dic of genes -> column index
	list2 := make([]string, len(list))
	copy(list2, list)

	for i, val := range header {
		for j, l := range list2 {
			if val == l {
				indDic[val] = i
				list2 = append(list2[:j], list2[j+1:]...) // remove found item from list
				break
			}
		}
	}

	indexes = append(indexes, 0) // append the first column containing cells names

	for _, v := range list {
		value, exist := indDic[v]
		if exist {
			indexes = append(indexes, value)
		}
	}
	return indexes
}

// selByIndex select item in a slice accordiing to indexes
// we use it to select in the header of score table only columns
// corresponding to indexes positions
func selByIndex(header []string, indexes []int) []string {
	var selection []string

	for _, i := range indexes {
		selection = append(selection, header[i])
	}
	return selection
}

func mergeSubTable(pathway, listsName string, tsneDict map[string][]string, tsneHeader []string, ch1 chan<- string, list []string) {
	path := "scores/" + pathway
	header := readHeader(path)            // read signature/gene table header
	colIndex := getColIndex(header, list) // index of columns in header corresponding to list items

	// open result file for write
	baseNamePW, _ := remExt(pathway)
	fout := "results/" + baseNamePW + "_" + listsName + ".tsv"
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
	pathwayHeader = selByIndex(pathwayHeader, colIndex)
	mergedHeader := append(pathwayHeader, tsneHeader[1:]...)

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
	ch1 <- baseNamePW + "_" + listsName
}
