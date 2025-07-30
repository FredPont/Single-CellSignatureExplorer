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
 (c) Frederic Pont 2025
*/

package fileUtil

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

// MetaDataSummary prints the metadata in a json format
func MetaDataSummary(db, dataFile string) {
	// Create a new JSON file
	// The file name is the result file name with the extension .json
	resultFileWithoutExt, _ := remExt(Config.ResultFileName)
	jsonFileName := resultFileWithoutExt + ".json"
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer jsonFile.Close()

	metadata := map[string]interface{}{
		"AnalysisDate":    time.Now().Format("2006-01-02"),
		"AnalysisTime":    time.Now().Format("15:04:05"),
		"SoftwareVersion": softVersion(),
		"OperatingSystem": runtime.GOOS,
		"CPUArch":         runtime.GOARCH,
		"ComputerCPU":     runtime.NumCPU(),
		"Database":        db,
		"DataFile":        dataFile,
		"RemoveLog2":      Config.RemLog2,
		"Server":          Config.Server,
		"DBserver":        Config.DBserver,
		"ScoreAlgo":       Config.ScoreAlgo,
		"ChunkSize":       Config.ChunkSize,
		"ChunkInParallel": Config.ChunkInParallel,
		"NbCPU":           Config.NbCPU,
		"ResultFileName":  Config.ResultFileName,
		"User":            os.Getenv("USER"), // user name
		"GoVersion":       runtime.Version(), // Version of Go
	}

	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(metadata)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
