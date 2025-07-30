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

// SparseMatrix stores a sparse matrix
type SparseMatrix struct {
	Rows     int
	Cols     int
	ColNames []string
	RowNames []string
	Data     map[int]map[int]float64
}

// NewSparseMatrix creates a new sparse matrix
func NewSparseMatrix(rows, cols int, colNames []string) *SparseMatrix {
	return &SparseMatrix{
		Rows:     rows,
		Cols:     cols,
		ColNames: colNames,
		RowNames: make([]string, rows),
		Data:     make(map[int]map[int]float64),
	}
}

// SetValue sets a value in the sparse matrix
func (sm *SparseMatrix) SetValue(row, col int, value float64) {
	if sm.Data[row] == nil {
		sm.Data[row] = make(map[int]float64)
	}
	sm.Data[row][col] = value
}

// GetValue returns a value from the sparse matrix
func (sm *SparseMatrix) GetValue(row, col int) float64 {
	if sm.Data[row] == nil {
		return 0
	}
	return sm.Data[row][col]
}
