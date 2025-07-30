// Testing file for scSignatureScorer
//  this file covers all the critical parts of the score calculation
// to use it, install Go programming language. the command line is : go test -v -cover

package fileUtil

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

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

// test the read pathway function, in particular when some genes have special char or spaces
func TestReadPW(t *testing.T) {
	tests := []struct {
		path string
		res  []string
	}{
		{
			"../test_files/Test_metabolism.txt",
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

func TestUMISumLoop(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		rowData        map[int]float64
		geneExpressMap map[string]float64
		pwGenes        []string
		expected       float64
	}{
		{
			name:    "Single gene addition",
			rowData: map[int]float64{},
			geneExpressMap: map[string]float64{
				"GeneA": 10.0,
			},
			pwGenes:  []string{"GeneA"},
			expected: 10.0,
		},
		{
			name:    "Single gene subtraction",
			rowData: map[int]float64{},
			geneExpressMap: map[string]float64{
				"GeneA": 10.0,
			},
			pwGenes:  []string{"-GeneA"},
			expected: -10.0,
		},
		{
			name:    "Multiple genes with addition and subtraction",
			rowData: map[int]float64{},
			geneExpressMap: map[string]float64{
				"GeneA": 10.0,
				"GeneB": 20.0,
				"GeneC": 5.0,
			},
			pwGenes:  []string{"GeneA", "GeneB", "-GeneC"},
			expected: 25.0, // 10 + 20 - 5
		},
		{
			name:           "No genes",
			rowData:        map[int]float64{},
			geneExpressMap: map[string]float64{},
			pwGenes:        []string{},
			expected:       0.0,
		},
		{
			name:    "Gene not found in map",
			rowData: map[int]float64{},
			geneExpressMap: map[string]float64{
				"GeneA": 10.0,
			},
			pwGenes:  []string{"GeneA", "-GeneB"}, // GeneB is not in the map
			expected: 10.0,                        // 10 + 0 (GeneB not found)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := UMISumLoop(tt.rowData, tt.geneExpressMap, tt.pwGenes)
			if got != tt.expected {
				t.Errorf("UMISumLoop() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestScoreUMIsum(t *testing.T) {
	// Test cases
	tests := []struct {
		name   string
		UMISum float64
		sumUMI float64
		want   string
	}{
		{
			name:   "Positive UMISum and sumUMI",
			UMISum: 50.0,
			sumUMI: 100.0,
			want:   "50.00",
		},
		{
			name:   "Zero UMISum",
			UMISum: 0.0,
			sumUMI: 100.0,
			want:   "0",
		},
		{
			name:   "Zero sumUMI",
			UMISum: 50.0,
			sumUMI: 0.0,
			want:   "0",
		},
		{
			name:   "Negative UMISum",
			UMISum: -5.3,
			sumUMI: 100.0,
			want:   "-5.30",
		},
		{
			name:   "Both UMISum and sumUMI are zero",
			UMISum: 0.0,
			sumUMI: 0.0,
			want:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scoreUMIsum(tt.UMISum, tt.sumUMI)
			if got != tt.want {
				t.Errorf("scoreUMIsum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScoreUMImean(t *testing.T) {
	// Test cases
	Config.ColNames = []string{"CellName", "GeneA", "GeneB", "GeneC"} // Mocking Config.ColNames for testing
	tests := []struct {
		name    string
		UMISum  float64
		sumUMI  float64
		pwGenes []string
		want    string
	}{
		{
			name:    "Positive UMISum and sumUMI",
			UMISum:  50.0,
			sumUMI:  100.0,
			pwGenes: []string{"GeneA", "GeneB"},
			want:    "0.75",
		},
		{
			name:    "Zero UMISum",
			UMISum:  0.0,
			sumUMI:  100.0,
			pwGenes: []string{"GeneA", "GeneB"},
			want:    "0",
		},
		{
			name:    "Zero sumUMI",
			UMISum:  50.0,
			sumUMI:  0.0,
			pwGenes: []string{"GeneA", "GeneB"},
			want:    "0",
		},
		{
			name:    "Negative UMISum",
			UMISum:  -5.3,
			sumUMI:  100.0,
			pwGenes: []string{"GeneA", "GeneB"},
			want:    "-0.08",
		},
		{
			name:    "Both UMISum and sumUMI are zero",
			UMISum:  0.0,
			sumUMI:  0.0,
			pwGenes: []string{"GeneA", "GeneB"},
			want:    "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scoreUMImean(tt.UMISum, tt.sumUMI, len(tt.pwGenes))
			if got != tt.want {
				t.Errorf("scoreUMIsum() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two string slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Test for getKeys function
func TestGetKeys(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		rowData  map[int]float64
		expected []int
	}{
		{
			name: "Basic case with multiple keys",
			rowData: map[int]float64{
				0: 10.5,
				4: 25.0,
				1: 20.3,
				2: 15.0,
			},
			expected: []int{0, 4, 1, 2},
		},
		{
			name: "Case with non-sequential keys",
			rowData: map[int]float64{
				1: 5.0,
				3: 12.0,
				7: 8.0,
			},
			expected: []int{1, 3, 7},
		},
		{
			name:     "Empty map",
			rowData:  map[int]float64{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getKeys(tt.rowData)
			sort.Ints(got)
			sort.Ints(tt.expected) // Sort expected for comparison
			// Check if the keys match

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getKeys() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test for getGeneExpress function
func TestGetGeneExpress(t *testing.T) {
	// Test cases
	Config.ColNames = []string{"CellName", "GeneA", "GeneB", "GeneC"} // Mocking Config.ColNames for testing
	tests := []struct {
		name     string
		rowData  map[int]float64
		expected map[string]float64
	}{
		{
			name: "Basic case with all genes",
			rowData: map[int]float64{
				0: 10.5,
				1: 20.3,
				2: 15.0,
			},
			expected: map[string]float64{
				"GeneA": 10.5,
				"GeneB": 20.3,
				"GeneC": 15.0,
			},
		},
		{
			name: "Case with missing genes",
			rowData: map[int]float64{
				0: 5.0,
				2: 12.0,
			},
			expected: map[string]float64{
				"GeneA": 5.0,
				"GeneC": 12.0,
			},
		},
		{
			name:     "Empty rowData",
			rowData:  map[int]float64{},
			expected: map[string]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getGeneExpress(tt.rowData)
			if !equalMaps(got, tt.expected) {
				t.Errorf("getGeneExpress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to compare two maps
func equalMaps(a, b map[string]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valueA := range a {
		if valueB, exists := b[key]; !exists || valueA != valueB {
			return false
		}
	}
	return true
}
