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

// CONF stores software parameters
// always use maj for conf variables
type CONF struct {
	RemLog2         int                 `json:"removeLog2"`
	Server          int                 `json:"server"`
	DBserver        string              `json:"DBserver"`
	ScoreAlgo       string              `json:"scoreAlgo"`       // if true, score is the mean of the genes in the pathway / mean of all genes in the cell
	ChunkSize       int                 `json:"ChunkSize"`       // number of lines in a chunk
	ChunkInParallel int                 `json:"ChunkInParallel"` // number of chunks in parallel
	NbCPU           int                 `json:"NbCPU"`           // number of CPUs
	ColNames        []string            `json:-`                 // column names of the count matrix
	GeneIndex       map[string]int      `json:-`                 // gene name => index in the count matrix
	NbCols          int                 // number of columns in the count matrix
	DataBase        map[string][]string `json:-` // pathway name -> genes names
	PathwayNames    []string            `json:-` // pathway names
	ResultFileName  string              // name of the result file
}

// Declare a global variable of type Conf
var Config CONF
