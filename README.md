# rebalance.go

Rebalance your portfolio

rebalance.go tells you what stocks to buy when investing more money into
a portfolio using an asset allocation strategy

## Usage

rebalance prints information about its decisions to standard error and
only prints the new allocation to standard out.

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
    Prices map[SCHC:34.78 SCHV:43.25 SCHB:48.26 SCHP:55.1]
    Starting allocation: map[SCHC:10 SCHB:300 SCHV:50 SCHP:25]
    Final allocation: map[SCHC:59 SCHB:302 SCHV:50 SCHP:38]
    Final percentages
    SCHB 69.79184823520173
    SCHV 10.355392274230901
    SCHP 10.026414031807937
    SCHC 9.826345458759441
    Buys to make
    SCHP 13
    SCHC 49
    SCHB 2
    Total Investment: 20882.839999999997
    New allocation
    {"SCHB":302,"SCHC":59,"SCHP":38,"SCHV":50}

You can pipe standard output to save the new allocations and feed into the next round of rebalancing.

    $ rebalance -amount=2500 > new-allocation.json
    $ rebalance -amount=2500 -current=new-allocation.json

Or pipe it ftw

    $ cat new-allocation.json | rebalance -amount=2500
