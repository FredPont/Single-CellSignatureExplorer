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

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// CONF stores software parameters
// always use maj for conf variables
type CONF struct {
	RemLog2  int    `json:"removeLog2"`
	Server   int    `json:"server"`
	DBserver string `json:"DBserver"`
}

// ReadConfig reads conf.json
func ReadConfig() CONF {
	file, err1 := os.Open("conf.json")
	if err1 != nil {
		fmt.Println(err1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := CONF{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(conf)
	return conf
}

func serverDB(conf CONF, allDBnames []string) []string {
	var selectedDB []string
	input := parseNB(conf.DBserver)
	for _, dbIndex := range input {
		selectedDB = append(selectedDB, allDBnames[dbIndex])
	}
	return selectedDB
}
