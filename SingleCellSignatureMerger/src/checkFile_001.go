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
	"errors"
	"strconv"
)

func checkTSNE(tsneHeader []string) {
	lenTsne := len(tsneHeader)

	// check nb of columns
	if lenTsne < 3 {
		msg := "error in t-SNE coordinates ! the number of columns found in t-SNE table is " + strconv.Itoa(lenTsne) + " \nThe separator must be TAB and the number of columns must be at least 3"
		err := errors.New(msg)
		check(err)
	}
	// check column names
	// if contains(tsneHeader, "tSNE_1") == false || contains(tsneHeader, "tSNE_2") == false {
	// 	msg := "error in t-SNE coordinates ! t-SNE table must have columns labeled  t§NE_1 and t§NE_2 !"
	// 	err := errors.New(msg)
	// 	check(err)
	// }
}

func checkScoreHeader(tableHeader []string, pathway string) {
	lenTsne := len(tableHeader)

	// check nb of columns
	if lenTsne < 2 {
		msg := "error in score Table ! the number of columns found is " + strconv.Itoa(lenTsne) + " in " + pathway + "\nThe separator must be TAB and the number of columns must be at least 2"
		err := errors.New(msg)
		check(err)
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
