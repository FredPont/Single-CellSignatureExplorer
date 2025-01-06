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
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar"
)

func header() {

	fmt.Println("   ┌───────────────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │       scSignature Merger (c)Frederic PONT 2018    │")
	fmt.Println("   │       Free Software GNU General Public License    │")
	fmt.Println("   └───────────────────────────────────────────────────┘")
}

func main() {
	isList := false // is there lists of genes to extract in table detected in the lists directory
	header()
	//t0 := time.Now()
	nbCPU := runtime.NumCPU()
	fmt.Println("   ", nbCPU, "CPUs detected")

	dirPw := ListFiles("scores")   // list files in scores dir (pathways)
	dirTsne := ListFiles("tsne")   // list files in t-sne dir - only one file is allowed
	dirLists := ListFiles("lists") // list files with genes lists
	var listsArray [][]string      // arrays of array of genes ex: [[CD4 CD5 CD8] [CD56 CD69 isotype_Ctrl]]
	var listsNames []string        // lists files names without extention
	var sort bool                  // if true sort table colnames

	if len(dirTsne) > 1 {
		msg := "error ! more than one file found in t-sne" + strings.Join(dirTsne, " ")
		err := errors.New(msg)
		check(err)
	}
	if len(dirLists) > 0 {
		listsArray, listsNames = readAllLists(dirLists)
		isList = true
	} else {
		sort = askuser()
	}
	t0 := time.Now()
	tsneDict, tsneHeader := IndexTsne(dirTsne[0]) // read t-sne file and index it in a map cell name => XY
	checkTSNE(tsneHeader)                         // check table length and column names

	if isList == false {
		if sort == false {
			fastMerge(tsneDict, dirPw, tsneHeader, nbCPU)
		} else {
			sortMerge(tsneDict, dirPw, tsneHeader, nbCPU)
		}
	} else {
		listMerge(tsneDict, dirPw, tsneHeader, dirLists, listsNames, nbCPU, listsArray)
	}

	fmt.Println("Finished !")
	fmt.Printf("Elapsed time : %v.\n", time.Since(t0))
	fmt.Print("Press enter to close window ")
	fmt.Scanln() // saisie clavier

}

func fastMerge(tsneDict map[string][]string, dirPw, tsneHeader []string, nbCPU int) {
	ch1 := make(chan string, 2*nbCPU+1)
	for _, pwFile := range dirPw {
		go ConcatFiles(pwFile, tsneDict, tsneHeader, ch1)
	}

	bar := progressbar.New(len(dirPw)) // Add a new progress bar

	for i := 0; i < len(dirPw); i++ {
		bar.Add(1) // show progress bar
		msg := <-ch1
		fmt.Println(msg, "merged !")
	}
}

func listMerge(tsneDict map[string][]string, dirPw, tsneHeader, dirLists, listsNames []string, nbCPU int, listsArray [][]string) {
	ch1 := make(chan string, 2*nbCPU+1)

	for _, pwFile := range dirPw {
		for i, list := range listsArray {
			go mergeSubTable(pwFile, listsNames[i], tsneDict, tsneHeader, ch1, list)
		}
	}

	bar := progressbar.New(len(dirPw) * len(listsNames)) // Add a new progress bar

	for i := 0; i < len(dirPw)*len(listsNames); i++ {
		bar.Add(1) // show progress bar
		msg := <-ch1
		fmt.Println(msg, "merged !")

	}

}

func sortMerge(tsneDict map[string][]string, dirPw, tsneHeader []string, nbCPU int) {
	ch1 := make(chan string, 2*nbCPU+1)

	for _, pwFile := range dirPw {
		go sortedMerge(pwFile, tsneDict, tsneHeader, ch1)
	}

	bar := progressbar.New(len(dirPw)) // Add a new progress bar

	for i := 0; i < len(dirPw); i++ {
		bar.Add(1) // show progress bar
		msg := <-ch1
		fmt.Println(msg, "merged !")
	}

}

// ask user for sorted table
func askuser() bool {

	sort := false
	fmt.Print("Sort table y/n (n) ? (usefull for genes, useless for signatures) : ")

	// lecture du clavier
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}

	if s == "q" {
		os.Exit(1)
	}
	if s != "" { // si s not nul then input = s
		if s == "y" {
			sort = true
		}
	}
	return sort
}
