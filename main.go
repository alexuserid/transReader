package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
)

type bt struct {
	Block string
	Tr string
}

var m = make(map[string]bt)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	j, err := json.Marshal(m[trans])
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

	b, err := ioutil.ReadAll(gz)
	if err != nil {
		fmt.Println("ioutil.ReadAll: ", err)
	}

	fieldsB := bytes.Fields(b)
	for i:=0; i<len(fieldsB); i+=3 {
		m[string(fieldsB[i])] = bt{string(fieldsB[i+1]), string(fieldsB[i+2])}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
