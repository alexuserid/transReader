package main

import (
	"compress/gzip"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"os"
	"net/http"
)

type bt struct {
	Block uint32
	Tr uint16
}
type key [32]byte

var m = make(map[key]bt)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	hTrans, err := hex.DecodeString(trans)
	if err != nil {
		fmt.Println("hex.DecodeString: ", err)
	}
	var k key
	copy(k[:], hTrans)

	j, err := json.Marshal(m[k])
	if err != nil {
		fmt.Println("json.Marshal: ", err)
		return
	}
	j = append(j, '\n')
	w.Write(j)
}


func main() {
	f, err := os.Open("tr.txt.gz")
	if err != nil {
		fmt.Println("os.Open: ", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println("gzip.NewReader: ", err)
	}

	var (
		kk key
		k string
		v1 uint32
		v2 uint16
	)

	for {
		_, err := fmt.Fscanln(gz, &k, &v1, &v2)
		if err != nil {
			fmt.Println("fmt.Fscanln: ", err)
			break
		}
		hk, err := hex.DecodeString(k)
		if err != nil {
			fmt.Println("hex.Decode: ", err)
		}
		copy(kk[:], hk)
		m[kk] = bt{v1, v2}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
