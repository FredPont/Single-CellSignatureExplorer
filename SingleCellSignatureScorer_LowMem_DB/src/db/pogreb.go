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
//(c) Frederic Pont 2023

package db

import (
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/akrylysov/pogreb"
)

func fakeData(db *pogreb.DB) {

	mutex := &sync.Mutex{}

	for i := 0; i < 2000; i++ {

		key := "key_" + strconv.Itoa(i)
		val := []string{strconv.Itoa(i), strconv.Itoa(i * i), strconv.FormatFloat(math.Pow(float64(i), 3), 'f', -1, 64)}
		go loadData(db, val, key, mutex)
	}

}

// loadData load a CSV file into the database
func loadData(db *pogreb.DB, values []string, key string, mutex *sync.Mutex) {

	//fmt.Println(column)
	str := []byte(strings.Join(values, "\t"))
	//insertCol(db, []byte(header[i]), str)
	mutex.Lock()
	defer mutex.Unlock()
	insertCol(db, []byte(key), str) //

}

// insertcol inserts one column in the database
func insertCol(db *pogreb.DB, key, column []byte) {

	err := db.Put([]byte(key), column)
	if err != nil {
		log.Fatal(err)
	}

}
