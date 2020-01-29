package main

import (
	"eurex-juliana/converter"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {

	// Set input fields and helper
	amountPtr := flag.String("amount", "", "Amount (e.g. 10)")
	sourcePtr := flag.String("source", "", "Source Currency (e.g. BRL)")
	targetPtr := flag.String("target", "", "Target Currency (e.g. CHF)")
	datePtr := flag.String("date", "", "Date (e.g. YYYY-MM-DD)")

	flag.Usage = func() {
		fmt.Println("----------------------")
		fmt.Println("- Currency Conversor -")
		fmt.Println("----------------------")
		fmt.Println("Convert currencies")
		fmt.Println("To run, call")
		fmt.Printf("\t- go run main.go -amount 10 -source BRL -target EUR -date YYYY-MM-DD\n")
		fmt.Printf("\t- or just go run main.go to input the data\n")
		fmt.Println("")
	}

	flag.Parse()

	var date, source, target, amountStr string

	if *amountPtr != "" {
		amountStr = *amountPtr
	} else {
		var amountInput string
		fmt.Print("Amount: ")
		_, _ = fmt.Scanln(&amountInput)
		amountStr = amountInput
	}

	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil {
		log.Error("Invalid amount.")
		os.Exit(1)
	}

	if *sourcePtr != "" {
		source = *sourcePtr
	} else {
		var sourceInput string
		fmt.Print("Source Currency: ")
		_, _ = fmt.Scanln(&sourceInput)
		source = sourceInput
	}

	if *targetPtr != "" {
		target = *targetPtr
	} else {
		var targetInput string
		fmt.Print("Target Currency: ")
		_, _ = fmt.Scanln(&targetInput)
		target = targetInput
	}

	if *datePtr != "" {
		date = *datePtr
	} else {
		var dateInput string
		fmt.Print("Date: ")
		_, _ = fmt.Scanln(&dateInput)
		date = dateInput
	}

	res := converter.Convert(source, target, date, amount)
	fmt.Printf("%.2f %s = %.2f %s in %s\n", amount, source, res, target, date)
}