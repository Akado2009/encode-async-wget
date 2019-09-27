package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

func checkForControl(URL string) []string {
	cResponse := &ControlResponse{}
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, cResponse)
	if len(cResponse.Controls) > 0 {
		return cResponse.Controls[0].Files
	}
	return make([]string, 0)
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

		file, err := os.OpenFile(AppConfig.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		_, _ = datawriter.WriteString("Accession\tDataset\tTissue\tCellLine\tPrimaryCell\tLab\tLink\tDataType\tControls\n")

		result := getExperiments(AppConfig.MainURL)

		previousDataset := ""
		for _, experiment := range result.Graph {
			files := getFiles(fmt.Sprintf(AppConfig.ExperimentURL, experiment.Accession))
			for _, file := range files.Graph {
				tissue := "."
				cellLine := "."
				primaryCell := "."
				switch file.Ontology.Classification {
				case "cell line":
					cellLine = file.Ontology.TermName
				case "tissue":
					tissue = file.Ontology.TermName
				case "primary cell":
					primaryCell = file.Ontology.TermName
				}
				for _, subFile := range file.Data {
					fResp := getFileResponse(fmt.Sprintf(AppConfig.FileURL, subFile.ID))
					if fResp.OutputType == "signal" || fResp.OutputType == "raw signal" {

						// check for controls
						// create a better table with controls [control1, control2] to just sum them?
						if previousDataset != fResp.Dataset {
							previousDataset = fResp.Dataset

							controls := checkForControl(fmt.Sprintf(AppConfig.ExperimentControlURL, fResp.Dataset))
							//AddControls! (can add bam, since we have bam2wig)
							output := fmt.Sprintf(
								"%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
								fResp.Accession,
								fResp.Dataset,
								tissue,
								cellLine,
								primaryCell,
								fResp.Lab.Title,
								fmt.Sprintf(AppConfig.EncodeRoot, fResp.Href),
								fResp.OutputType,
								strings.Join(controls, ", "))
							_, _ = datawriter.WriteString(output)
						}
					}
				}
			}
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
