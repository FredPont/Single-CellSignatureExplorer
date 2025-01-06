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
	"fmt"
)

// ###########################################
func Header() {
	fmt.Println("")
	fmt.Println("   ┌─────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │ Single-Cell Scorer (c)Frederic PONT │")
	fmt.Println("   │ 2018-2020 - Free Software GNU GPL   │")
	fmt.Println("   └─────────────────────────────────────┘")
	//fmt.Println("")
}

// ###########################################
func DisplayDB(DBnames []string) {

	fmt.Println("")
	fmt.Println("   ┌─────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │           DataBases list            │")
	fmt.Println("   ├─────────────────────────────────────┤")
	fmt.Println("   │                                     │")
	for i, d := range DBnames {
		if i < 10 {
			fmt.Printf("   │ %-d  - %-31s│\n", i, d) // insert space for better alignment :)
		} else {
			fmt.Printf("   │ %-d - %-31s│\n", i, d)
		}
	}
	fmt.Println("   └─────────────────────────────────────┘")
	fmt.Println("   select working databases, default = 0")
	fmt.Println("   example : 1-3,5")
}
