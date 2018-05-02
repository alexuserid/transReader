package main

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"flag"
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

var (
	transMap = make(map[key]bt)
	port     = flag.String("port", ":8080", "Server port")
	file     = flag.String("file", "tr.txt.gz", "File with transactions in .gzip format")
)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	hexTrans, err := hex.DecodeString(trans)
	if err != nil {
		http.Error(w, "The request is not in a hex format. Wrong key.", http.StatusBadRequest)
		return
	}

	var k32b key
	copy(k32b[:], hexTrans)

	v, ok := transMap[k32b]
	if !ok {
		http.Error(w, "Wrong key", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Json error", http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	log.Println("Reading data...")
	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("gzip.NewReader: %v", err)
	}

	var (
		k32byte key
		k       string
		v1      uint32
		v2      uint16
	)

	for {
		_, err := fmt.Fscanln(gz, &k, &v1, &v2)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("fmt.Fscanln: %v", err)
		}
		hexKey, err := hex.DecodeString(k)
		if err != nil {
			log.Fatalf("hex.Decode: %v", err)
		}
		copy(k32byte[:], hexKey)
		transMap[k32byte] = bt{v1, v2}
	}
	f.Close()
	log.Println("Success. Launch server.")

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
