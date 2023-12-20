# Finance Helper

## Purpose
Helps user see where their money is being spent based on CSV file uploaded.
CSV rows formatted as:Date|amount|description.
User inputs start date and end date they would like to see transactions for. Application outputs how much was spent where based on the
values of objects in placesWithPatterns.txt 

## How to Use
Install latest version of go
copy csv file into project directory
name csv file "CreditCard3"
run "go run main.go" 