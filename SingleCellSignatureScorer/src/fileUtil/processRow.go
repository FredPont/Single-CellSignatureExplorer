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

// ###########################################################
// process one data line = one cell
func processLine(cellName string, colNames, geneValues []string, database map[string][]string, ch1 chan<- map[idPW]float64, conf CONF) {
	geneNames, genesExpress := cleanZero(colNames, geneValues) // remove genes with null expression
	var geneExpressNoLog []float64

	// remove log 2 transformation if needed
	if conf.RemLog2 == 1 {
		geneExpressNoLog = unLog2(genesExpress) // remove log2 for all genes
	} else {
		geneExpressNoLog = genesExpress
	}

	sumUMI := sliceSum(geneExpressNoLog)                 // UMI sum of all genes
	GNnormUMI := GeneExpDic(geneExpressNoLog, geneNames) // gene name -> gene express no log2 no zero

	tmpPWvalues := make(map[idPW]float64) // [cell name ; pw name] -> UMI
	// for each pathway in database
	for pwName, pwGenes := range database {

		UMISum := 0.0
		if len(geneNames) > 0 {
			genesInPW := intersectionNeg(geneNames, pwGenes)
			if len(genesInPW) > 0 {
				for _, g := range genesInPW {
					if g[:1] == "-" { //if a sign "-" is detected then UMI of the gene without "-" g[1:]) is substracted
						UMISum = UMISum - GNnormUMI[g[1:]]
					} else {
						UMISum = UMISum + GNnormUMI[g]
					}
				}
			}
			tmpPWvalues[idPW{cellName, pwName}] = UMISum * 100. / sumUMI // the score is the UMI %

		} else {
			tmpPWvalues[idPW{cellName, pwName}] = 0.0
		}
	}
	ch1 <- tmpPWvalues

}

// GeneExpDic create dict : gene name -> gene express no log2 no zero
func GeneExpDic(geneExpressNoLog []float64, geneNames []string) map[string]float64 {

	GNnormUMI := make(map[string]float64, len(geneNames)) // gene name -> gene express no log2 no zero
	for i, v := range geneNames {
		GNnormUMI[v] = geneExpressNoLog[i]
	}
	return GNnormUMI
}

// remove pathways with zero score in all cells
func removeNullPW(allCellNames []string, allPW map[idPW]float64) []string {
	allPWnames := mapKeysPWid(allPW)
	var nonNullPWnames []string

	for _, pwName := range allPWnames {
		var pathwayScores []float64
		for _, cellName := range allCellNames {
			pathwayScores = append(pathwayScores, allPW[idPW{cellName, pwName}])
		}
		if sliceSum(pathwayScores) != 0. {
			nonNullPWnames = append(nonNullPWnames, pwName)
		}
	}
	return nonNullPWnames
}
