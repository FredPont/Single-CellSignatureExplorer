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
//(c) Frederic Pont 2020

package fileUtil

import (
	"ScorerLowMem/src/userTypes"
	"os"
)

// add current datafile and database to conf
func UpdateConf(file, db string, cfig *userTypes.CONF) {
	f, _ := remExt(file)
	cfig.DBname = db
	cfig.DataName = f
}

// create temp dir and return dir name
func TmpDirName(file, db string) string {
	f, _ := remExt(file)
	path := "tmp/" + f + "_" + db + "/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	return path
}
