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
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"sort"

	"os"
	"path/filepath"
)

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
	files, _ := os.ReadDir(path)
	for _, f := range files {
		file := f.Name()
		fileSlice = append(fileSlice, file)
	}
	return fileSlice
}

// count file line
// from https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently

func LineCounter(inPath string) (int, error) {
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

// BaseFile remove file extension and return name without extension
func BaseFile(filename string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name
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

// read one database into map pathway name -> []genes
func ReadDB(path string) (map[string][]string, []string) {
	dataBase := make(map[string][]string, 0) // pathway name -> genes
	var allPWnames []string                  // all pathways files names without extension
	// read DB files
	files, _ := os.ReadDir(path)

	for _, f := range files {
		file := f.Name()
		name, _ := remExt(file)
		allPWnames = append(allPWnames, name)
		pwSlice := readPW(path + file)
		dataBase[file] = pwSlice
	}
	sort.Strings(allPWnames) // sort pw names
	return dataBase, allPWnames
}

// https://stackoverflow.com/questions/33450980/how-to-remove-all-contents-of-a-directory-using-golang
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// remove dir
func removeDir(dir string) {
	err := os.Remove(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// remove temp dir and content
func CleanTMP(dir string) {
	RemoveContents(dir)
	// removeDir(dir)
	// err := os.Mkdir(dir, 0755)
	// if err != nil {
	// 	log.Println(err)
	// }
}
