package main

import (
	"fmt"
	"time"
)

// Example Transaction
// Date 11/13/2023
// Amount 12.00
// Place: Amazon
// CharacterPatterns are the patterns of characters that are associated with the Place
// Transaction at place Amazon the character pattern could be: [AMAZON.COM*TO3Q13XI0, AMZNAMZN.COM/BILLWA]
type Transaction struct {
	Date                     time.Time
	Amount                   float64
	Place                    string
	WordsAssociatedWithPlace []string
}

const other string = "other"

// ignore these since they are payments to the credit card and don't represent actual transactions outside of user bank account
var ignorePayment1 = []string{"online", "payment", "thank", "you"}

var ignorePayment2 = []string{"automatic", "payment", "thank", "you"}

func (t *Transaction) printTransaction() {
	fmt.Println("Place: ", t.Place)
	fmt.Println("CharacterPatterns", t.WordsAssociatedWithPlace, "\n Date: ", t.Date.Day(), t.Date.Month(), t.Date.Year(), "\n Amount", t.Amount)
}
