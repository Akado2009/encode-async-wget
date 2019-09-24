package main

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
