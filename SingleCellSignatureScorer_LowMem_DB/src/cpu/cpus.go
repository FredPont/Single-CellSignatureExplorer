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

package cpu

import (
	"runtime"
)

// this function compute the optimal nb of CPUs to use
// by default all CPU are used if the user does not limit it
func CalcCPU(confCPU int) int {

	totalCPU := runtime.NumCPU()

	if confCPU != 0 && confCPU < totalCPU {
		return confCPU
	}
	return totalCPU
}
