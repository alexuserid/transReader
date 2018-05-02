package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type st struct {
	k      string
	v1, v2 int
	status int
}
type bt struct {
	Block, Tr int
}

var (
	addr     = flag.String("address", "http://localhost:8080", "Server address for test")
	dur      = flag.Duration("duration", 5*time.Second, "Test duration")
	testFile = flag.String("test", "../tr.txt.gz", "Test file")
	n        = flag.Int("number", 10, "Test items number")
)

func RASample(gzReader *gzip.Reader, n *int) ([]st, error) {
	var (
		key              string
		val1, val2, i, j int
		arr              []st
	)
	for {
		_, err := fmt.Fscanln(gzReader, &key, &val1, &val2)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if i < *n {
			arr = append(arr, st{key, val1, val2, http.StatusOK})
		} else {
			j = rand.Intn(i)
			if j < *n {
				arr[j] = st{key, val1, val2, http.StatusOK}
			}
		}
		i++
	}
	return arr, nil
}

func shuf(arr []st) []st {
	warr := wareturner()
	for _, v := range warr {
		arr = append(arr, v)
	}
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func main() {
	flag.Parse()
	log.Println("Reading data...")

	openf, err := os.Open(*testFile)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer openf.Close()
	gzReader, err := gzip.NewReader(openf)
	if err != nil {
		log.Fatalf("gzip.NewReader: %v", err)
	}
	arr, err := RASample(gzReader, n)
	if err != nil {
		log.Fatalf("fmt.Scanln: %v", err)
	}
	arr = shuf(arr)

	timer := time.NewTimer(*dur)
	log.Printf("Timer for %v is setted. Testing...\n", *dur)

	var i int
	var dec bt
	for {
		fmt.Print(i, "\r")
		select {
		case <-timer.C:
			fmt.Print(i, " requests\n")
			fmt.Printf("Timeout. Duration %v. %vrps.\n", *dur, float64(i)/dur.Seconds())
			return
		default:
			for j, v := range arr {
				resp, err := http.Get(*addr + "?t=" + v.k)
				if err != nil {
					log.Fatalf("http.Get: %v", err)
				}

				err = json.NewDecoder(resp.Body).Decode(&dec)
				if err != nil && resp.StatusCode != v.status {
					fmt.Printf("Wrong answer.\nTest %d: %q.\nServer answer is %q.\nTest answer is: %q.\n", j, v.v1, dec, http.StatusText(v.status))
					return
				}
				i++
			}
		}
	}
}
