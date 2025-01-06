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

package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar"
)

// idPW used as key in map idPW -> pathway score
type idPW struct {
	cellN string // cell name
	pwN   string // pathway name
}

// ###########################################
func header() {
	fmt.Println("")
	fmt.Println("   ┌─────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │ Single-Cell Scorer (c)Frederic PONT │")
	fmt.Println("   │ 2018-2020 - Free Software GNU GPL   │")
	fmt.Println("   └─────────────────────────────────────┘")
	//fmt.Println("")
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

func listFiles(path string) []string {
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
	file, err := ioutil.ReadFile(path)
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
func readDB(path string) map[string][]string {
	dataBase := make(map[string][]string, 0) // pathway name -> genes
	// read DB files
	files, _ := ioutil.ReadDir(path)

	for _, f := range files {
		file := f.Name()
		pwSlice := readPW(path + file)
		dataBase[file] = pwSlice
	}
	return dataBase
}

func readTable(path string, database map[string][]string, nbCPU int, conf CONF) (map[idPW]float64, []string) {
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

//###########################################
// process number list entered by user
// to choose database

func parseNB(l string) []int {
	var list []int

	switch strings.Contains(l, ",") {
	case true:
		l2 := strings.Split(l, ",")
		for _, l1 := range l2 {
			if strings.Contains(l1, "-") {
				l3 := parseRG(l1)
				list = append(list, l3...)
			} else {
				nb, err := strconv.Atoi(l1)
				if err != nil {
					fmt.Println("Impossible to read number : ", l1, err.Error())
				}
				list = append(list, nb)
			}
		}
	case false:
		switch strings.Contains(l, "-") {
		case true:
			l3 := parseRG(l)
			list = append(list, l3...)
		default:
			nb, err := strconv.Atoi(l)
			if err != nil {
				fmt.Println("Impossible to read number : ", l, err.Error())
			}
			list = append(list, nb) // l ne contient ni "," , ni "-", ie un seul nombre
		}
	}
	return list
}

// parse range nb 1-5
func parseRG(r string) []int {

	var listNB []int
	rg := strings.Split(r, "-")

	start, err := strconv.Atoi(rg[0])
	if err != nil {
		fmt.Println("Impossible to read number : ", rg[0], err.Error())
	}
	end, err := strconv.Atoi(rg[1])
	if err != nil {
		fmt.Println("Impossible to read number : ", rg[1], err.Error())
	}

	for i := start; i <= end; i++ {
		listNB = append(listNB, i)
	}
	return listNB
}

// ###########################################
func displayDB(DBnames []string) {

	fmt.Println("")
	fmt.Println("   ┌─────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │           DataBases list            │")
	fmt.Println("   ├─────────────────────────────────────┤")
	fmt.Println("   │                                     │")
	for i, d := range DBnames {
		if i < 10 {
			fmt.Printf("   │ %-d  - %-31s│\n", i, d) // insert space for better alignment :)
		} else {
			fmt.Printf("   │ %-d - %-31s│\n", i, d)
		}
	}
	fmt.Println("   └─────────────────────────────────────┘")
	fmt.Println("   select working databases, default = 0")
	fmt.Println("   example : 1-3,5")
}

// ask user for database choice
func criteria(allDBnames []string) []string {
	var selectedDB []string
	input := []int{0} // input

	fmt.Print("   database number : (q quit): ")

	// lecture du clavier
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	fmt.Println("   ____________________________________")
	fmt.Println()

	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}

	if s == "q" {
		os.Exit(1)
	}
	if s != "" { // si s not nul then input = s
		input = parseNB(s)

	}
	if len(input) > 1 {
		fmt.Println("   databases ", input, "are selected")
	} else {
		fmt.Println("   database ", input, "is selected")
	}
	fmt.Println("   ____________________________________")
	fmt.Println()

	for _, dbIndex := range input {
		selectedDB = append(selectedDB, allDBnames[dbIndex])
	}

	return selectedDB
}

// ###########################################################
// process one data line = one cell
func processLine(cellName string, colNames, geneValues []string, database map[string][]string, ch1 chan<- map[idPW]float64, conf CONF) {
	geneNames, genesExpress := cleanZero(colNames, geneValues) // remove genes with null expression
	var geneExpressNoLog []float64

	// remove log 2 transformation if needed
	if conf.RemLog2 == 1 {
		geneExpressNoLog = unLog2(genesExpress) // remove log2 for all genes
	} else {
		geneExpressNoLog = genesExpress
	}

	sumUMI := sliceSum(geneExpressNoLog)                 // UMI sum of all genes
	GNnormUMI := GeneExpDic(geneExpressNoLog, geneNames) // gene name -> gene express no log2 no zero

	tmpPWvalues := make(map[idPW]float64) // [cell name ; pw name] -> UMI
	// for each pathway in database
	for pwName, pwGenes := range database {

		UMISum := 0.0
		if len(geneNames) > 0 {
			genesInPW := intersectionNeg(geneNames, pwGenes)
			if len(genesInPW) > 0 {
				for _, g := range genesInPW {
					if g[:1] == "-" { //if a sign "-" is detected then UMI of the gene without "-" g[1:]) is substracted
						UMISum = UMISum - GNnormUMI[g[1:]]
					} else {
						UMISum = UMISum + GNnormUMI[g]
					}
				}
			}
			tmpPWvalues[idPW{cellName, pwName}] = UMISum * 100. / sumUMI // the score is the UMI %

		} else {
			tmpPWvalues[idPW{cellName, pwName}] = 0.0
		}
	}
	ch1 <- tmpPWvalues

}

// GeneExpDic create dict : gene name -> gene express no log2 no zero
func GeneExpDic(geneExpressNoLog []float64, geneNames []string) map[string]float64 {

	GNnormUMI := make(map[string]float64, len(geneNames)) // gene name -> gene express no log2 no zero
	for i, v := range geneNames {
		GNnormUMI[v] = geneExpressNoLog[i]
	}
	return GNnormUMI
}

// remove pathways with zero score in all cells
func removeNullPW(allCellNames []string, allPW map[idPW]float64) []string {
	allPWnames := mapKeysPWid(allPW)
	var nonNullPWnames []string

	for _, pwName := range allPWnames {
		var pathwayScores []float64
		for _, cellName := range allCellNames {
			pathwayScores = append(pathwayScores, allPW[idPW{cellName, pwName}])
		}
		if sliceSum(pathwayScores) != 0. {
			nonNullPWnames = append(nonNullPWnames, pwName)
		}
	}
	return nonNullPWnames
}

func writeResults(path string, allCellNames []string, allPW map[idPW]float64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	nonNullPWnames := removeNullPW(allCellNames, allPW) // remove PW with null scores for all cells

	fmt.Fprintln(w, "id"+"\t"+strings.Join(remExtFslice(nonNullPWnames), "\t"))

	var cellScores []float64 // cell score for all pathways > 0

	for _, c := range allCellNames {
		for _, pwName := range nonNullPWnames {
			cellScores = append(cellScores, allPW[idPW{c, pwName}]) // merge all scores for one cell
		}
		line := c + "\t" + strings.Join(floatStoStringS(cellScores), "\t")
		fmt.Fprintln(w, line)
		cellScores = nil
	}

	return w.Flush()
}

func main() {
	// read cmb line arguments
	cmdlineDB := flag.String("db", "", "databases number selected by user in command line : for example 0-4,6-12")
	flag.Parse()

	conf := ReadConfig() // read json config file
	if conf.Server == 0 {
		header()
	} else if *cmdlineDB != "" {
		conf.DBserver = *cmdlineDB // priority is given to cmd line database to json file databases
	}

	//
	allDBnames := listFiles("databases/")
	dataFileNames := listFiles("data/")
	nbCPU := runtime.NumCPU()
	fmt.Println("   ", nbCPU, "CPUs detected")

	var DBnames []string
	if conf.Server == 0 {
		displayDB(allDBnames)          // display databases
		DBnames = criteria(allDBnames) // user select databases
	} else {
		DBnames = serverDB(conf, allDBnames)
	}

	totalCalc := len(DBnames) * len(dataFileNames)
	count := 0
	t0 := time.Now()

	// for each database
	for _, db := range DBnames {
		dataBase := readDB("databases/" + db + "/")
		// for each data file
		for _, dataFile := range dataFileNames {
			count++
			fmt.Println("\n", count, "/", totalCalc, " File: ", dataFile, " DB: ", db)

			//allPW := make(map[idPW]float64, 0) // cell name + pathway -> []scores
			allPW, allCellNames := readTable("data/"+dataFile, dataBase, nbCPU, conf)

			resFile := "results/" + db + "_" + dataFile
			fmt.Println("\n        Write:", resFile)
			writeResults(resFile, allCellNames, allPW)
		}
	}
	fmt.Println("Finished !")
	fmt.Printf("Elapsed time : %v.\n", time.Since(t0))

}
