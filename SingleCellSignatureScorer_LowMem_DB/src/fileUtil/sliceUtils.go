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
	"ScorerLowMem/src/userTypes"
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/akrylysov/pogreb"
)

// sum of floats in slice
func SliceSum(s []float64) float64 {
	sum := 0.
	for _, x := range s {
		sum += x
	}
	return sum

}

// float64 -> string
func FloatString(a float64) string {
	s := fmt.Sprintf("%.3f", a)
	return (s) //

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

// remove duplicates in []string
func UniqueStrings(input []string) []string {
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

// return all pathway names in map[idPW]
// func mapKeysPWid(mymap map[idPW]float64) []string {
// 	keys := make([]string, len(mymap))
// 	i := 0
// 	for k := range mymap {
// 		keys[i] = k.pwN
// 		i++
// 	}
// 	keys = uniqueStrings(keys)
// 	sort.Strings(keys)
// 	return keys
// }

// sort slice of PWscore

// BypwN implements sort.Interface based on the Age field.
type BypwN []userTypes.PWscore

func (a BypwN) Len() int           { return len(a) }
func (a BypwN) Less(i, j int) bool { return a[i].PwN < a[j].PwN }
func (a BypwN) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func SortPWscore(s []userTypes.PWscore) {
	sort.Sort(BypwN(s))
}

func FloatToByte(f float64) []byte {

	bits := math.Float64bits(f)
	byteSlice := make([]byte, 8) // float64 is 8 bytes
	binary.LittleEndian.PutUint64(byteSlice, bits)

	return byteSlice
}

// FloatSliceToByteSlice convert []float64 to []byte
func FloatSliceToByteSlice(floatArray []float64) []byte {

	// Convert []float64 to []byte
	byteArray := make([]byte, len(floatArray)*8)
	for i, f := range floatArray {
		binary.LittleEndian.PutUint64(byteArray[i*8:], math.Float64bits(f))
	}
	return byteArray
}

// ByteSliceToStringslice convert []byte to []float64
func ByteSliceToStringslice(byteArray []byte) []float64 {

	// Convert []byte back to []float64
	FloatArray := make([]float64, len(byteArray)/8)
	for i := 0; i < len(byteArray); i += 8 {
		bits := binary.LittleEndian.Uint64(byteArray[i : i+8])
		FloatArray[i/8] = math.Float64frombits(bits)
	}
	return FloatArray
}

// ByteSliceToStringslice convert []float64 to []string
func FloatSliceToStringslice(FloatArray []float64) []string {

	// Convert []float64 to []string
	stringArray := make([]string, len(FloatArray))
	for i, f := range FloatArray {
		stringArray[i] = strconv.FormatFloat(f, 'f', -1, 64)
	}
	return stringArray
}

func ByteSliceToRow(byteArray []byte) string {
	// Convert []byte back to []float64
	FloatArray := ByteSliceToStringslice(byteArray)
	// Convert []float64 to []string
	stringArray := FloatSliceToStringslice(FloatArray)
	// Convert []string to string
	row := ""
	for _, s := range stringArray {
		row += s + "	"
	}
	return strings.TrimSpace(row) // remove last tabulation
}

// DBtoCSV export the database to a csv result file
func DBtoCSV(db *pogreb.DB, allPWnames []string, resFile string) error {
	path := resFile
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("write", path)

	w := bufio.NewWriter(file)
	header := "id" + "\t" + strings.Join(allPWnames, "\t")
	fmt.Fprintln(w, header)

	it := db.Items()
	for {
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		row := ByteSliceToRow(val)
		row = string(key) + "\t" + row
		//fmt.Println(row)
		fmt.Fprintln(w, row)
		//log.Printf("%s %s", key, row)
	}
	w.Flush()
	return nil
}

// SortMap sort the database map to have always the same order of pathways
func SortKeysMap(m map[string][]string) []string {

	// Extract keys from map
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Strings(keys)

	// Print sorted map
	// for _, k := range keys {
	// 	fmt.Printf("%s: %d\n", k, m[k])
	// }

	return keys
}
