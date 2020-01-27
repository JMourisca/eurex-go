package main

import (
	"eurex-juliana/converter"
	"flag"
	"fmt"
)

func main() {

	//queuePtr := flag.String("queue", "all", "queue name: a string")
	sourcePtr := flag.String("source", "", "Source Currency")
	targetPtr := flag.String("target", "", "Target Currency")
	datePtr := flag.String("date", "", "Date YYYY-MM-DD")

	flag.Usage = func() {
		fmt.Println("----------------------")
		fmt.Println("- Currency Conversor -")
		fmt.Println("----------------------")
		fmt.Println("Convert currencies")
		fmt.Println("To run, just call")
		fmt.Printf("\t- go run main.go -source BRL -target EUR -date YYYY-MM-DD\n")
		fmt.Printf("\t- or just go run main.go to input the data\n")
		fmt.Println("")
	}

	flag.Parse()

	var date, source, target string

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

	converter.Convert(source, target, date)
}