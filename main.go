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

var monthsOfYear = map[string]int{"January": 1, "February:": 2, "March": 3, "April": 4, "May": 5, "June": 6, "July": 7,
	"Auguest": 8, "September": 9, "October": 10, "November": 11, "December": 12}

const other string = "other"

var ignorePayment = []string{"online", "payment", "thank", "you"}

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

func firstDateLessThanSecondDate(date1 time.Time, date2 time.Time) (bool, error) {
	fmt.Println("First date: ", date1, "\n", "Second Date: ", date2)
	if date1.Year() < date2.Year() {
		return true, nil
	} else if date1.Year() == date2.Year() {
		if int(date1.Month()) < int(date2.Month()) {
			return true, nil
		} else if int(date1.Month()) == int(date2.Month()) {
			if date1.Day() < date2.Day() {
				return true, nil
			}
		}
	}
	return false, nil
}

func findTransactionRangeToCalculate(transactionsAtPlaces []Transaction) ([]Transaction, error) {
	lastTransactionDate := transactionsAtPlaces[0].Date
	firstTransactionDate := transactionsAtPlaces[len(transactionsAtPlaces)-1].Date

	fmt.Println("Please enter the dates you would like to calculate transactions ",
		"between the range of ", firstTransactionDate.Day(), firstTransactionDate.Month(), firstTransactionDate.Year(), " to",
		lastTransactionDate.Day(), lastTransactionDate.Month(), lastTransactionDate.Year(),
		"in the form of from: mm/dd/yyyy to: mm/dd/yyyy")

	//if fromInput <fromDate err
	//if toInput > toDate err

	// scanning the input by the user
	var fromInput, toInput string
	fmt.Scanln(&fromInput, &toInput)
	fmt.Println("from date", fromInput)
	fmt.Println("to date", toInput)

	fromInputDate, err := time.Parse("01/02/2006", fromInput)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return []Transaction{}, err
	}
	toInputDate, err := time.Parse("01/02/2006", toInput)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return []Transaction{}, err
	}
	inputChronological, err := firstDateLessThanSecondDate(fromInputDate, toInputDate)
	if err != nil {
		fmt.Println("Invalid date range entered:", err)
		return []Transaction{}, err
	}
	firstInputLessThanLastTransaction, err := firstDateLessThanSecondDate(fromInputDate, lastTransactionDate)
	if err != nil {
		fmt.Println("Invalid date range entered:", err)
		return []Transaction{}, err
	}

	fmt.Println("inputChronological", inputChronological, "firstInputLessThanLastTransaction", firstInputLessThanLastTransaction)
	//most reacent closer to front of life

	if inputChronological && firstInputLessThanLastTransaction {
		fmt.Println("valid dates entered")
		back := len(transactionsAtPlaces) - 1
		front := 0
		for back > front {
			inputGreaterThanCurrent, err := firstDateLessThanSecondDate(transactionsAtPlaces[back].Date, fromInputDate)
			if err != nil {
				fmt.Println("Invalid date range entered:", err)
				return []Transaction{}, err
			}
			if inputGreaterThanCurrent {
				back--
			}

			inputLessThanCurrent, err := firstDateLessThanSecondDate(transactionsAtPlaces[front].Date, toInputDate)
			if err != nil {
				fmt.Println("Invalid date range entered:", err)
				return []Transaction{}, err
			}
			if !inputLessThanCurrent {
				front++
			}

			if !inputGreaterThanCurrent && inputLessThanCurrent {
				break
			}
			fmt.Println("Back ", back, " Front ", front)
		}
		return transactionsAtPlaces[front : back+1], nil
	}
	return []Transaction{}, nil

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
				transactionsAtPlacesMap[place] += transaction.Amount

			}
		}
		if !foundMatch {
			if len(transaction.CharacterPatterns) == len(ignorePayment) {
				ignoreString := strings.Join(ignorePayment, " ")
				patternsStrings := strings.Join(transaction.CharacterPatterns, " ")
				if ignoreString == patternsStrings {
					continue
				}

			}
			unmatchedTransactions = append(unmatchedTransactions, transaction.CharacterPatterns)
			transactionsAtPlacesMap[other] += transaction.Amount

			fmt.Println("Other Transaction", transaction)
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
				patternsOfPlace = append(patternsOfPlace, strings.ToLower(val))
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
		fmt.Println("val: \n CharacterPattern", val.CharacterPatterns, "\n Date: ", val.Date.Day(), val.Date.Month(), val.Date.Year(), "\n Amount", val.Amount)
	}
	return data, nil
	// Print the resulting struct
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
	// fmt.Println("TransactionsAtPlacesMap")
	// for key, val := range transactionsAtPlacesMap {
	// 	fmt.Println(key, ":", val)
	// }
	//logging
	listOfTransactions, err := createTransactionObjects()

	listOfTransactions, err = findTransactionRangeToCalculate(listOfTransactions)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Transactions: ")
	fmt.Println(listOfTransactions)
	unmatched := [][]string{}
	calculateTransactionsAtPlaces(transactionsAtPlacesMap, listOfTransactions, placesMappedToPatterns, unmatched)
	fmt.Println(transactionsAtPlacesMap)
	//fmt.Printf("%.2f", k) to get rounded 2 decimal places
	//fmt.Println(data)
	return
}
