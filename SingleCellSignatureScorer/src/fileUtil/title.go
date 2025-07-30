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
	"fmt"
	"time"
)

// ###########################################
func Header() {
	version := softVersion()
	fmt.Println("")
	fmt.Println("   ┌──────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │ Single-Cell Signature Scorer v" + version + " │")
	fmt.Println("   │             (c)Frederic PONT             │")
	fmt.Println("   │   2018-2025 - Free Software GNU GPL      │")
	fmt.Println("   └──────────────────────────────────────────┘")
	//fmt.Println("")
}

func softVersion() string {
	// Get the current date and time
	now := time.Now()

	// Format the date as "YYYY-MM-DD"
	formattedDate := now.Format("2006-01-02")
	return formattedDate
}
