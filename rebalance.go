package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func checkTarget(target map[string]float64) bool {
	total := 0.0
	for _, allocation := range target {
		total += allocation
	}

	return total == 100
}

type StockQuote struct {
	Price string `json:"l"`
}

func getPrices(stocks []string) map[string]float64 {
	prices := make(map[string]float64)

	var wg sync.WaitGroup
	for _, _symbol := range stocks {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()

			resp, err := http.Get("http://www.google.com/finance/info?q=" + symbol)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			bodyString := strings.TrimLeft(string(body), "/ \n")
			var quotes []StockQuote
			err = json.Unmarshal([]byte(bodyString), &quotes)
			if err != nil {
				log.Fatal(err)
			}

			prices[symbol], err = strconv.ParseFloat(quotes[0].Price, 64)
		}(_symbol)
	}
	wg.Wait()
	return prices
}

func parseAllocation(allocationData []byte) (map[string]float64, error) {
	var allocations = make(map[string]float64)
	err := json.Unmarshal(allocationData, &allocations)
	if err != nil {
		return nil, err
	}
	return allocations, nil
}

func readAllocation(fileName string) (map[string]float64, error) {
	allocationData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return parseAllocation(allocationData)
}

func readAllocationStdin() (map[string]float64, error) {
	allocationData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return parseAllocation(allocationData)
}

func balanceAllocations(investAmount int, currentAllocation, targetAllocation map[string]float64) {
	stocks := make([]string, 0, len(targetAllocation))
	for symbol := range targetAllocation {
		stocks = append(stocks, symbol)
	}

	prices := getPrices(stocks)
	fmt.Fprintln(os.Stderr, "Prices", prices)
	buys := make(map[string]int)
	amountInvested := 0.0

	fmt.Fprintln(os.Stderr, "Starting allocation:", currentAllocation)
	for {
		for symbol, allocation := range currentAllocation {
			currentValue := allocation * prices[symbol]
			//calculate total value each round b.c it changes
			totalValue := 0.0
			for symbol, allocation := range currentAllocation {
				totalValue += allocation * prices[symbol]
			}
			//if our current percent is under, buy one of these
			currentPercent := currentValue / totalValue
			if currentPercent <= targetAllocation[symbol]/100 {
				currentAllocation[symbol] += 1
				buys[symbol] += 1
				amountInvested += prices[symbol]
			}
		}

		if amountInvested > float64(investAmount) {
			break
		}
	}

	totalValue := 0.0
	for symbol, allocation := range currentAllocation {
		totalValue += allocation * prices[symbol]
	}

	fmt.Fprintln(os.Stderr, "Final allocation:", currentAllocation)
	fmt.Fprintln(os.Stderr, "Final percentages")
	for symbol, allocation := range currentAllocation {
		value := allocation * prices[symbol]
		fmt.Fprintln(os.Stderr, symbol, value/totalValue*100)
	}

	fmt.Fprintln(os.Stderr, "Buys to make")
	for symbol, buys := range buys {
		fmt.Fprintln(os.Stderr, symbol, buys)
	}
	fmt.Fprintln(os.Stderr, "Total Investment:", totalValue)

	fmt.Fprintln(os.Stderr, "New allocation")
	allocationJSON, err := json.Marshal(currentAllocation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stdout, string(allocationJSON))
}

func main() {
	var investAmount int
	flag.IntVar(&investAmount, "amount", 1000, "amount to invest")
	var targetAllocationFile string
	flag.StringVar(&targetAllocationFile, "target",
		"./target-allocation.json", "json file of stock: percent")
	var currentAllocationFile string
	flag.StringVar(&currentAllocationFile, "current",
		"./current-allocation.json", "json file of stock: number of stocks")
	flag.Parse()

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var currentAllocation map[string]float64
	if fi.Size() > 0 {
		currentAllocation, err = readAllocationStdin()
	} else {
		currentAllocation, err = readAllocation(currentAllocationFile)
	}

	if err != nil {
		log.Fatal(err)
	}

	targetAllocation, err := readAllocation(targetAllocationFile)
	if err != nil {
		log.Fatal(err)
	}

	if !checkTarget(targetAllocation) {
		log.Fatal("Target allocation does not add up to 100%")
	}
	balanceAllocations(investAmount, currentAllocation, targetAllocation)
}
