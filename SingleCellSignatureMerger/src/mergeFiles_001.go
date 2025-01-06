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

// merge t-sne coordinate with score tables
package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// remove file extension
func remExt(filename string) (string, string) {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name, extension
}

// ListFiles lists all files in a directory
func ListFiles(dir string) []string {
	var filesList []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		filesList = append(filesList, f.Name())
	}
	return filesList
}

// ConcatFiles merge 2 tables according to the first column
func ConcatFiles(pathway string, tsneDict map[string][]string, tsneHeader []string, ch1 chan<- string) {
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
			line := append(record, value...)
			row := strings.Join(line, "\t") + "\n"
			writeOneLine(out, row)
		}
	}
	ch1 <- pathway
}

//###########################################
func writeOneLine(f *os.File, line string) {
	_, err := f.WriteString(line)
	check(err)
}
