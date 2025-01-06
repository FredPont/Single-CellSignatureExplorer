// Testing file for scSignatureScorer
//  this file covers all the critical parts of the score calculation
// to use it, install Go programming language. the command line is : go test -v -cover

package main

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestRemoveNullPW(t *testing.T) {

	db := make(map[idPW]float64) // [cell name ; pw name] -> UMI
	db[idPW{"cell1", "Glycolyse"}] = 0.
	db[idPW{"cell2", "Glycolyse"}] = 0.
	db[idPW{"cell3", "Glycolyse"}] = 0.
	db[idPW{"cell4", "Glycolyse"}] = 0.
	db[idPW{"cell6", "Glycolyse"}] = 0.
	db[idPW{"cell1", "Reactome_G2"}] = 100.
	db[idPW{"cell2", "Reactome_G2"}] = 0.
	db[idPW{"cell3", "Reactome_G2"}] = 100.
	db[idPW{"cell4", "Reactome_G2"}] = 0.451612903225806
	db[idPW{"cell6", "Reactome_G2"}] = 0.

	tests := []struct {
		cellNames []string
		mp        map[idPW]float64
		want      []string
	}{
		{[]string{"cell1", "cell2", "cell3", "cell4", "cell5", "cell6"},
			db,
			[]string{"Reactome_G2"},
		},
	}

	for _, s := range tests {
		got := removeNullPW(s.cellNames, s.mp)
		sort.Strings(got)
		output := s.want
		sort.Strings(output)
		if !reflect.DeepEqual(got, output) {
			t.Errorf("removeNullPW was incorrect, got: %v, want : %v.", got, output)
		}

	}

}

func TestMapKeysPWid(t *testing.T) {

	db := make(map[idPW]float64) // [cell name ; pw name] -> UMI
	db[idPW{"cell1", "Reactome_G2"}] = 100.
	db[idPW{"cell2", "Glycolyse"}] = 27.884615384615387
	db[idPW{"cell3", "MyPW"}] = 100.
	db[idPW{"cell4", "Verif_UMI"}] = 0.451612903225806
	db[idPW{"cell6", "OxPhos"}] = 27.884615384615387

	tests := []struct {
		mp   map[idPW]float64
		want []string
	}{
		{
			db,
			[]string{"Reactome_G2", "Glycolyse", "MyPW", "Verif_UMI", "OxPhos"},
		},
	}

	for _, s := range tests {
		got := mapKeysPWid(s.mp)
		sort.Strings(got)
		output := s.want
		sort.Strings(output)
		if !reflect.DeepEqual(got, output) {
			t.Errorf("mapKeysPWid was incorrect, got: %v, want : %v.", got, output)
		}

	}

}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		s1   []string
		want []string
	}{
		{
			[]string{"g1", "g2", "E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"g1", "g2", "E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"},
		},
		{
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2", "E2F1"},
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
		},
	}

	for _, s := range tests {
		got := uniqueStrings(s.s1)
		if !reflect.DeepEqual(got, s.want) {
			t.Errorf("uniqueStrings was incorrect, got: %v, want : %v.", got, s.want)
		}

	}

}

func TestIntersectionNeg(t *testing.T) {
	tests := []struct {
		s1   []string
		db   []string
		want []string
	}{
		{
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"E2F3", "-CDK2", "CCNA1", "CCNA2"}, // test sign "-" in the database only, never in sample table gene names
			[]string{"E2F3", "-CDK2", "CCNA1", "CCNA2"},
		},
		{
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"E2F3", "CDK2", "CCNA1", "-CCNA2"}, // test sign "-" in the database only, never in sample table gene names
			[]string{"E2F3", "CDK2", "CCNA1", "-CCNA2"},
		},
		{
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"E2F3", "CDK2", "CCNA1", "CCNA2"},
			[]string{"E2F3", "CDK2", "CCNA1", "CCNA2"},
		},
		{
			[]string{"18S rRNA", "28S rRNA", "5.8S rRNA"},
			[]string{"18S rRNA", "28S rRNA", "5.8S rRNA", "5S rRNA", "AAAS"},
			[]string{"18S rRNA", "28S rRNA", "5.8S rRNA"},
		},
		{
			[]string{"18S rRNA", "28S rRNA", "5.8S rRNA", "5S rRNA", "AAAS"},
			[]string{"5S rRNA", "18S rRNA", "28S rRNA"},
			[]string{"5S rRNA", "18S rRNA", "28S rRNA"},
		},
	}

	for _, s := range tests {
		got := intersectionNeg(s.s1, s.db)
		if !reflect.DeepEqual(got, s.want) {
			t.Errorf("intersectionNeg was incorrect, got: %v, want : %v.", got, s.want)
		}

	}

}

func TestCleanZero(t *testing.T) {
	tests := []struct {
		colNames   []string
		geneValues []string
		wantCol    []string
		wantExp    []float64
	}{
		{
			[]string{"E2F1", "E2F3", "-CDK2", "CCNA1", "CCNA2", "g1", "g2"}, // test sign "-"
			[]string{"0", "0.5", "2", "3", "0.3", "10", "5"},
			[]string{"E2F3", "-CDK2", "CCNA1", "CCNA2", "g1", "g2"}, // test sign "-"
			[]float64{0.5, 2, 3, 0.3, 10, 5},
		},
		{
			[]string{"E2F1", "E2F3", "-CDK2", "CCNA1", "CCNA2", "g1", "g2"}, // test sign "-"
			[]string{"0", "0.5", "2", "3", "0.3", "0", "0"},
			[]string{"E2F3", "-CDK2", "CCNA1", "CCNA2"}, // test sign "-"
			[]float64{0.5, 2, 3, 0.3},
		},
	}

	for _, s := range tests {
		geneNames, genesExpress := cleanZero(s.colNames, s.geneValues)
		if !reflect.DeepEqual(geneNames, s.wantCol) {
			t.Errorf("Clean Zero was incorrect, got: %v, want : %v.", geneNames, s.wantCol)
		}
		if !reflect.DeepEqual(genesExpress, s.wantExp) {
			t.Errorf("Clean Zero was incorrect, got: %v, want : %v.", genesExpress, s.wantExp)
		}

	}

}

func TestSliceSum(t *testing.T) {
	tests := []struct {
		slice []float64
		res   float64
	}{
		{
			[]float64{1., 2., 3.5},
			6.5,
		},
		{
			[]float64{1, -2, 3.5},
			2.5,
		},
		{
			[]float64{4, 3, 2, 0, 0, 5, 0, 0, 0, 1},
			15,
		},
		{
			[]float64{16, 8, 4, 1, 1, 32, 1, 1, 1, 2},
			67,
		},
	}

	for _, s := range tests {
		result := sliceSum(s.slice)
		if result != s.res {
			t.Errorf("sliceSum was incorrect, got: %f, want: %f.", result, s.res)
		}

	}

}

func TestUnLog2(t *testing.T) {

	tests := []struct {
		slice []float64
		res   []float64
	}{
		{
			[]float64{1., 2., 0},
			[]float64{2., 4., 1},
		},
		{
			[]float64{4, 3, 2, 0, 0, 5, 0, 0, 0, 1},
			[]float64{16, 8, 4, 1, 1, 32, 1, 1, 1, 2},
		},
	}

	for _, s := range tests {
		result := unLog2(s.slice)
		for i, r := range result {
			if r != s.res[i] {
				t.Errorf("unLog2 was incorrect, got: %f, want: %f.", r, s.res[i])
			}
		}

	}

}

// function return the idPW keys of a map[idPW]float64
func resKeys(mymap map[idPW]float64) []idPW {
	keys := make([]idPW, len(mymap))
	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}
	return keys
}

func TestProcessLine(t *testing.T) {
	// return keys in map

	db := make(map[string][]string)
	db["Reactome_G2"] = []string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"}
	db2 := make(map[string][]string)
	db2["Verif_UMI"] = []string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"}
	db1 := make(map[string][]string)
	db1["Verif_UMI"] = []string{"gn1", "gn2", "gn3", "gn4"}
	db3 := make(map[string][]string)
	db3["Verif_UMI"] = []string{"gn1", "gn2", "gn6"}
	db4 := make(map[string][]string)
	db4["Reactome_G2"] = []string{"E2F1", "E2F3", "-CDK2", "CCNA1", "CCNA2"}

	tests := []struct {
		cellName   string
		colNames   []string
		geneValues []string
		conf       CONF
		database   map[string][]string
		want       float64
		id         idPW
	}{
		{
			"cell1",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"},
			[]string{"0", "0.5", "2", "3", "0.3"},
			CONF{0, 0, "0"},
			db,
			100.,
			idPW{"cell1", "Reactome_G2"},
		},
		{
			"cell2",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"0", "0.5", "2", "3", "0.3", "10", "5"},
			CONF{0, 0, "0"},
			db,
			27.884615384615387,
			idPW{"cell2", "Reactome_G2"},
		},
		{
			"cell3",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{1, 0, "1"},
			db2,
			100,
			idPW{"cell3", "Verif_UMI"},
		},
		{
			"cell4",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{1, 0, "1"},
			db1,
			45.1612903225806,
			idPW{"cell4", "Verif_UMI"},
		},
		{
			"cell5",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{1, 0, "1"},
			db3,
			90.3225806451613,
			idPW{"cell5", "Verif_UMI"},
		},
		{
			"cell6",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"0", "0.5", "2", "3", "0.3", "10", "5"},
			CONF{0, 0, "0"},
			db4, // test sign "-" in custom database only
			8.653846153846153,
			idPW{"cell6", "Reactome_G2"},
		},
		{
			"cell7",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6"},
			[]string{"0", "0.0", "0", "0", "0.", "0", "0"},
			CONF{0, 0, "0"},
			db3,
			0.,
			idPW{"cell7", "Verif_UMI"},
		},
		{
			"cell8",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6"},
			[]string{"0", "0.0", "1", "-1", "0.", "0", "0"},
			CONF{0, 0, "0"},
			db3,
			0.,
			idPW{"cell8", "Verif_UMI"},
		},
		{
			"cell9",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"},
			[]string{"0", "0.5", "2", "3", "0.3"},
			CONF{0, 0, "0"},
			db, // test sign "-" in custom database only
			100.,
			idPW{"cell9", "Reactome_G2"},
		},
		{
			"cell10",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{0, 0, "0"},
			db, // test sign "-" in custom database only
			0.,
			idPW{"cell10", "Reactome_G2"},
		},
		{
			"cell11",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{0, 0, "0"},
			db1,
			60.,
			idPW{"cell11", "Verif_UMI"},
		},
		{
			"cell12",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			CONF{1, 0, "1"},
			db1,
			45.16129032258064,
			idPW{"cell12", "Verif_UMI"},
		},
		{
			"cell13",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"1409.7", "16985.1", "26938", "29796.6", "37908.5", "35658", "2813.1", "34903.5", "25912.4", "38446.3"},
			CONF{0, 0, "1"},
			db1,
			29.959341423576547,
			idPW{"cell13", "Verif_UMI"},
		},
	}
	ch1 := make(chan map[idPW]float64, 8)
	for _, id := range tests {
		//fmt.Println(id.cellName, id.colNames, id.geneValues, db, ch1, id.conf)
		go processLine(id.cellName, id.colNames, id.geneValues, id.database, ch1, id.conf)
	}

	for x := 0; x < len(tests); x++ {
		res := <-ch1

		//id := tests[x]
		var id struct {
			cellName   string
			colNames   []string
			geneValues []string
			conf       CONF
			database   map[string][]string
			want       float64
			id         idPW
		}
		keys := resKeys(res)
		k := keys[0]

		for _, t := range tests {
			if k == t.id {
				id = t
			}
		}

		got := math.Round(res[k]*1000) / 1000
		want := math.Round(id.want*1000) / 1000

		fmt.Println("got: ", got, " want: ", want, " ", id.cellName, " db: ", id.database, keys)

		if got != want {
			t.Errorf("Score was incorrect for %s, got: %f, want %f.", id.cellName, got, want)
		}

	}

}

// test the read pathway function, in particular when some genes have special char or spaces
func TestReadPW(t *testing.T) {
	tests := []struct {
		path string
		res  []string
	}{
		{
			"src/Test_metabolism.txt",
			[]string{"18S rRNA", "28S rRNA", "5.8S rRNA", "5S rRNA", "AAAS"},
		},
	}

	for _, s := range tests {
		got := readPW(s.path)
		want := s.res
		fmt.Println(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("reading pathw was incorrect, got: %v, want : %v.", got, want)
		}

	}
}
