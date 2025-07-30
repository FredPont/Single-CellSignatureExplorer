package fileUtil

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// ReadMatrix reads a matrix from a CSV file
func ReadMatrix(dataFile string) {
	t0 := time.Now()
	//fmt.Println("Begin data processing...")

	// Open the file for writing, create it if it doesn't exist
	fileOut, err := os.OpenFile(Config.ResultFileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileOut.Close()

	// Create a new CSV writer with tab as the delimiter
	writer := csv.NewWriter(bufio.NewWriter(fileOut))
	writer.Comma = '\t' // Set the column delimiter for TSV
	defer writer.Flush()

	// Open the data file
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader with a buffer of 4KB
	reader := csv.NewReader(bufio.NewReaderSize(file, 4096)) // Buffer de 4KB
	reader.Comma = '\t'                                      // Set the column delimiter for TSV
	reader.FieldsPerRecord = -1                              // Allow a variable number of fields per record
	reader.Comment = '#'                                     // Ignore lines of comment

	// Read the header
	Config.ColNames, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	geneColumnIndex() // Fill the map geneIndex with gene names and their index

	// Read the file by chunks of n lines
	chunkSize := Config.ChunkSize
	//fmt.Printf("Reading file by chunks of %d lines...\n", chunkSize)

	var wg sync.WaitGroup
	chunkChan := make(chan [][]string, Config.ChunkInParallel) // Buffer of nb ChunkInParallel chunks

	// Launch goroutines to process chunks
	for i := 0; i < Config.ChunkInParallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				ProcessChunk(chunk, writer)
			}
		}()
	}

	// Read and send chunks to goroutines for processing
	go func() {
		for {
			chunk, err := readChunk(reader, chunkSize)
			if err != nil {
				if err == io.EOF {
					close(chunkChan) // Close the channel when EOF is reached
					break
				}
				log.Fatal(err)
			}
			chunkChan <- chunk
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	t1 := time.Now()
	fmt.Printf("Elapsed time: %v\n", t1.Sub(t0))
}

// readChunk reads a chunk of lines from the CSV reader
func readChunk(reader *csv.Reader, chunkSize int) ([][]string, error) {
	var chunk [][]string
	for i := 0; i < chunkSize; i++ {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				if len(chunk) > 0 {
					return chunk, nil
				}
				return nil, io.EOF
			}
			return nil, err
		}
		chunk = append(chunk, record)
	}
	return chunk, nil
}

// ProcessChunk processes a chunk of lines
func ProcessChunk(chunk [][]string, writer *csv.Writer) {
	sm := chunkToSparseMatrix(chunk)
	ProcessSparse(sm, writer)
}

// chunkToSparseMatrix creates a sparse matrix from a chunk of lines
func chunkToSparseMatrix(chunk [][]string) *SparseMatrix {
	// create a new sparse matrix
	sm := NewSparseMatrix(len(chunk), len(Config.ColNames)-1, Config.ColNames[1:])

	// fill the sparse matrix with the data from the CSV file
	for row, record := range chunk {

		//if row < len(sm.RowNames) {
		sm.RowNames[row] = record[0] // suppose that the first column contains the row names
		//}
		for col, value := range record[1:] { // ignore the first column
			// convert the value to float64 and add it to the sparse matrix
			// check if the value is not empty or null
			if value != "" && value != "0" && value != "0.0" {
				val, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Fatal(err)
				}
				sm.SetValue(row, col, val)
			}
		}
	}
	//fmt.Println("sparse matrix created", sm)
	return sm
}

// geneColumnIndex creates a map with gene name => column index in the sparse matrix
// Config.ColNames contains the column names of the TSV file (first column is the cell name)
// so geneIndex[0] is the name of the first column (ID for example)
func geneColumnIndex() error {
	if Config.ColNames == nil || len(Config.ColNames) == 0 {
		return fmt.Errorf("ColNames cannot be nil or empty")
	}

	geneIndex := make(map[string]int)
	for i, gene := range Config.ColNames {
		if _, exists := geneIndex[gene]; exists {
			return fmt.Errorf("duplicate detected for gene: %s", gene)
		}
		geneIndex[gene] = i
	}

	Config.GeneIndex = geneIndex // Update the global configuration
	//fmt.Println("Config.ColNames", Config.ColNames)
	//fmt.Println("Index des gènes :", geneIndex)
	return nil
}
