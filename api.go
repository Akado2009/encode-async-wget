package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

func main() {
	url := "https://www.encodeproject.org/search/?type=Experiment&assay_term_name=Histone%20ChIP-seq&replicates.library.biosample.donor.organism.scientific_name=Homo+sapiens&status=released&limit=all&format=json"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	strData := string(data)
	result := gjson.Get(strData, "@graph.#.accession")

	encodeExps := make([]string, 0)

	result.ForEach(func(key, value gjson.Result) bool {
		encodeExps = append(encodeExps, value.String())
		return true // keep iterating
	})

	fmt.Println(len(encodeExps))
}
