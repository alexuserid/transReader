package main

import (
//	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
//	"strings"
)

type bt struct {
	block string
	tr string
}

var m = make(map[string]bt)

func handler(w http.ResponseWriter, r *http.Request) {
	trans := r.URL.Query().Get("t")
	j, err := json.Marshal(m[trans])
	if err != nil {
		fmt.Println("json.Marshal: ", err)
	}
	w.Write(j)
}


func main() {
	reader, err := ioutil.ReadFile("tr")
	if err != nil {
		fmt.Println("ioutil.Readfile: ", err)
	}

	fieldsB := bytes.Fields(reader)
	for i:=0; i<len(fieldsB); i+=3 {
		m[string(fieldsB[i])] = bt{string(fieldsB[i+1]), string(fieldsB[i+2])}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
