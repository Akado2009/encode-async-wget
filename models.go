package main

import (
	"github.com/BurntSushi/toml"
)

var (
	configPath = "config.toml"
	//AppConfig ...
	AppConfig Config
)

//Config ...
type Config struct {
	File                 string `toml:"File"`
	EncodeRoot           string `toml:"EncodeRoot"`
	FileURL              string `toml:"FileURL"`
	ExperimentURL        string `toml:"ExperimentURL"`
	ExperimentControlURL string `toml:"ExperimentControlURL"`
	MainURL              string `toml:"MainURL"`
}

//LoadConfig ..
func LoadConfig(c *Config) error {
	_, err := toml.DecodeFile(configPath, c)
	return err
}

//Experiments ...
type Experiments struct {
	Graph []Experiment `json:"@graph"`
}

//Experiment ...
type Experiment struct {
	Accession string `json:"accession"`
}

//Files ...
type Files struct {
	Graph []SubGraph `json:"@graph"`
}

//SubGraph ...
type SubGraph struct {
	Data []File `json:"files"`
}

//File ...
type File struct {
	ID string `json:"@id"`
}

// ("Accession\tDataset\tTissue\tLab\tLink\tDataType\n")

//FileResponse ...
type FileResponse struct {
	Accession  string    `json:"accession"`
	Dataset    string    `json:"dataset"`
	Lab        LabStruct `json:"lab"`
	OutputType string    `json:"output_type"`
	Href       string    `json:"href"`
}

//LabStruct ...
type LabStruct struct {
	Title string `json:"title"`
}

//ControlResponse ...
type ControlResponse struct {
	Controls []PossibleControls `json:"possible_controls"`
}

//PossibleControls ...
type PossibleControls struct {
	Files []string `json:"files"`
}
