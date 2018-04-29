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
	timeout := time.NewTimer(*dur)
	var i, rps int

	for {
		select {
		case <-timeout.C:
			fmt.Println("Timeout. Test was repeated %d times. Duration %v. %vrps.\n", i, *dur, float64(rps)/dur.Seconds())
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

				ans := strings.TrimFunc(string(by), func(r rune) bool {
					return !unicode.IsNumber(r) && !unicode.IsSpace(r)
				})
				if fi := strings.Fields(ans); fi[0] != v.v1 || fi[1] != v.v2 && v.v1 != "wa" {
					if !timeout.Stop() {
						to := <-timeout.C
						fmt.Printf("Wrong answer after %v.\nTest %d: %q.\nServer answer is: %q.\nRight answer is: %q.\n", to, j, v.v1, string(by), v.v2)
						return
					}
				}
				i++
			}
		}
	}
}
