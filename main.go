package main

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type bt struct {
	Block uint32
	Tr    uint16
}
type key [32]byte

var m = make(map[key]bt)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	hTrans, err := hex.DecodeString(trans)
	if err != nil {
		log.Fatalf("hex.DecodeString: %v", err)
	}
	var k key
	copy(k[:], hTrans)

	v, ok := m[k]
	if !ok {
		w.Write([]byte("Wrong key\n"))
		return
	}
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Fatalf("json.Encode: %v", err)
	}
	w.Write([]byte("\n"))
}

func main() {
	f, err := os.Open("tr.txt.gz")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("gzip.NewReader: %v", err)
	}

	var (
		kk key
		k  string
		v1 uint32
		v2 uint16
	)

	for {
		_, err := fmt.Fscanln(gz, &k, &v1, &v2)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("fmt.Fscanln: %v", err)
		}
		hk, err := hex.DecodeString(k)
		if err != nil {
			log.Fatalf("hex.Decode: %v", err)
		}
		copy(kk[:], hk)
		m[kk] = bt{v1, v2}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
