# eurex-juliana

## Introduction
This program fetches information from https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml, saves 
into a file and uses this file as source. If the file is more than 24h old, it fetches a new version from 
https://www.ecb.europa.eu.

## Usage
```bash
$ go run main.go -amount 10 -source BRL -target EUR -date 2020-01-24
10.00 BRL = 2.32 CHF in 2020-01-24
```
or simply ```go run main.go``` to insert the inputs one by one as asked.

When an source or target don't exist, it defaults to EUR. The amount and date must be correct

To see the help:
```bash
go run main.go -h 
```

## To Do
* If a date doesn't exist, use the latest one;
* Improve tests;
* Improve data input to be more flexible with formats;
* When reading data from the site, append to the existing data instead of replacing it in order to have a source with 
more than 90 days;
* Maybe use golang.org/x/text/currency to handle the currency;

## Task
Write a Go library to exchange money from one currency into another, using the ECB reference exchange rate for a particular day (within the last 90 days)

Example: convert 14.50 USD to CHF on Dec 2nd

ECB reference rates as XML files https://www.ecb.europa.eu/stats/exchange/eurofxref/html/index.en.html

90 days history https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml

## Timeline
This task should take you a couple of hours. There is no time pressure. Just send me an email when you are done.

## Expectations
We are going to look at the correctness and style of your code and your choice of data structures. Please also write some tests and an imaginary todo list of what could be changed or added in the future.
