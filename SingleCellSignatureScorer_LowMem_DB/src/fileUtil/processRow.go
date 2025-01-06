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
	"ScorerLowMem/src/pogrebdb"
	"ScorerLowMem/src/userTypes"
	"math"
	"sync"

	"github.com/akrylysov/pogreb"
)

// ###########################################################
// process one data line = one cell
func ProcessLine(fileDB *pogreb.DB, cellName string, colNames, sortedPWnames, geneValues []string, database map[string][]string, conf userTypes.CONF, wg *sync.WaitGroup) {
	defer wg.Done()
	geneNames, genesExpress := CleanZero(colNames, geneValues) // remove genes with null expression
	var geneExpressNoLog []float64

	// remove log 2 transformation if needed
	if conf.RemLog2 == 1 {
		geneExpressNoLog = UnLog2(genesExpress) // remove log2 for all genes
	} else {
		geneExpressNoLog = genesExpress
	}

	sumUMI := SliceSum(geneExpressNoLog)                 // UMI sum of all genes
	GNnormUMI := GeneExpDic(geneExpressNoLog, geneNames) // gene name -> gene express no log2 no zero

	cellScore := make([]float64, len(database))
	OneCellScore(geneNames, sortedPWnames, sumUMI, cellScore, database, GNnormUMI)

	// register the cell name and the scores in the file database
	pogrebdb.InsertColDB(fileDB, []byte(cellName), FloatSliceToByteSlice(cellScore))
}

// GeneExpDic create dict : gene name -> gene express no log2 no zero
func GeneExpDic(geneExpressNoLog []float64, geneNames []string) map[string]float64 {

	GNnormUMI := make(map[string]float64, len(geneNames)) // gene name -> gene express no log2 no zero
	for i, v := range geneNames {
		GNnormUMI[v] = geneExpressNoLog[i]
	}
	return GNnormUMI
}

// score computes the pathway score from UMI and the sum of all UMI
func Score(UMISum, sumUMI float64) float64 {
	return math.Round(1000*UMISum*100/sumUMI) / 1000
}

// OneCellScore computes the pathway score for one cell
func OneCellScore(geneNames, sortedPWnames []string, sumUMI float64, cellScore []float64, database map[string][]string, GNnormUMI map[string]float64) {
	nbPW := 0
	// for each sorted pathway in database (to have always the same pathway order)
	//for pwName, pwGenes := range database {
	for _, pwName := range sortedPWnames {
		pwGenes := database[pwName]
		UMISum := 0.0
		if len(geneNames) > 0 {
			genesInPW := IntersectionNeg(geneNames, pwGenes)
			if len(genesInPW) > 0 {
				for _, g := range genesInPW {
					if g[:1] == "-" { //if a sign "-" is detected then UMI of the gene without "-" g[1:]) is substracted
						UMISum = UMISum - GNnormUMI[g[1:]]
					} else {
						UMISum = UMISum + GNnormUMI[g]
					}
				}
			}
			//PWvalues = append(PWvalues, userTypes.PWscore{PwN: pwName, PwScore: UMISum * 100 / sumUMI})
			// UMISum * 100 / sumUMI to round it to 3 decimal : math.Round(f*1000)/1000
			cellScore[nbPW] = Score(UMISum, sumUMI)
		} else {
			//PWvalues = append(PWvalues, userTypes.PWscore{PwN: pwName, PwScore: 0.0})
			cellScore[nbPW] = 0.0
		}
		//fmt.Print(pwName + " ")
		nbPW++
	}
}
