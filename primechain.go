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

	knownPrimes := []int64{}
	fmt.Println("walking old data")
	err = c.Walk(func(b *block.Block) error {
		oneNum := &PrimeNum{}
		err = json.Unmarshal(b.Data, oneNum)
		if err != nil {
			fmt.Println(b.Headers)
			fmt.Printf(string(b.Data))
			return err
		}
		knownPrimes = append([]int64{oneNum.Num}, knownPrimes...)
		return nil
	}, func(b *block.Block) bool {
		// Want the whole thing
		return false
	})
	if err != nil {
		panic(fmt.Sprintf("Could not read data: ", err))
	}
	oneNum := &PrimeNum{}

	if len(knownPrimes) <= 1 {
		fmt.Println("seeding the sequence")

		knownPrimes = append(knownPrimes, 2)
		knownPrimes = append(knownPrimes, 3)

		oneNum.Num = 2
		by, err := json.Marshal(oneNum)
		if err != nil {
			panic(fmt.Sprintf("json error:%v", err))
		}
		err = c.AddBlock(by)
		if err != nil {
			panic(fmt.Sprintf("addblock error: %v", err))
		}
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
	fmt.Println("starting computation")

	if err != nil {
		panic(fmt.Sprintf("Could not close chain: %v", err))
	}
	trapSignal := make(chan os.Signal, 1)
	signal.Notify(trapSignal, os.Interrupt, syscall.SIGTERM)
hunt:
	for nextPrime := knownPrimes[len(knownPrimes)-1] + 2; ; nextPrime += 2 {
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
			newPrimes := []int64{oneNum.Num}

			err = c.Walk(func(b *block.Block) error {
				err = json.Unmarshal(b.Data, oneNum)
				if err != nil {
					return err
				}
				newPrimes = append([]int64{oneNum.Num}, newPrimes...)
				return nil
			}, func(b *block.Block) bool {
				return oneNum.Num < nextPrime
			})
			if err != nil {
				panic(fmt.Sprintf("iterator error: %v", err))
			}
			err = c.Close()
			if err != nil {
				panic(fmt.Sprintf("Could not close chain: %v", err))
			}
			knownPrimes = append(knownPrimes, newPrimes...)
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
	old := int64(0)
	for _, prime := range knownPrimes {
		if old > prime {
			fmt.Println(old, " > ", prime)
			panic("Out of order, something is terribly wrong")
		}
		old = prime
	}
	fmt.Println("Here are the primes: ", knownPrimes)

}
