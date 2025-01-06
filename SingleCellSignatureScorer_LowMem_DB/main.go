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
	"ScorerLowMem/src/cpu"
	"ScorerLowMem/src/fileUtil"
	"ScorerLowMem/src/pogrebdb"
	"flag"
	"fmt"
	"runtime"
	"time"
)

func main() {

	// clean tmp dir
	fileUtil.RemoveContents("tmp/")
	// read cmb line arguments
	cmdlineDB := flag.String("db", "", "databases number selected by user in command line : for example 0-4,6-12")
	flag.Parse()

	conf := fileUtil.ReadConfig() // read json config file
	if conf.Server == 0 {
		fileUtil.Header()
	} else if *cmdlineDB != "" {
		conf.DBserver = *cmdlineDB // priority is given to cmd line database to json file databases
	}

	allDBnames := fileUtil.ListFiles("databases/")
	dataFileNames := fileUtil.ListFiles("data/")
	nbCPU := cpu.CalcCPU(conf.CPU)
	fmt.Println("   ", nbCPU, "CPUs detected")
	runtime.GOMAXPROCS(nbCPU) // limit the number of CPUs

	var DBnames []string
	if conf.Server == 0 {
		fileUtil.DisplayDB(allDBnames)          // display databases
		DBnames = fileUtil.Criteria(allDBnames) // user select databases
	} else {
		DBnames = fileUtil.ServerDB(conf, allDBnames)
	}

	totalCalc := len(DBnames) * len(dataFileNames)
	count := 0
	t0 := time.Now()

	// for each database
	for _, db := range DBnames {
		dataBase, allPWnames := fileUtil.ReadDB("databases/" + db + "/")
		// for each data file
		for _, dataFile := range dataFileNames {
			count++
			fmt.Println("\n", count, "/", totalCalc, " File: ", dataFile, " DB: ", db)

			fileUtil.UpdateConf(dataFile, db, &conf) // add current datafile and database to conf

			// create a new Pogreb database and return the temp files database
			fileDB := pogrebdb.CreateDataBase(db + "_" + fileUtil.BaseFile(dataFile))

			fileUtil.ReadTable(fileDB, "data/"+dataFile, dataBase, nbCPU, conf)

			resFile := "results/" + db + "_" + dataFile

			fileUtil.DBtoCSV(fileDB, allPWnames, resFile)
			fmt.Printf("Elapsed time : %v.\n", time.Since(t0))

			fileDB.Close()
		}
	}
	// clean tmp dir
	fileUtil.RemoveContents("tmp/")
	fmt.Println("All files analyzed !")
	fmt.Printf("Elapsed time : %v.\n", time.Since(t0))

}
