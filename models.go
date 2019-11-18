package main

import (
	"github.com/BurntSushi/toml"
)

var (
	configPath = "config.toml"
	//AppConfig ...
	AppConfig Config
)

//ErrorReport
type ErrorReport struct {
	Err  error
	Data string
}

//Sample
type Sample struct {
	Accession string
	Dataset   string
	Tissue    string
	CellLine  string
	Link      string
}

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
	Ontology OntologyStruct `json:"biosample_ontology"`
	Data     []File         `json:"files"`
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

//Ontology ...
type OntologyStruct struct {
	Classification string `json:"classification"` //can be either tissue/cell line (everything else = garbage (est est primary cell))
	TermName       string `json:"term_name"`
}

//ControlResponse ...
type ControlResponse struct {
	Controls []PossibleControls `json:"possible_controls"`
}

//PossibleControls ...
type PossibleControls struct {
	Files []string `json:"files"`
}

//FileDescription ...
type FileDescription struct {
	FileType string `json:"file_type"`
	Href     string `json:"href"`
}
