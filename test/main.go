package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

var (
	addr = flag.String("a", "http://localhost:8080", "Server address for test")
	dur  = flag.Duration("d", 5*time.Second, "Test duration")
)

func main() {
	flag.Parse()
	timer := time.NewTimer(*dur)
	var i float64

	for {
		select {
		case <-timer.C:
			fmt.Printf("Timeout. Test was repeated %g times. Duration %v. %vrps.\n", i, *dur, i/dur.Seconds())
			return
		default:
			for j, v := range mass {
				resp, err := http.Get(*addr + "?t=" + v.k)
				if err != nil {
					log.Fatalf("http.Get: %v", err)
				}

				by, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("ioutil.ReadAll: %v", err)
				}

				trims := strings.TrimFunc(string(by), func(r rune) bool {
					return !unicode.IsNumber(r)
				})
				sp := strings.Split(trims, `,"Tr":`)
				if (sp[0] != v.v1 || sp[1] != v.v2) && v.v1 != "wa" {
					fmt.Printf("Wrong answer.\nTest %d: %q.\nServer answer is: %q.\nRight answer is: %q.\n", j, v.v1, string(by), v.v2)
					return
				}
				i++
			}
		}
	}
}
