package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Define a struct for your CSV data

// Example Transaction
// Date 11/13/2023
// Amount 12.00
// Place: Amazon
// CharacterPatterns are the patterns of characters that are associated with the Place
// Transaction at place Amazon the character pattern could be: [AMAZON.COM*TO3Q13XI0, AMZNAMZN.COM/BILLWA]
type Transaction struct {
	Date              time.Time
	Amount            float64
	Place             string
	CharacterPatterns []string
}

const other string = "other"

func generateTransactionsMap(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	transactionsAtPlaces := make(map[string]int)

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
				value, err := strconv.Atoi(keyValuePair[1])
				if err != nil {
					fmt.Println("Error converting to integer:", err)
					return nil, err
				}
				// Add key-value pair to the map
				transactionsAtPlaces[key] = value
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return transactionsAtPlaces, nil
}

// generatePlaceDict reads a text file with lines of strings
// in the form of ${key:value,key:value,key:value} and returns a map with key-value pairs.
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

func initializeTransactionsAtPlacesMap(keywordMapToPlaces map[string]string) (map[string]float64, error) {
	transactionsAtPlaces := make(map[string]float64)
	for _, value := range keywordMapToPlaces {
		_, alreadyInMap := transactionsAtPlaces[value]
		if !alreadyInMap {
			//initialize all transactions at places to 0
			transactionsAtPlaces[value] = 0
		}
	}
	return transactionsAtPlaces, nil
}

func main() {
	placesWithPatterns := "placesWithPatterns.txt"
	placesMappedToPatterns, err := generatePlaceMap(placesWithPatterns)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//logging
	fmt.Println("Keyword map to places")

	for key, val := range placesMappedToPatterns {
		fmt.Println(key, ":", val)
	}
	//logging
	transactionsAtPlacesMap, err := initializeTransactionsAtPlacesMap(placesMappedToPatterns)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//logging
	fmt.Println("TransactionsAtPlacesMap")
	for key, val := range transactionsAtPlacesMap {
		fmt.Println(key, ":", val)
	}
	//logging
	listOfTransactions, err := createTransactionObjects()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	unmatched := [][]string{}
	calculateTransactionsAtPlaces(transactionsAtPlacesMap, listOfTransactions, placesMappedToPatterns, unmatched)

	//fmt.Println(data)
}

func calculateTransactionsAtPlaces(transactionsAtPlacesMap map[string]float64,
	listOfTransactions []Transaction, placesMappedToPatterns map[string]string,
	unmatchedTransactions [][]string) {
	for _, transaction := range listOfTransactions {
		foundMatch := false
		for _, pattern := range transaction.CharacterPatterns {
			place, ok := placesMappedToPatterns[pattern]
			if ok {
				foundMatch = true
				transactionsAtPlacesMap[place] -= transaction.Amount
			}
		}
		if !foundMatch {
			unmatchedTransactions = append(unmatchedTransactions, transaction.CharacterPatterns)
			transactionsAtPlacesMap[other] -= transaction.Amount
		}
	}
}

// creates transactions objects for each row in csv file
func createTransactionObjects() ([]Transaction, error) {
	// Open the CSV file
	file, err := os.Open("CreditCard3.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	// Create a slice to store the records
	var data []Transaction

	// Iterate through the CSV records and convert them to the struct
	for _, row := range records[0:] { // Assuming the first row contains headers
		date, err := time.Parse("01/02/2006", row[0])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return nil, err
		}

		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println("Error parsing amount:", err)
			return nil, err
		}
		fmt.Println(row[4])
		stringsInDescription := strings.Fields(row[4])
		patternsOfPlace := []string{}
		for _, val := range stringsInDescription {
			if unicode.IsLetter(rune(val[0])) {
				patternsOfPlace = append(patternsOfPlace, val)
			}
		}
		//get all values that do not start with a number
		record := Transaction{
			Date:              date,
			Amount:            amount,
			CharacterPatterns: patternsOfPlace,
		}

		data = append(data, record)

	}
	for _, val := range data {
		fmt.Println("val: \n Place", val.Place, "\n Date: ", val.Date.Day(), val.Date.Month(), val.Date.Year(), "\n Amount", val.Amount)
	}
	return data, nil
	// Print the resulting struct
}
