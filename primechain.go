package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	FLAGS "github.com/jessevdk/go-flags"
	block "github.com/weberr13/blockchain/block"
)

type Opts struct {
	Filename string `short:"f" description:"Filename of chain db" required:"true"`
}

var opts Opts
var parser = FLAGS.NewParser(&opts, FLAGS.Default)

type PrimeNum struct {
	Num int64
}

func main() {
	_, err := parser.Parse()
	if nil != err {
		panic(err)
	}
	p := block.NewSha256Pow(2)

	c := block.NewBlockChain(p, opts.Filename)
	if c == nil {
		panic("Couldn't build or open chain")
	}

	i := c.Iterator()
	knownPrimes := []int64{}
	oneNum := &PrimeNum{}
	for b, err := i.Next(); b != nil && err == nil && !b.IsGenesis(); b, err = i.Next() {
		proof := p.GetPOW(b)
		if !proof.Validate() {
			panic("Blockchain corrupt!")
		}
		err = json.Unmarshal(b.Data, oneNum)
		if err != nil {
			break
		}
		knownPrimes = append(knownPrimes, oneNum.Num)
	}
	if len(knownPrimes) <= 1 {
		knownPrimes = append(knownPrimes, 2)
		oneNum.Num = 2
		by, err := json.Marshal(oneNum)
		if err != nil {
			panic(fmt.Sprintf("json error:%v", err))
		}
		err = c.AddBlock(by)
		if err != nil {
			panic(fmt.Sprintf("addblock error: %v", err))
		}
		knownPrimes = append(knownPrimes, 3)
		oneNum.Num = 3
		by, err = json.Marshal(oneNum)
		if err != nil {
			panic(fmt.Sprintf("json error:%v", err))
		}
		err = c.AddBlock(by)
		if err != nil {
			panic(fmt.Sprintf("addblock error: %v", err))
		}
	}

	err = c.Close()
	if err != nil {
		panic(fmt.Sprintf("Could not close chain: %v", err))
	}
	trapSignal := make(chan os.Signal, 1)
	signal.Notify(trapSignal, os.Interrupt, syscall.SIGTERM)
hunt:
	for nextPrime := knownPrimes[0] + 2; ; nextPrime += 2 {
	check:
		for _, p := range knownPrimes {
			if int(nextPrime)%int(p) == 0 {
				select {
				case <-trapSignal:
					break hunt
				default:
					break
				}
				continue hunt
			}
			if p > nextPrime/2 {
				continue check
			}
		}
		c = block.NewBlockChain(p, opts.Filename)
		oldMax, err := c.Iterator().Next()
		if err != nil {
			panic(fmt.Sprintf("iterator error: %v", err))
		}
		proof := p.GetPOW(oldMax)
		if !proof.Validate() {
			panic("Blockchain corrupt!")
		}
		err = json.Unmarshal(oldMax.Data, oneNum)
		if oneNum.Num >= nextPrime {
			knownPrimes = append(knownPrimes, oneNum.Num)
			catchupIterator := c.Iterator()
		catchup:
			for oldMax, err := catchupIterator.Next(); oldMax != nil && err == nil && !oldMax.IsGenesis(); oldMax, err = catchupIterator.Next() {
				proof := p.GetPOW(oldMax)
				if !proof.Validate() {
					panic("Blockchain corrupt!")
				}
				err := json.Unmarshal(oldMax.Data, oneNum)
				if err != nil {
					panic(fmt.Sprintf("json error:%v", err))
				}
				if oneNum.Num < nextPrime {
					break catchup
				}
				knownPrimes = append(knownPrimes, oneNum.Num)
			}
			if err != nil {
				panic(fmt.Sprintf("iterator error: %v", err))
			}
			err = c.Close()
			if err != nil {
				panic(fmt.Sprintf("Could not close chain: %v", err))
			}
			continue hunt
		}
		knownPrimes = append(knownPrimes, nextPrime)
		oneNum.Num = nextPrime
		by, err := json.Marshal(oneNum)
		if err != nil {
			panic(fmt.Sprintf("json error:%v", err))
		}
		err = c.AddBlock(by)
		if err != nil {
			panic(fmt.Sprintf("addblock error: %v", err))
		}

		err = c.Close()
		if err != nil {
			panic(fmt.Sprintf("Could not close chain: %v", err))
		}
		select {
		case <-trapSignal:
			break hunt
		default:
			break
		}
	}
	fmt.Println("Here are the primes: ", knownPrimes)

}
