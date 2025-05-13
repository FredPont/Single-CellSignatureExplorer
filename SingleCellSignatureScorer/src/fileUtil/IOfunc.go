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
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// idPW used as key in map idPW -> pathway score
type idPW struct {
	cellN string // cell name
	pwN   string // pathway name
}

//###########################################

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//###########################################
// 			I/O functions
//###########################################

func ListFiles(path string) []string {
	var fileSlice []string
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		file := f.Name()
		fileSlice = append(fileSlice, file)
	}
	return fileSlice
}

// count file line
// from https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently

func lineCounter(inPath string) (int, error) {
	//inPath := inDir + "/" + file
	csvfile, err := os.Open(inPath) //open file
	check(err)
	// make a read buffer
	r := bufio.NewReader(csvfile)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// remove file extension
func remExt(filename string) (string, string) {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name, extension
}

// remove file extension in slice of filenames
func remExtFslice(s []string) []string {
	var noExtS []string
	for _, n := range s {
		name, _ := remExt(n)
		noExtS = append(noExtS, name)
	}
	return noExtS
}

// read one pathway into a slice of strings
func slurpFile(path string) []string {
	file, err := os.ReadFile(path)
	check(err)
	s := strings.Split(string(file), "\n")
	s = s[:len(s)-1] // remove last empty element of slice
	return s
}

// read one pathway into a slice of strings
func readPW(path string) []string {
	var genes []string
	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.FieldsPerRecord = 1
	for {

		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		genes = append(genes, record[0])
	}
	return genes
}

// read one database into map pathway name -> genes
func ReadDB(path string) map[string][]string {
	dataBase := make(map[string][]string, 0) // pathway name -> genes
	// read DB files
	files, _ := os.ReadDir(path)

	for _, f := range files {
		file := f.Name()
		pwSlice := readPW(path + file)
		dataBase[file] = pwSlice
	}
	return dataBase
}

//###########################################

// remove duplicates in []string
func uniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// slice intersection that remove "-" temporarly the signe - in the database genenames
// caution it is not a symetrical function, the sign - MUST be in the slice b, NEVER in a
func intersectionNeg(a, b []string) (c []string) {
	m := make(map[string]bool, len(a))
	negSign := false

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {

		if item[:1] == "-" { // if sign "-" is detected , the genes is append in the intersection list with a sign "-"
			item = item[1:]
			negSign = true
		}
		if _, ok := m[item]; ok {
			if negSign == true {
				c = append(c, "-"+item)
			} else {
				c = append(c, item)
			}
		}
		negSign = false
	}
	return
}

// return all pathway names in map[idPW]
func mapKeysPWid(mymap map[idPW]float64) []string {
	keys := make([]string, len(mymap))
	i := 0
	for k := range mymap {
		keys[i] = k.pwN
		i++
	}
	keys = uniqueStrings(keys)
	sort.Strings(keys)
	return keys
}

// sum of floats in slice
func sliceSum(s []float64) float64 {
	sum := 0.
	for _, x := range s {
		sum += x
	}
	return sum

}

// []float64 -> []string
func floatStoStringS(a []float64) []string {
	if len(a) == 0 {
		return []string{}
	}
	b := make([]string, len(a))
	for i, v := range a {
		val := strconv.FormatFloat(v, 'f', 3, 64)
		b[i] = val
	}
	return b
}

// remove genes with expression = 0
func cleanZero(colNames, geneValues []string) ([]string, []float64) {
	var geneNames []string
	var genesExpress []float64

	for i, s := range geneValues {
		v, err := strconv.ParseFloat(s, 64)
		check(err)
		if v > 0 {
			geneNames = append(geneNames, string(colNames[i]))
			genesExpress = append(genesExpress, v)
		}
	}

	return geneNames, genesExpress
}

// remove log2 transformation
func unLog2(log2 []float64) []float64 {
	var noLog2 []float64
	for _, i := range log2 {
		noLog2 = append(noLog2, math.Pow(2, i))
	}
	return noLog2
}
