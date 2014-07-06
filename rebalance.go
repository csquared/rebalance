package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"
  "strconv"
)

func getPrices(stocks []string) map[string]float64{
  prices := make(map[string]float64);

  for _,symbol := range(stocks) {
    resp, err := http.Get("http://www.google.com/finance/info?q=" + symbol)
    if(err != nil){
      panic("HTTP price lookup failed")
    }
    defer resp.Body.Close()

    //this is ugly
    body, err := ioutil.ReadAll(resp.Body)
    bodyString := strings.TrimLeft(string(body), "/ \n")
    var f interface{}
    err = json.Unmarshal([]byte(bodyString), &f)
    m := f.([]interface{})
    first := m[0].(map[string]interface{})
    prices[symbol], err = strconv.ParseFloat(first["l"].(string), 64)
  }
  return prices;
}

var targetAllocation = map[string]float64 {
  "SCHP" : 10,
  "SCHC" : 10,
  "SCHB" : 70,
  "SCHV" : 10,
}

var currentAllocation = map[string]float64 {
  "SCHP" : 16.0766,
  "SCHC" : 31,
  "SCHB" : 305.7877,
  "SCHV" : 48.2886,
}

func main() {
  stocks := make([]string, 0, len(targetAllocation))
  for symbol := range(targetAllocation){
    stocks = append(stocks, symbol)
  }
  prices := getPrices(stocks)
  fmt.Println(prices);
  buys   := make(map[string]int)
  amountInvested := 0.0
  investLimit := 2500.0

  for {
    for symbol, allocation := range(currentAllocation) {
      currentValue := allocation * prices[symbol]
      totalValue := 0.0
      for s, a := range(currentAllocation) {
        totalValue += a * prices[s]
      }
      currentPercent := currentValue / totalValue;
      if(currentPercent <= targetAllocation[symbol] / 100){
        fmt.Println(currentPercent, "less than", targetAllocation[symbol], "so buy 1", symbol)
        currentAllocation[symbol] += 1
        buys[symbol] += 1;
        amountInvested += prices[symbol];
      }
  //    fmt.Println(symbol, allocation, currentValue, currentPercent);
    }

    if amountInvested > investLimit {
      break
    }
  }

  totalValue := 0.0
  for s, a := range(currentAllocation) {
    totalValue += a * prices[s]
  }

  fmt.Println("Final percentages")
  for symbol, allocation := range(currentAllocation) {
    value := allocation * prices[symbol]
    fmt.Println(symbol, value/totalValue*100)
  }

  fmt.Println("Final allocation:",currentAllocation)
  fmt.Println("Buys to make:",buys)
  fmt.Println("Total Investment:",totalValue)
}
