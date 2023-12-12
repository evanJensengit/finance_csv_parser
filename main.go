package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Define a struct for your CSV data
type Record struct {
	Date   time.Time
	Amount float64
	Place  string
}

// generatePlaceDict reads a text file with lines of strings
// in the form of "key":"value" and returns a map with key-value pairs.
func generatePlaceMap(filename string) (map[string]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialize a map to store key-value pairs
	placeDict := make(map[string]string)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through each line
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by ","
		parts := strings.Split(line, ",")
		for _, placeAndPattern := range parts {

			keyValuePair := strings.Split(placeAndPattern, ":")
			if len(keyValuePair) == 2 {
				// Trim leading and trailing spaces and commas from key and value
				key := strings.TrimSpace(keyValuePair[0])
				value := strings.TrimSpace(keyValuePair[1])

				// Add key-value pair to the map
				placeDict[key] = value
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return placeDict, nil
}

func generateTransactionsAtPlaces(keywordMapToPlaces map[string]string) (map[string]int, error) {
	transactionsAtPlaces := make(map[string]int)
	for _, value := range keywordMapToPlaces {
		_, alreadyInMap := transactionsAtPlaces[value]
		if !alreadyInMap {
			transactionsAtPlaces[value] = 0
		}
	}
	return transactionsAtPlaces, nil
}

func main() {
	placesWithPatterns := "placesWithPatterns.txt"
	keywordMapToPlaces, err := generatePlaceMap(placesWithPatterns)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Keyword map to places")

	for key, val := range keywordMapToPlaces {
		fmt.Println(key, ":", val)
	}
	transactionsAtPlacesMap, err := generateTransactionsAtPlaces(keywordMapToPlaces)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("TransactionsAtPlacesMap")
	for key, val := range transactionsAtPlacesMap {
		fmt.Println(key, ":", val)
	}
	// Open the CSV file
	file, err := os.Open("CreditCard3.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Create a slice to store the records
	var data []Record

	// Iterate through the CSV records and convert them to the struct
	for _, row := range records[0:] { // Assuming the first row contains headers
		date, err := time.Parse("01/02/2006", row[0])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println("Error parsing amount:", err)
			return
		}
		fmt.Println(row[4])
		record := Record{
			Date:   date,
			Amount: amount,
			Place:  row[4],
		}

		data = append(data, record)

	}

	// Print the resulting struct
	fmt.Println(data)
}
