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
//(c) Frederic Pont 2023

package userTypes

// PWscore  cell name -> pathway score
type PWscore struct {
	PwN     string  // pathway name
	PwScore float64 // pathway score
}

// CONF stores software parameters
// always use maj for conf variables
type CONF struct {
	RemLog2  int    `json:"removeLog2"`
	Server   int    `json:"server"`
	DBserver string `json:"DBserver"`
	CPU      int    `json:"CPU"`
	//ResSze   int    `json:"ResultSize"`
	DataName string // name of the csv datafile
	DBname   string // name of the pathway database
}
