# rebalance.go

Rebalance your portfolio

rebalance.go tells you what stocks to buy when investing more money into
a portfolio using an asset allocation strategy

## Usage

### options

* `amount`: amount to invest in USD
* `current`: json file of current investments (defaults to
  `./current-allocation.json`)
* `target`: json file of target allocation (defaults to
  `./target-allocation.json`)

### example

    $ rebalance -amount=2500
    Prices map[SCHC:34.78 SCHV:43.25 SCHB:48.26 SCHP:55.1]
    Starting allocation: map[SCHP:16.0766 SCHC:31 SCHB:305.7877
    SCHV:48.2886]
    Final allocation: map[SCHV:49.2886 SCHP:39.0766 SCHC:62 SCHB:309.7877]
    Final percentages
    SCHP 10.065277867638994
    SCHC 10.080420937794552
    SCHB 69.88900997114104
    SCHV 9.965291223425405
    Buys to make: map[SCHB:4 SCHP:23 SCHC:31 SCHV:1]
    Total Investment: 21391.567012
