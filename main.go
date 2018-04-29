package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
)

type bt struct {
	Block uint32
	Tr uint16
}
type b32 struct {
	bs []byte
}

var m = make(map[b32]bt)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	var x = b32{[]byte(trans)}
	j, err := json.Marshal(m[x])
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
		k b32
		v1 uint32
		v2 uint16
	)

	for {
		_, err := fmt.Fscanln(gz, &k, &v1, &v2)
		if err != nil {
			fmt.Println("fmt.Fscanln: ", err)
			break
		}
		m[k] = bt{v1, v2}
	}
	fmt.Println(m)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
