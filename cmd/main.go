package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const debug bool = true

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

// keywordMapToPlaces contains words associated with places i.e "amz":"amazon"
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

// creates transactions objects for each row in csv file
func createTransactionObjects() ([]Transaction, error) {
	var pathToCSV string
	pathToCSV = "/Users/evanjensen/go/src/finance_csv_sorter/CreditCard3.csv"

	if !debug {
		fmt.Println("Please enter the path to the csv file you would like to use ")

		// scanning the input by the user
		fmt.Scanln(&pathToCSV)
	}
	// Open the CSV file
	file, err := os.Open(pathToCSV)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	transactions, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	// Create a slice to store the records
	var data []Transaction

	// Iterate through the CSV records and convert them to the struct
	for _, row := range transactions[0:] { // Assuming the first row contains headers
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

		stringsInDescription := strings.Fields(row[4])
		patternsOfPlace := []string{}

		for _, val := range stringsInDescription {

			//append each character that is not a number or special character
			//abcd456abc would return abcd,abc
			currentString := ""
			for _, character := range val {
				if unicode.IsLetter(character) {
					currentString += string(character)
				} else {

					if currentString != "" {
						//no one letter words associated with place
						if len(currentString) > 1 {
							patternsOfPlace = append(patternsOfPlace, strings.ToLower(currentString))
						}
						currentString = ""
					}
				}
			}
			if currentString != "" {
				if len(currentString) > 2 {
					patternsOfPlace = append(patternsOfPlace, strings.ToLower(currentString))
				}
			}
		}

		record := Transaction{
			Date:                     date,
			Amount:                   amount,
			WordsAssociatedWithPlace: patternsOfPlace,
		}

		data = append(data, record)

	}

	return data, nil
}

// receives input from user of what dates to calculate transactions between (inclusive)
func findTransactionRangeToCalculate(transactionsAtPlaces []Transaction) ([]Transaction, error) {
	lastTransactionDate := transactionsAtPlaces[0].Date
	firstTransactionDate := transactionsAtPlaces[len(transactionsAtPlaces)-1].Date
	fromInput := "12/01/2023"
	toInput := "12/31/2023"
	if !debug {
		fmt.Println("Please enter the dates you would like to calculate transactions ",
			"between the range of ", firstTransactionDate.Day(), firstTransactionDate.Month(), firstTransactionDate.Year(), " to",
			lastTransactionDate.Day(), lastTransactionDate.Month(), lastTransactionDate.Year(),
			"in the form of from: mm/dd/yyyy to: mm/dd/yyyy")

		// scanning the input by the user
		var fromInput, toInput string
		fmt.Scanln(&fromInput, &toInput)
	}
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

	if debug {
		fmt.Println("inputChronological", inputChronological, "firstInputLessThanLastTransaction", firstInputLessThanLastTransaction)
	}

	if inputChronological && firstInputLessThanLastTransaction {
		back := len(transactionsAtPlaces) - 1
		front := 0
		//can be optimized to be log(n) time instead of O(n)
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

		}
		return transactionsAtPlaces[front : back+1], nil
	}
	return []Transaction{}, nil
}

// checks if the first date passed to function is chronologically before the second date
// returns true if so false otherwise
func firstDateLessThanSecondDate(date1 time.Time, date2 time.Time) (bool, error) {
	//fmt.Println("First date: ", date1, "\n", "Second Date: ", date2)
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

// populates transactionsAtPlaces map from Transaction objects in listOfTransactions and uses wordsAssociatedWithPlaces
// to determine which transactions are correlated with which place in
func calculateTransactionsAtPlaces(transactionsAtPlacesMap map[string]float64,
	listOfTransactions []Transaction, wordsAssociatedWithPlaces map[string]string,
	unmatchedTransactions []Transaction) {

	for index, transaction := range listOfTransactions {
		foundMatch := false
		for _, pattern := range transaction.WordsAssociatedWithPlace {
			place, ok := wordsAssociatedWithPlaces[pattern]
			if ok {
				foundMatch = true
				transactionsAtPlacesMap[place] += transaction.Amount
				listOfTransactions[index].Place = place
			}
		}
		if !foundMatch {
			if len(transaction.WordsAssociatedWithPlace) == len(ignorePayment1) {
				ignoreString1 := strings.Join(ignorePayment1, " ")
				ignoreString2 := strings.Join(ignorePayment2, " ")
				patternsStrings := strings.Join(transaction.WordsAssociatedWithPlace, " ")
				if ignoreString1 == patternsStrings || ignoreString2 == patternsStrings {
					continue
				}

			}
			//goes in "Other" category
			unmatchedTransactions = append(unmatchedTransactions, transaction)
			transactionsAtPlacesMap[other] += transaction.Amount
			listOfTransactions[index].Place = other
			if debug {
				fmt.Println("Other Transaction", transaction)
			}
		}
	}
}

// prints map given in alphabetical order
func printMapInOrder(myMap map[string]float64) {
	// Extract keys from the map
	var keys []string
	for key := range myMap {
		keys = append(keys, key)
	}

	// Sort the keys in alphabetical order
	sort.Strings(keys)

	// Iterate over sorted keys and print corresponding values
	var total float64 = 0
	for _, key := range keys {
		fmt.Printf("%s: %.2f\n", key, myMap[key])

		total += myMap[key]
	}

	fmt.Printf("Total: %.2f \n", total)
}

func getTotalSpentAtPlaces(transactionsAtPlaces map[string]float64) float64 {
	var total float64 = 0
	for key, _ := range transactionsAtPlaces {
		total += transactionsAtPlaces[key]
	}
	return math.Round(total*100) / 100
}
func saveMapToFile(transactionsAtPlaces map[string]float64) error {
	// Convert the map to JSON format

	jsonData, err := json.MarshalIndent(transactionsAtPlaces, "", "    ")
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	filename := currentDir + "/transactionsAtPlaces.txt"
	if err != nil {
		return err
	}
	// Write JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	totalSpentAtPlaces := make(map[string]float64)
	totalSpentAtPlaces["Total"] = getTotalSpentAtPlaces(transactionsAtPlaces)

	// Iterate through the map and append each key-value pair to the file
	for key, value := range totalSpentAtPlaces {
		line := fmt.Sprintf("\n%s: %.2f\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
func getTransactionsInCategory(listOfTransactions []Transaction, place string) []Transaction {
	var categoryTransactions []Transaction

	for _, t := range listOfTransactions {
		t.printTransaction()
		if t.Place == place {
			categoryTransactions = append(categoryTransactions, t)
		}
	}
	return categoryTransactions
}

func printKeysOfMapInOrder(transactionsAtPlacesMap map[string]float64) {
	var keys []string
	for key := range transactionsAtPlacesMap {
		keys = append(keys, key)
	}

	// Sort the keys in alphabetical order
	sort.Strings(keys)

	// Iterate over sorted keys and print corresponding values
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func loopThroughTransactionsInOther(listOfTransactions []Transaction, transactionsAtPlacesMap map[string]float64) {
	scanner := bufio.NewScanner(os.Stdin)
	otherTransactions := getTransactionsInCategory(listOfTransactions, "other")
	fmt.Println("loopThroughTransactionsInOther")

	for _, t := range otherTransactions {
		for {
			fmt.Println("\n(l) list current categories")
			fmt.Println("(a) add transaction to current category")
			fmt.Println("(n) add transaction to new category")
			fmt.Println("(s) skip this transaction")
			fmt.Println("(b) go back to previous menu ")

			fmt.Print("Transaction: \n")
			t.printTransaction()

			scanner.Scan()
			input := strings.ToLower(scanner.Text())
			if input == "l" {
				printKeysOfMapInOrder(transactionsAtPlacesMap)
			}
			if input == "a" {
				fmt.Println("Which place do you want to associate this transaction with?")
				printKeysOfMapInOrder(transactionsAtPlacesMap)
				scanner.Scan()
				input = strings.ToLower(scanner.Text())
				_, alreadyInMap := transactionsAtPlacesMap[input]
				if alreadyInMap {
					fmt.Println("Which term(s) from this transaction do you want to associate with", input)
					for _, val := range t.WordsAssociatedWithPlace {
						fmt.Print("\"", val, "\" ")
					}
					fmt.Println()

					scanner.Scan()
					input := strings.ToLower(scanner.Text())

					parts := strings.Split(input, " ")
					for _, word := range parts {
						fmt.Println(word)
					}
					break

				} else {
					fmt.Printf("Place %s is not recognized, would you like to add %s place to possible categories? (y/n) \n ", input, input)
					scanner.Scan()
					input := strings.ToLower(scanner.Text())
					if input == "y" {
						//create new category in wordsAssociatedWithTransactions
					} else {
						continue
					}
				}
				continue
			}
			if input == "n" {
				//TODO create func to ask user what new category should be called
				//and add transaction to new category
				continue
			}
			if input == "s" {
				//skip
				break
			}
			if input == "b" {
				//go back to previous menu
				return
			}
		}

	}
}

func main() {
	wordsAssociatedWithPlacesFile := "wordsAssociatedWithPlaces.txt"
	wordsAssociatedWithPlaces, err := generatePlaceMap(wordsAssociatedWithPlacesFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//logging
	if debug {
		fmt.Println("Keyword map to places")

		for key, val := range wordsAssociatedWithPlaces {
			fmt.Println(key, ":", val)
		}
	}

	//create transactions at places from values of wordsAssociatedWithPlaces map
	transactionsAtPlacesMap, err := initializeTransactionsAtPlacesMap(wordsAssociatedWithPlaces)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//if debug {
	// fmt.Println("TransactionsAtPlacesMap")
	// for key, val := range transactionsAtPlacesMap {
	// 	fmt.Println(key, ":", val)
	// }
	//}

	//read through csv file and create a list of transaction objects
	listOfTransactions, err := createTransactionObjects()

	//ask user what range they would like to calculate transactions for
	listOfTransactions, err = findTransactionRangeToCalculate(listOfTransactions)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//if debug {
	//fmt.Println("Transactions: ")

	// for _, val := range listOfTransactions {
	// 	val.printTransaction()
	// }
	//}

	var unmatched []Transaction

	calculateTransactionsAtPlaces(transactionsAtPlacesMap, listOfTransactions, wordsAssociatedWithPlaces, unmatched)

	if debug {
		printMapInOrder(transactionsAtPlacesMap)
	}

	err = saveMapToFile(transactionsAtPlacesMap)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter the letter with the action associated with what you want to do \n",
			"(a) look through the 'other' category\n",
			"(b) enter a new range of dates to calculate transactions\n",
			"(c) look through categories to see expenses in each category\n",
			"(q) quit the program\n")
		scanner.Scan()
		input := strings.ToLower(scanner.Text())
		fmt.Printf("User input %s %T\n", input, input)
		if input == "q" {
			fmt.Println("Exiting Finance Helper.")
			break
		}

		if input == "a" {
			loopThroughTransactionsInOther(listOfTransactions, transactionsAtPlacesMap)
		}

	}

	return
}
