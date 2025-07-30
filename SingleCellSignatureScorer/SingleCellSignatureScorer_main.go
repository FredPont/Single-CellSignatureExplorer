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
	"fmt"
	_ "net/http/pprof"
	"scorer/src/fileUtil"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {

	// start the profiling server
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	fileUtil.Config = fileUtil.InitializeConfig()

	// read configuration from JSON file
	fileUtil.Config = fileUtil.ReadConfig()

	// Apply command line parameters, switch between server and client mode
	fileUtil.ApplyCommandLineParams(&fileUtil.Config)

	// read all databases names
	allDBnames := fileUtil.ListFiles("databases/")
	// read all data files names
	dataFileNames := fileUtil.ListFiles("data/")

	// Display configuration
	fileUtil.ConfigSummary(fileUtil.Config)

	// Select databases to use
	DBnames := fileUtil.SelectDatabases(fileUtil.Config, allDBnames)

	// Display databases to use
	fmt.Println("Selected Databases:", DBnames)
	totalCalc := len(DBnames) * len(dataFileNames)
	count := 0
	t0 := time.Now()

	// for each database
	bar := progressbar.New(len(DBnames) * len(dataFileNames))
	for _, db := range DBnames {
		fileUtil.Config.DataBase = fileUtil.ReadDB("databases/" + db + "/")
		// for each data file
		for _, dataFile := range dataFileNames {

			count++
			fmt.Println("\n", count, "/", totalCalc, " File: ", dataFile, " DB: ", db)

			fileUtil.CreateResultFile(db, dataFile) // create the result file

			fileUtil.ReadMatrix("data/" + dataFile) // read the data matrix

			fileUtil.MetaDataSummary(db, dataFile) // write the metadata in a json file
			bar.Add(1)                             // show progress bar
		}
		fileUtil.CloseResultFile()
	}
	fmt.Println("Finished !")
	fmt.Printf("Elapsed time : %v.\n", time.Since(t0))

}
