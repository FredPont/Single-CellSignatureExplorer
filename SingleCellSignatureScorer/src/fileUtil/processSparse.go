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
	"encoding/csv"
	"fmt"
	"sync"
)

func ProcessSparse(sm *SparseMatrix, writer *csv.Writer) {
	//fmt.Println(Config)
	// remove log 2 transformation if needed
	if Config.RemLog2 == 1 {
		sm.Data = unlog2Sparse(sm.Data) // remove log2 for all genes
	}

	// process sparse matrix

	var wg sync.WaitGroup
	maxGoroutines := Config.NbCPU
	// create a channel with the max number of job allowed
	ch := make(chan struct{}, maxGoroutines)
	// create a channel to signal the end of the printing
	done := make(chan bool)
	// create a channel to send the lines to write
	dataChan := make(chan []string, Config.ChunkSize*Config.ChunkInParallel) // Buffer of chunks lines size

	for row, rowData := range sm.Data {
		wg.Add(1) // Increment the WaitGroup counter
		ch <- struct{}{}
		cellName := sm.RowNames[row]
		GeneExpressMap := getGeneExpress(rowData)

		// compute the score for one cell and each pathway in the database
		switch Config.ScoreAlgo {
		case "sum":
			go score_sum(cellName, GeneExpressMap, rowData, &wg, ch, dataChan)
		case "mean":
			go score_mean(cellName, GeneExpressMap, rowData, &wg, ch, dataChan)
		default:
			// Handle the default case, which is "sum"
			go score_sum(cellName, GeneExpressMap, rowData, &wg, ch, dataChan)
		}
	}

	// close the channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	// Write CSV body
	go writeCSVFromChannel(dataChan, done, len(sm.Data), writer)
	// Notify main goroutine process is finished
	<-done
}

// built a map of gene name -> gene express for a row of the sparse matrix
func getGeneExpress(rowData map[int]float64) map[string]float64 {
	geneExpress := make(map[string]float64, len(rowData))
	geneColumnIndex := getKeys(rowData)
	//fmt.Println("rowData", rowData)
	for _, geneIndex := range geneColumnIndex {
		geneExpress[Config.ColNames[geneIndex+1]] = rowData[geneIndex] // +1 because the first column is the cell name
		//fmt.Println("gene", geneIndex, Config.ColNames[geneIndex+1], rowData[geneIndex])
	}
	return geneExpress

}

// score_sum calculates the score  for each pathway in database using the UMI sum of the genes in the pathway / the sum of all UMI of the cell
func score_sum(cellName string, GeneExpressMap map[string]float64, rowData map[int]float64, wg *sync.WaitGroup, ch chan struct{}, dataChan chan []string) {
	defer func() {
		<-ch // get the token back to free up a slot
		wg.Done()
	}()
	//fmt.Println(rowData)
	// final table row
	resultTableRow := make([]string, len(Config.PathwayNames)+1) // +1 because the first column is the cell name
	resultTableRow[0] = cellName
	sumUMI := mapSum(rowData)
	//fmt.Println("sumUMI", sumUMI)
	// for each pathway in database
	for i, pwGeneName := range Config.PathwayNames {
		pwGenes := Config.DataBase[pwGeneName]
		//fmt.Println("pwGenes", pwGenes)
		UMISum, _ := UMISumLoop(rowData, GeneExpressMap, pwGenes)
		//fmt.Println("UMISum", UMISum)
		resultTableRow[i+1] = scoreUMIsum(UMISum, sumUMI)

	}
	dataChan <- resultTableRow
	//fmt.Println(resultTableRow)
}

// UMISumLoop sums the UMI of a pathway
func UMISumLoop(rowData map[int]float64, GeneExpressMap map[string]float64, pwGenes []string) (float64, int) {
	UMISum := 0.0
	expressedGeneCount := 0 // count the number of expressed genes in the pathway
	// pwGenes is a list of genes in the pathway
	// if a gene is prefixed with "-" then the UMI of the gene is subtracted
	// if a gene is not prefixed with "-" then the UMI of the gene is added
	//fmt.Println("pwGenes", pwGenes)
	for _, pwGene := range pwGenes {
		if pwGene[:1] == "-" { //if a sign "-" is detected then UMI of the gene without "-" g[1:]) is substracted
			UMISum = UMISum - GeneExpressMap[pwGene[1:]]
		} else {
			UMISum = UMISum + GeneExpressMap[pwGene]
		}
		// increment the expressed gene count if the gene is expressed
		if _, exists := GeneExpressMap[pwGene]; exists {
			expressedGeneCount++
		}
	}
	return UMISum, expressedGeneCount
}

func scoreUMIsum(UMISum, sumUMI float64) string {
	if UMISum != 0. && sumUMI != 0. {
		score := UMISum * 100. / sumUMI // the score is the UMI %
		return fmt.Sprintf("%.2f", score)
	}
	return "0"
}

func score_mean(cellName string, GeneExpressMap map[string]float64, rowData map[int]float64, wg *sync.WaitGroup, ch chan struct{}, dataChan chan []string) {
	defer func() {
		<-ch // get the token back to free up a slot
		wg.Done()
	}()
	//fmt.Println(rowData)
	// final table row
	resultTableRow := make([]string, len(Config.PathwayNames)+1) // +1 because the first column is the cell name
	resultTableRow[0] = cellName
	sumUMI := mapSum(rowData)

	// for each pathway in database
	for i, pwGeneName := range Config.PathwayNames {
		pwGenes := Config.DataBase[pwGeneName] // pwGeneName is the pathway name, pwGenes is the list of genes in the pathway
		//fmt.Println("pwGenes", pwGenes)
		UMISum, _ := UMISumLoop(rowData, GeneExpressMap, pwGenes)
		//fmt.Println("UMISum", UMISum)
		resultTableRow[i+1] = scoreUMImean(UMISum, sumUMI, len(pwGenes)) // len(pwGenes) is the number of genes in the pathway
		//fmt.Println(pwGeneName, resultTableRow[0], resultTableRow[i+1])
	}
	dataChan <- resultTableRow
	//fmt.Println(resultTableRow)
}

func scoreUMImean(UMISum, sumUMI float64, expressedGeneCount int) string {
	//fmt.Println("UMISum", UMISum, "sumUMI", sumUMI, "pwGenes", expressedGeneCount, "nbTotalGenes", len(Config.ColNames)-1)
	if UMISum != 0. && sumUMI != 0. && expressedGeneCount > 0 {
		score := UMISum / float64(expressedGeneCount) * float64(len(Config.ColNames)-1) / sumUMI // the score is the MEAN(UMI of the pathway)/MEAN(all UMI)
		return fmt.Sprintf("%.2f", score)
	}
	return "0"
}
