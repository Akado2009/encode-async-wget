// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// func checkForControl(URL string) []string {
// 	fmt.Println(URL)
// 	controls := make([]string, 0)
// 	cResponse := &ControlResponse{}
// 	response, err := http.Get(URL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	data, _ := ioutil.ReadAll(response.Body)
// 	err = json.Unmarshal(data, cResponse)
// 	fmt.Printf("%+v", cResponse)
// 	return controls
// }

// func main() {
// 	if err := LoadConfig(&AppConfig); err != nil {
// 		log.Fatalf("Config loading failed: %v", err)
// 	}
// 	controls := checkForControl(fmt.Sprintf(AppConfig.ExperimentURL, "ENCSR370RJP"))

// 	fmt.Println(controls)
// }
