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
	"math"
	"strconv"
)

//###########################################

// slice intersection that remove "-" temporarly the signe - in the database genenames
// caution it is not a symetrical function, the sign - MUST be in the slice b, NEVER in a
func IntersectionNeg(a, b []string) (c []string) {
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

// remove genes with expression = 0
func CleanZero(colNames, geneValues []string) ([]string, []float64) {
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
func UnLog2(log2 []float64) []float64 {
	var noLog2 []float64
	for _, i := range log2 {
		noLog2 = append(noLog2, math.Pow(2, i))
	}
	return noLog2
}
