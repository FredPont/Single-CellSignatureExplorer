// Testing file for scSignatureScorer
//  this file covers all the critical parts of the score calculation
// to use it, install Go programming language. the command line is : go test -v -cover

package main

import (
	"ScorerLowMem/src/fileUtil"
	"ScorerLowMem/src/userTypes"
	"fmt"
	"reflect"
	"testing"
)

// func TestMapKeysPWid(t *testing.T) {

// 	db := make(map[idPW]float64) // [cell name ; pw name] -> UMI
// 	db[idPW{"cell1", "Reactome_G2"}] = 100.
// 	db[idPW{"cell2", "Glycolyse"}] = 27.884615384615387
// 	db[idPW{"cell3", "MyPW"}] = 100.
// 	db[idPW{"cell4", "Verif_UMI"}] = 0.451612903225806
// 	db[idPW{"cell6", "OxPhos"}] = 27.884615384615387

// 	tests := []struct {
// 		mp   map[idPW]float64
// 		want []string
// 	}{
// 		{
// 			db,
// 			[]string{"Reactome_G2", "Glycolyse", "MyPW", "Verif_UMI", "OxPhos"},
// 		},
// 	}

// 	for _, s := range tests {
// 		got := mapKeysPWid(s.mp)
// 		sort.Strings(got)
// 		output := s.want
// 		sort.Strings(output)
// 		if !reflect.DeepEqual(got, output) {
// 			t.Errorf("mapKeysPWid was incorrect, got: %v, want : %v.", got, output)
// 		}

// 	}

// }

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
		got := fileUtil.UniqueStrings(s.s1)
		if !reflect.DeepEqual(got, s.want) {
			t.Errorf("uniqueStrings was incorrect, got: %v, want : %v.", got, s.want)
		}

	}

}

func TestFloatString(t *testing.T) {
	tests := []struct {
		f    float64
		want string
	}{
		{
			3.1459,
			"3.146",
		},
		{
			0.6184,
			"0.618",
		},
	}

	for _, s := range tests {
		got := fileUtil.FloatString(s.f)
		if !reflect.DeepEqual(got, s.want) {
			t.Errorf("FloatString was incorrect, got: %v, want : %v.", got, s.want)
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
	}

	for _, s := range tests {
		got := fileUtil.IntersectionNeg(s.s1, s.db)
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
		geneNames, genesExpress := fileUtil.CleanZero(s.colNames, s.geneValues)
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
	}

	for _, s := range tests {
		result := fileUtil.SliceSum(s.slice)
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
	}

	for _, s := range tests {
		result := fileUtil.UnLog2(s.slice)
		for i, r := range result {
			if r != s.res[i] {
				t.Errorf("unLog2 was incorrect, got: %f, want: %f.", r, s.res[i])
			}
		}

	}

}

func TestScore(t *testing.T) {
	tests := []struct {
		UMISum, sumUMI float64
		want           float64
	}{
		{5, 10, 50},
		{0, 10, 0},
		{10.5, 10.5, 100},
		{1.05, 10.5, 10},
		{1.05, 120.5, 0.871},
	}

	for _, s := range tests {

		got := fileUtil.Score(s.UMISum, s.sumUMI)
		if got != s.want {
			t.Errorf("Score was incorrect, got: %v, want : %v.", got, s.want)
		}

	}
}

func TestOneCellScore(t *testing.T) {
	db := make(map[string][]string)
	db["Reactome_G2"] = []string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"}
	dbneg := make(map[string][]string)
	dbneg["Reactome_G2"] = []string{"-E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"}
	tests := []struct {
		geneNames, sortedPWnames []string
		sumUMI                   float64
		//cellScore                []float64
		database  map[string][]string
		GNnormUMI map[string]float64
		want      []float64
	}{
		{[]string{"E2F1", "E2F3"}, []string{"Reactome_G2"}, 100, db, map[string]float64{"E2F1": 1., "E2F3": 4.}, []float64{5}},
		{[]string{"E2F1", "E2F3"}, []string{"Reactome_G2"}, 100, db, map[string]float64{"E2F1": 6., "E2F3": 4.}, []float64{10}},
		{[]string{"E2F1", "E2F3"}, []string{"Reactome_G2"}, 100, dbneg, map[string]float64{"E2F1": 4., "E2F3": 4.}, []float64{0}},
		{[]string{"E2F1", "E2F3"}, []string{"Reactome_G2"}, 100, dbneg, map[string]float64{"E2F1": 1., "E2F3": 4.}, []float64{3}},
		//{[]string{"E2F1", "E2F3"}, []string{"Reactome_G2"}, 100, dbneg, map[string]float64{"E2F1": 1., "E2F3": 4.}, []float64{4}},
		{[]string{"gx", "gy"}, []string{"Reactome_G2"}, 100, db, map[string]float64{"E2F1": 1., "E2F3": 4.}, []float64{0}},
	}

	for _, s := range tests {

		cellScore := make([]float64, len(s.database))

		fileUtil.OneCellScore(s.geneNames, s.sortedPWnames, s.sumUMI, cellScore, s.database, s.GNnormUMI)
		fmt.Println(cellScore, s.want)
		if !reflect.DeepEqual(cellScore, s.want) {
			t.Errorf("Score was incorrect, got: %v, want : %v.", cellScore, s.want)
		}

	}
}

// function return the idPW keys of a map[idPW]float64
// func resKeys(mymap map[idPW]float64) []idPW {
// 	keys := make([]idPW, len(mymap))
// 	i := 0
// 	for k := range mymap {
// 		keys[i] = k
// 		i++
// 	}
// 	return keys
// }

/*func TestProcessLine(t *testing.T) {
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
		conf       userTypes.CONF
		database   map[string][]string
		want       float64
	}{
		{
			"cell1",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2"},
			[]string{"0", "0.5", "2", "3", "0.3"},
			userTypes.CONF{0, 0, "0", 0, 0, "", ""},
			db,
			100.,
		},
		{
			"cell2",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"0", "0.5", "2", "3", "0.3", "10", "5"},
			userTypes.CONF{0, 0, "0", 0, 0, "", ""},
			db,
			27.884615384615387,
		},
		{
			"cell3",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			userTypes.CONF{1, 0, "1", 0, 0, "", ""},
			db2,
			100,
		},
		{
			"cell4",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			userTypes.CONF{1, 0, "1", 0, 0, "", ""},
			db1,
			45.1612903225806,
		},
		{
			"cell5",
			[]string{"gn1", "gn2", "gn3", "gn4", "gn5", "gn6", "gn7", "gn8", "gn9", "gn10"},
			[]string{"4", "3", "2", "0", "0", "5", "0", "0", "0", "1"},
			userTypes.CONF{1, 0, "1", 0, 0, "", ""},
			db3,
			90.3225806451613,
		},
		{
			"cell6",
			[]string{"E2F1", "E2F3", "CDK2", "CCNA1", "CCNA2", "g1", "g2"},
			[]string{"0", "0.5", "2", "3", "0.3", "10", "5"},
			userTypes.CONF{0, 0, "0", 0, 0, "", ""},
			db4, // test sign "-" in custom database only
			8.653846153846153,
		},
	}
	ch1 := make(chan int, 8)
	for _, id := range tests {
		//fmt.Println(id.cellName, id.colNames, id.geneValues, db, ch1, id.conf)
		go fileUtil.ProcessLine("tmp/", id.cellName, id.colNames, id.geneValues, id.database, ch1, id.conf)
	}

	for i := range tests {
		res := <-ch1
		fmt.Println("goroutine", i, "done ! chan =", res)
	}

	for _, id := range tests {

		csvFile, err := os.Open("tmp/" + id.cellName)
		check(err)
		defer csvFile.Close()
		reader := csv.NewReader(bufio.NewReader(csvFile))
		reader.Comma = '\t'
		row := ""
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("rec = ", record)
			row = record[1]
		}

		score, err := strconv.ParseFloat(row, 64)
		if err != nil {
			check(err)
		}

		got := math.Round(score*1000) / 1000
		want := math.Round(id.want*1000) / 1000

		fmt.Println("got: ", got, " want: ", want, " ", id.cellName, " db: ", id.database)

		if got != want {
			t.Errorf("Score was incorrect for %s, got: %f, want %f.", id.cellName, got, want)
		}

	}

}*/

/*func TestCalcTable(t *testing.T) {
	tests := []struct {
		nbPW     int
		colMax   int
		nbTables int
	}{
		{
			100,
			500,
			1,
		},
		{
			100,
			50,
			3,
		},
		{
			100,
			51,
			2,
		},
	}

	for _, s := range tests {
		got := fileUtil.CalcTable(s.nbPW, s.colMax)
		if !reflect.DeepEqual(got, s.nbTables) {
			t.Errorf("calcTable was incorrect, for %v pathawys and %v colmax, got: %v, want : %v.", s.nbPW, s.colMax, got, s.nbTables)
		}

	}

}*/

func TestSortPWscore(t *testing.T) {
	tests := []struct {
		unsorted []userTypes.PWscore
		sorted   []userTypes.PWscore
	}{
		{
			[]userTypes.PWscore{userTypes.PWscore{"x", 9.4}, userTypes.PWscore{"a", 0.26}},
			[]userTypes.PWscore{userTypes.PWscore{"a", 0.26}, userTypes.PWscore{"x", 9.4}},
		},
		{
			[]userTypes.PWscore{userTypes.PWscore{"GO_MITOCHONDRION.txt", 9.448708156257128}, userTypes.PWscore{"HALLMARK_SPERMATOGENESIS.txt", 0.2697506394971464}, userTypes.PWscore{"KEGG_GLYCOLYSIS_GLUCONEOGENESIS.txt", 0.525544029067155}, userTypes.PWscore{"KEGG_OXIDATIVE_PHOSPHORYLATION.txt", 1.496830635977891}, userTypes.PWscore{"GO_EMBRYONIC_ORGAN_DEVELOPMENT.txt", 1.6454186276786122}, userTypes.PWscore{"GO_GLUCOSE_CATABOLIC_PROCESS.txt", 0.3177883749335531}, userTypes.PWscore{"GO_GLUCOSE_METABOLIC_PROCESS.txt", 0.840903996402429}, userTypes.PWscore{"GO_MALE_MEIOSIS.txt", 0.03284925819750949}},
			[]userTypes.PWscore{userTypes.PWscore{"GO_EMBRYONIC_ORGAN_DEVELOPMENT.txt", 1.6454186276786122}, userTypes.PWscore{"GO_GLUCOSE_CATABOLIC_PROCESS.txt", 0.3177883749335531}, userTypes.PWscore{"GO_GLUCOSE_METABOLIC_PROCESS.txt", 0.840903996402429}, userTypes.PWscore{"GO_MALE_MEIOSIS.txt", 0.03284925819750949}, userTypes.PWscore{"GO_MITOCHONDRION.txt", 9.448708156257128}, userTypes.PWscore{"HALLMARK_SPERMATOGENESIS.txt", 0.2697506394971464}, userTypes.PWscore{"KEGG_GLYCOLYSIS_GLUCONEOGENESIS.txt", 0.525544029067155}, userTypes.PWscore{"KEGG_OXIDATIVE_PHOSPHORYLATION.txt", 1.496830635977891}},
		},
	}

	for _, s := range tests {
		fileUtil.SortPWscore(s.unsorted)
		if !reflect.DeepEqual(s.unsorted, s.sorted) {
			t.Errorf("SortPWscore was incorrect, got: %v, want : %v.", s.unsorted, s.sorted)
		}

	}

}

/*func TestRowBlocks(t *testing.T) {
	tests := []struct {
		allPWnames []string
		tableNB    int
		colMax     int
		bloks      [][]string
	}{
		{
			[]string{"a", "b", "c", "d", "e", "f"},
			2,
			4,
			[][]string{[]string{"a", "b", "c"}, []string{"d", "e", "f"}},
		},
		{
			[]string{"a", "b", "c", "d", "e", "f", "g", "h"},
			3,
			4,
			[][]string{[]string{"a", "b", "c"}, []string{"d", "e", "f"}, []string{"g", "h"}},
		},
		{
			[]string{"a", "b", "c", "d", "e", "f", "g", "h"},
			4,
			3,
			[][]string{[]string{"a", "b"}, []string{"c", "d"}, []string{"e", "f"}, []string{"g", "h"}},
		},
		{
			[]string{"a", "b", "c", "d", "e", "f"},
			1,
			10,
			[][]string{[]string{"a", "b", "c", "d", "e", "f"}},
		},
		{
			[]string{"a", "b", "c", "d", "e", "f"},
			1,
			0,
			[][]string{[]string{"a", "b", "c", "d", "e", "f"}},
		},
	}

	for _, s := range tests {
		got := fileUtil.RowBlocks(s.allPWnames, s.tableNB, s.colMax)
		if !reflect.DeepEqual(got, s.bloks) {
			t.Errorf("rowBlocks was incorrect, for %v pathawys, %v tables and %v colmax, got: %v, want : %v.", s.allPWnames, s.tableNB, s.colMax, got, s.bloks)
		}

	}

}*/
