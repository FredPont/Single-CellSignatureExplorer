package fileUtil

import (
	"fmt"
	"sync"
	"testing"
)

// Test for score_sum function
func TestScoreSum(t *testing.T) {
	wg := &sync.WaitGroup{}
	ch := make(chan struct{}, 1) // Create a buffered channel
	dataChan := make(chan []string)

	// Mock Config structure for testing
	var Config = struct {
		PathwayNames []string
		DataBase     map[string][]string
	}{
		PathwayNames: []string{"Pathway1", "Pathway2"},
		DataBase: map[string][]string{
			"Pathway1": {"GeneA", "GeneB"},
			"Pathway2": {"GeneC", "-GeneD"},
		},
	}

	fmt.Println(Config)

	// Example input data
	cellName := "Cell1"
	GeneExpressMap := map[string]float64{
		"GeneA": 10.0,
		"GeneB": 20.0,
		"GeneC": 5.0,
		"GeneD": 15.0,
	}
	rowData := map[int]float64{
		0: 10.0,
		1: 20.0,
	}

	// Start the test
	wg.Add(1)
	go score_sum(cellName, GeneExpressMap, rowData, wg, ch, dataChan)

	// Wait for the result
	wg.Wait()

	// Check the result
	result := <-dataChan
	close(dataChan)
	expected := []string{"Cell1", "100.00", "0.00"} // Adjust based on expected scores

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expected[i], result[i])
		}
	}
}
