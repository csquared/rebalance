# rebalance.go

Rebalance your portfolio

rebalance.go tells you what stocks to buy when investing more money into
a portfolio using an asset allocation strategy

## Usage

### options

* `amount`: amount to invest in USD
* `current`: json file of current investments in number of
  stocks (defaults to `./current-allocation.json`)
* `target`: json file of target allocation in whole percentages
  (defaults to `./target-allocation.json`)

### example

    $ cat current-allocation.json
    {
      "SCHP" : 25,
      "SCHC" : 10,
      "SCHB" : 300,
      "SCHV" : 50
    }

    $ cat target-allocation.json
    {
      "SCHP" : 10,
      "SCHC" : 10,
      "SCHB" : 70,
      "SCHV" : 10
    }

    $ rebalance -amount=2500
    Prices map[SCHV:43.25 SCHB:48.26 SCHC:34.78 SCHP:55.1]
    Starting allocation: map[SCHC:10 SCHB:300 SCHV:50 SCHP:25]
    Final allocation: map[SCHP:38 SCHC:59 SCHB:303 SCHV:50]
    Final percentages
    SCHC 9.803689247101204
    SCHB 69.86149796236222
    SCHV 10.33151626049276
    SCHP 10.003296530043812
    Buys to make
    SCHB 3
    SCHP 13
    SCHC 49
    Total Investment: 20931.1

    $ rebalance -amount=5000
    Prices map[SCHP:55.1 SCHC:34.78 SCHV:43.25 SCHB:48.26]
    Starting allocation: map[SCHP:25 SCHC:10 SCHB:300 SCHV:50]
    Final allocation: map[SCHP:43 SCHC:68 SCHB:339 SCHV:55]
    Final percentages
    SCHP 10.09362580266968
    SCHC 10.07547746944072
    SCHB 69.69701229869088
    SCHV 10.133884429198709
    Buys to make
    SCHP 18
    SCHC 58
    SCHB 39
    SCHV 5
    Total Investment: 23473.23
