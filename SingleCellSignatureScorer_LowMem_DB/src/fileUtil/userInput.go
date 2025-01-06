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
	"fmt"
	"os"
	"strconv"
	"strings"
)

//###########################################
// process number list entered by user
// to choose database

func ParseNB(l string) []int {
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

// ask user for database choice
func Criteria(allDBnames []string) []string {
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
		input = ParseNB(s)

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
