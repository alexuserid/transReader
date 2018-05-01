package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"os"
	"math/rand"
	"compress/gzip"
	"encoding/json"
)

var (
	addr = flag.String("address", "http://localhost:8080", "Server address for test")
	dur  = flag.Duration("duration", 5*time.Second, "Test duration")
	testFile = flag.String("test", "../tr.txt.gz", "Test file")
	n = flag.Int("number", 10, "Test items number")
)

type st struct {
	k string
	v1, v2 int
}
type bt struct {
	Block,Tr int
}

func shuf(n *int, testFile *string) []st {
	openf, err := os.Open(*testFile)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer openf.Close()

	gzReader, err := gzip.NewReader(openf)
	if err != nil {
		log.Fatalf("gzip.NewReader: %v", err)
	}

	var (
		key string
		val1, val2, i, j int
		randMass []st
	)
	for {
		_, err := fmt.Fscanln(gzReader, &key, &val1, &val2)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("fmt.Fscanln: %v", err)
		}

		if i < *n {
			randMass = append(randMass, st{key, val1, val2})
		} else {
			j = rand.Intn(i)
			if j < *n {
				randMass[j] = st{key, val1, val2}
			}
		}
		i++
	}
	for _, value := range wrongMass {
		j = rand.Intn(i)
		if j < *n {
			randMass = append(randMass, value)
		}
		i++
	}
	return randMass
}

func main() {
	flag.Parse()

	mix := shuf(n, testFile)

	timer := time.NewTimer(*dur)
	var i int

	for {
		select {
		case <-timer.C:
			fmt.Printf("Timeout. Test was repeated %d times. Duration %v. %vrps.\n", i, *dur, float64(i)/dur.Seconds())
			return
		default:
			for j, v := range mix {
				resp, err := http.Get(*addr + "?t=" + v.k)
				if err != nil {
					log.Fatalf("http.Get: %v", err)
				}

				var dec bt
				json.NewDecoder(resp.Body).Decode(&dec)

				if (dec.Block != v.v1 || dec.Tr != v.v2) && v.v1 != 0 {
					fmt.Printf("Wrong answer.\nTest %d: %q.\nServer answer is: %q.\nRight answer is: %q.\n", j, v.v1, dec, v.v2)
					return
				}
				i++
			}
		}
	}
}
