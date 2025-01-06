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

package pogrebdb

import (
	"fmt"
	"log"

	"github.com/akrylysov/pogreb"
)

// InsertColDB inserts one column in the database.
// Pogreb is thread safe according to https://github.com/akrylysov/pogreb#key-characteristics
func InsertColDB(db *pogreb.DB, key, column []byte) {
	// mutex.Lock()
	// defer mutex.Unlock()
	err := db.Put([]byte(key), column)
	if err != nil {
		log.Fatal(err)
	}

}

// CreateDataBase create a new pogreb database in tmp dir
func CreateDataBase(dbName string) *pogreb.DB {
	db, err := pogreb.Open("tmp/"+dbName, nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println("cannot create database in tmp dir !")
	}
	//defer db.Close() // do not close the database before the end of the run

	fmt.Println("database created in tmp/" + dbName)
	return db
}
