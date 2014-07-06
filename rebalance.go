package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"
  "strconv"
  "flag"
  "sync"
  "log"
  "os"
)

func checkTarget(target map[string]float64) bool {
  total := 0.0
  for _,allocation := range(target){
    total += allocation;
  }

  return total == 100;
}

func getPrices(stocks []string) map[string]float64{
  prices := make(map[string]float64);

  var wg sync.WaitGroup
  for _,_symbol := range(stocks) {
    wg.Add(1)
    go func(symbol string){
      defer wg.Done()

      resp, err := http.Get("http://www.google.com/finance/info?q=" + symbol)
      if(err != nil){
        log.Fatal(err)
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
    }(_symbol)
  }
  wg.Wait()
  return prices;
}

func parseAllocation(fileName string) (map[string]float64, error) {
  var allocations = make(map[string]float64)

  allocationData, err := ioutil.ReadFile(fileName);
  if(err != nil){
    return nil, err;
  }
  err = json.Unmarshal(allocationData, &allocations)
  if(err != nil){
    return nil, err;
  }

  return allocations, nil;
}

func balanceAllocations(investAmount int, currentAllocation, targetAllocation map[string]float64) {
  stocks := make([]string, 0, len(targetAllocation))
  for symbol := range(targetAllocation){
    stocks = append(stocks, symbol)
  }

  prices := getPrices(stocks)
  fmt.Fprintln(os.Stderr, "Prices", prices);
  buys   := make(map[string]int)
  amountInvested := 0.0

  fmt.Fprintln(os.Stderr, "Starting allocation:",currentAllocation)
  for {
    for symbol, allocation := range(currentAllocation) {
      currentValue := allocation * prices[symbol]
      totalValue := 0.0
      for s, a := range(currentAllocation) {
        totalValue += a * prices[s]
      }
      currentPercent := currentValue / totalValue;
      if(currentPercent <= targetAllocation[symbol] / 100){
        currentAllocation[symbol] += 1
        buys[symbol] += 1;
        amountInvested += prices[symbol];
        //fmt.Println(currentPercent, "less than", targetAllocation[symbol], "so buy 1", symbol)
      }
    }

    if amountInvested > float64(investAmount) {
      break
    }
  }

  totalValue := 0.0
  for symbol, allocation := range(currentAllocation) {
    totalValue += allocation * prices[symbol]
  }

  fmt.Fprintln(os.Stderr, "Final allocation:",currentAllocation)
  fmt.Fprintln(os.Stderr, "Final percentages")
  for symbol, allocation := range(currentAllocation) {
    value := allocation * prices[symbol]
    fmt.Fprintln(os.Stderr, symbol, value/totalValue*100)
  }

  fmt.Fprintln(os.Stderr, "Buys to make")
  for symbol, buys := range(buys) {
    fmt.Fprintln(os.Stderr, symbol, buys)
  }
  //fmt.Fprintln(os.Stderr, "Buys to make:",buys)
  fmt.Fprintln(os.Stderr, "Total Investment:",totalValue)

  fmt.Fprintln(os.Stderr, "New allocation")
  allocationJSON, err := json.Marshal(currentAllocation)
  if(err != nil){
    log.Fatal(err)
  }
  fmt.Println(string(allocationJSON))
}

func main() {
  var investAmount int;
  flag.IntVar(&investAmount, "amount", 1000, "amount to invest")
  var targetAllocationFile string;
  flag.StringVar(&targetAllocationFile, "target",
    "./target-allocation.json", "json file of stock: percent")
  var currentAllocationFile string;
  flag.StringVar(&currentAllocationFile, "current",
    "./current-allocation.json", "json file of stock: number of stocks")
  flag.Parse()

  targetAllocation, err := parseAllocation(targetAllocationFile);
  if err != nil {
    log.Fatal(err)
  }
  currentAllocation, err := parseAllocation(currentAllocationFile);
  if err != nil {
    log.Fatal(err)
  }

  if !checkTarget(targetAllocation) {
    log.Fatal("Target allocation does not add up to 100%")
  }
  balanceAllocations(investAmount, currentAllocation, targetAllocation)
}
