package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getExperiments(URL string) *Experiments {
	experiments := &Experiments{}
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, experiments)
	return experiments
}

func getFiles(URL string) *Files {
	files := &Files{}
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, files)
	return files
}

func getFileResponse(URL string) *FileResponse {
	fResponse := &FileResponse{}
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, fResponse)
	return fResponse
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {

	if err := LoadConfig(&AppConfig); err != nil {
		log.Fatalf("Config loading failed: %v", err)
	}

	yes, err := exists(AppConfig.File)
	if err != nil {
		log.Fatal(err)
	}
	if !yes {
		ids := make([]string, 0)

		file, err := os.OpenFile(AppConfig.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		_, _ = datawriter.WriteString("Accession\tDataset\tTissue\tLab\tLink\tDataType\n")

		result := getExperiments(AppConfig.MainURL)
		for _, experiment := range result.Graph {
			files := getFiles(fmt.Sprintf(AppConfig.ExperimentURL, experiment.Accession))
			for _, file := range files.Graph {
				for _, subFile := range file.Data {
					ids = append(ids, subFile.ID)
					fResp := getFileResponse(fmt.Sprintf(AppConfig.FileURL, subFile.ID))
					if fResp.OutputType == "signal" || fResp.OutputType == "raw signal" {

						// check for controls
						// create a better table with controls [control1, control2] to just sum them?
						output := fmt.Sprintf(
							"%s\t%s\t%s\t%s\t%s\t%s\n",
							fResp.Accession,
							fResp.Dataset,
							"empty",
							fResp.Lab.Title,
							fmt.Sprintf(AppConfig.EncodeRoot, fResp.Href),
							fResp.OutputType)
						_, _ = datawriter.WriteString(output)
					}
				}
			}
		}

		for _, data := range ids {
			_, _ = datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
	}

	// for _, experiment := range encodeExps[:1] {
	// 	expResult := getRequestData(fmt.Sprintf(experimentURL, experiment), "@graph.#.files.#.@id")

	// 	files := make([]string, 0)
	// 	expResult.ForEach(func(key, value gjson.Result) bool {
	// 		fmt.Println(value.String())
	// 		files = append(files, value.String())
	// 		return true
	// 	})

	// 	fmt.Println(files)
	// 	for _, file := range files {
	// 		fileResult := getRequestData(fmt.Sprintf(fileURL, file), "")
	// 		fmt.Printf("%+v", fileResult)
	// 	}
	// }
}
