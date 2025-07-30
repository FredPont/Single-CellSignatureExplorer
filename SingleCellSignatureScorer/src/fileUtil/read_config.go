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

package fileUtil

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

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

func ServerDB(conf CONF, allDBnames []string) []string {
	var selectedDB []string
	input := parseNB(conf.DBserver)
	for _, dbIndex := range input {
		selectedDB = append(selectedDB, allDBnames[dbIndex])
	}
	return selectedDB
}

func InitializeConfig() CONF {
	return CONF{
		RemLog2:         1,
		Server:          0,
		DBserver:        "0",
		ScoreAlgo:       "sum",
		ColNames:        []string{},
		GeneIndex:       make(map[string]int),
		NbCols:          0,
		ChunkSize:       1000,
		ChunkInParallel: 2,
		NbCPU:           1,
	}
}

// SWitch to server mode
func ApplyCommandLineParams(conf *CONF) {
	// read cmb line arguments
	cmdlineDB := flag.String("db", "", "databases number selected by user in command line : for example 0-4,6-12")
	flag.Parse()

	if conf.Server == 0 {
		Header()
	} else if *cmdlineDB != "" {
		conf.DBserver = *cmdlineDB // set the DBserver parameter by default
	}
}

// Display databases to select
func SelectDatabases(conf CONF, allDBnames []string) []string {
	var DBnames []string
	if conf.Server == 0 {
		DisplayDB(allDBnames)          // show databases to select
		DBnames = Criteria(allDBnames) // user select some databases
	} else {
		DBnames = ServerDB(conf, allDBnames)
	}
	return DBnames
}

// ConfigSummary prints the configuration
func SimpleConfigSummary(conf CONF) {
	fmt.Println("Configuration :")
	fmt.Println("   Remove log2 :", conf.RemLog2)
	fmt.Println("   Server mode :", conf.Server)
	fmt.Println("   DBserver :", conf.DBserver)
	fmt.Println("   Score algorithm :", conf.ScoreAlgo)
	fmt.Println("   Chunk size :", conf.ChunkSize)
	fmt.Println("   Chunk in parallel :", conf.ChunkInParallel)
	fmt.Println("   Nb CPU :", conf.NbCPU)

}

// ConfigSummary prints the configuration with a table alignement
func ConfigSummary(conf CONF) {
	// Create a new tab writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Configuration:\t\t\t")
	fmt.Fprintln(w, "Remove log2:\t", conf.RemLog2)
	fmt.Fprintln(w, "Server mode:\t", conf.Server)
	fmt.Fprintln(w, "DBserver:\t", conf.DBserver)
	fmt.Fprintln(w, "Score algorithm:\t", conf.ScoreAlgo)
	fmt.Fprintln(w, "Chunk size:\t", conf.ChunkSize)
	fmt.Fprintln(w, "Chunk in parallel:\t", conf.ChunkInParallel)
	fmt.Fprintln(w, "Nb CPU:\t", conf.NbCPU)

	// Flush the writer to output the formatted table
	w.Flush()
}
