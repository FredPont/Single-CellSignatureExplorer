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
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// writeCSVFromChannel writes data from channel to CSV
func writeCSVFromChannel(dataChan <-chan []string, done chan bool, records int, w *csv.Writer) {
	// Write data from channel to CSV
	for data := range dataChan {
		// Write to CSV here
		// Read data from the channel and write to the file
		if err := w.Write(data); err != nil {
			log.Printf("Error writing line: %v", err)
			//return err
		}

		records--
		// Check if all records are processed, if yes then notify channel
		if records == 0 {
			done <- true
		}
	}
}

/*
func ListenAndServe(filename string, dataChan <-chan []string, wgP *sync.WaitGroup) {

	defer wgP.Done()
	for data := range dataChan {
		for {
			err := openAndWriteToFile(filename, strings.Join(data, "\t")+"\n")
			if err != nil {
				if err == syscall.EWOULDBLOCK {
					fmt.Println("File is currently locked, retrying...")
					time.Sleep(10 * time.Millisecond) // wait before retrying
					continue
				}

				//return
			}
			break
		}
		//openAndWriteToFile(filename, strings.Join(data, "\t")+"\n")
	}

}

func openAndWriteToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func appendToTSV(filename string, dataChan <-chan []string) error {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Iterate over the data and write it to the file
	for slice := range dataChan {

		for i, str := range slice {
			if i > 0 {
				_, err := file.WriteString("\t")
				if err != nil {
					return err
				}
			}
			_, err := file.WriteString(str)
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
		//mu.Unlock()
	}

	return nil
}

func WriteLine(filename string, slice []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	//mu.Lock()
	for i, str := range slice {
		if i > 0 {
			_, err := file.WriteString("\t")
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString(str)
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}
	//mu.Unlock()
	//file.Close()
	return nil
}

func WriteChanToFile(filename string, dataChan <-chan []string) error {
	// Ouvrir le fichier en écriture
	//file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// defer func() {
	// 	<-dataChan // free up a slot
	// }()

	// for slice := range dataChan {
	// 	fmt.Println("slice", slice)
	// }

	//Lire les données du canal et les écrire directement dans le fichier
	//mu.Lock()
	for slice := range dataChan {

		for i, str := range slice {
			if i > 0 {
				_, err := file.WriteString("\t")
				if err != nil {
					return err
				}
			}
			_, err := file.WriteString(str)
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
		//mu.Unlock()
	}
	//mu.Unlock()
	return nil

	// Créer un writer bufferisé
	// writer := bufio.NewWriter(file)
	// defer writer.Flush()

	// // Lire les données du canal et les écrire dans le fichier
	// for slice := range dataChan {
	// 	for i, str := range slice {
	// 		if i > 0 {
	// 			writer.WriteString("\t") // Séparateur entre les éléments de la slice
	// 		}
	// 		writer.WriteString(str)
	// 	}
	// 	writer.WriteString("\n") // Nouvelle ligne après chaque slice
	// }

	// return nil
}

func WriteResults(path string, allCellNames []string, allPW map[idPW]float64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	nonNullPWnames := removeNullPW(allCellNames, allPW) // remove PW with null scores for all cells

	fmt.Fprintln(w, "id"+"\t"+strings.Join(remExtFslice(nonNullPWnames), "\t"))

	var cellScores []float64 // cell score for all pathways > 0

	for _, c := range allCellNames {
		for _, pwName := range nonNullPWnames {
			cellScores = append(cellScores, allPW[idPW{c, pwName}]) // merge all scores for one cell
		}
		line := c + "\t" + strings.Join(floatStoStringS(cellScores), "\t")
		fmt.Fprintln(w, line)
		cellScores = nil
	}

	return w.Flush()
}

// WriteResultsTSV write results in TSV format
func WriteTSV(filename string, dataChan chan []string) error {
	// Open the file for writing, create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Create a new CSV writer with tab as the delimiter
	writer := csv.NewWriter(bufio.NewWriter(file))
	writer.Comma = '\t' // Set the column delimiter for TSV
	defer writer.Flush()

	// Read data from the channel and write to the file
	for data := range dataChan {
		if err := writer.Write(data); err != nil {
			log.Printf("Error writing line: %v", err)
			return err
		}
	}

	return nil
}
*/

// CreateResultFile create a result file
func CreateResultFile(db, dataFile string) {
	resFile := "results/" + db + "_" + dataFile
	Config.ResultFileName = resFile
	WriteHeader(resFile, Config.PathwayNames)
}

// close result file
func CloseResultFile() {
	file, err := os.OpenFile(Config.ResultFileName, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

// WriteHeader writes a slice of strings to a TSV file.
func WriteHeader(filename string, data []string) error {
	//fmt.Println("WriteHeader", filename)
	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, remove it
		err := os.Remove(filename)
		if err != nil {
			fmt.Println("Error removing file:", err)
			return err
		}
		//fmt.Println("Old file removed:", filename)
	} else if !os.IsNotExist(err) {
		// If there was an error other than "file does not exist"
		fmt.Println("Error checking file:", err)
		return err
	}
	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Create a new CSV writer with tab as the delimiter
	writer := csv.NewWriter(file)
	writer.Comma = '\t' // Set the delimiter to tab

	// append "ID" column
	data = remExtFslice(data)
	data = append([]string{"id"}, data...)
	// Write the data to the TSV file
	if err := writer.Write(data); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	// Write any buffered data to the underlying writer
	writer.Flush()

	// Check for any error during the flush
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
