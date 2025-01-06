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
 (c) Frederic Pont 2018
*/

package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"sort"
)

//###########################################

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// read header of table
func readHeader(path string) []string {

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	record, err := reader.Read() // read first line

	return record
}

// IndexTsne create a map with the first cell of each row => row
// cell name => XY
func IndexTsne(tsneFile string) (map[string][]string, []string) {
	tsneDict := make(map[string][]string)

	path := "tsne/" + tsneFile

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	tsneHeader, err := reader.Read() // skip first line

	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		tsneDict[record[0]] = record[1:] //remove first cell containing cell name in record
	}
	return tsneDict, tsneHeader
}

// read one list file and return []string containing list elements
func readList(fileName string) []string {
	var list []string

	csvFile, err := os.Open("lists/" + fileName)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = 1

	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		list = append(list, record[0]) //remove first cell containing cell name in record
	}

	return list
}

func readAllLists(dirLists []string) ([][]string, []string) {
	var listsArray [][]string
	var listsNames []string
	for _, fn := range dirLists {
		list := readList(fn)
		listsArray = append(listsArray, list)
		name, _ := remExt(fn)
		listsNames = append(listsNames, name)
	}
	return listsArray, listsNames
}

// remove empty string in []string
func remEmpty(s []string) []string {
	var clean []string

	for _, l := range s {
		if l != "" {
			clean = append(clean, l)
		}
	}
	sort.Strings(clean) // sort list items
	return clean
}
