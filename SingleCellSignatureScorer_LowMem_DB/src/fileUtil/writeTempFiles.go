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
	"ScorerLowMem/src/userTypes"
	"bufio"
	"fmt"
	"os"
)

func WriteTMP(tmpdir, cellName string, PWvalues []userTypes.PWscore) error {
	path := tmpdir + cellName
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	SortPWscore(PWvalues) // sort the slice of PWscore so the pathways are in the same order in all files

	for _, p := range PWvalues {
		line := p.PwN + "\t" + FloatString(p.PwScore)
		fmt.Fprintln(w, line)

	}

	return w.Flush()
}
