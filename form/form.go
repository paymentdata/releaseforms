package form

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

//ReleaseTemplateData is the encapulative struct for software release forms.
type ReleaseTemplateData struct {

	//Release Date
	Date string `json:"Date"`
	//Product
	Product string `json:"Product"`

	//Included changes
	Commit string `json:"Commit"`

	CommitterName string `json:"Committer"`
	Author        string `json:"Author"`

	BackOutProc string `json:"BackOutProc"`

	PCIImpact   string `json:"PCIImpact"`
	OWASPImpact string `json:"OWASPImpact"`
}

//Commit is the Change Item primitive
type Commit struct {
	Text        string
	RequestedBy string
	SummaryOfChanges,
	Notes,
	Developer,
	TestedBy,
	CodeReviewAndTesting,
	CodeReviewAndTestingNotes,
	ApprovedBy string
}

func (rtd *ReleaseTemplateData) Render() []byte {
	var (
		tpl bytes.Buffer
		e   error
	)
	t, e := template.New("thing").Parse(ReleaseTemplate)
	if e != nil {
		fmt.Printf("Err %v", e)
	}
	e = t.Execute(&tpl, rtd)
	if e != nil {
		log.Printf("TEMPLATE ERR: %v\n", e)
	}
	return tpl.Bytes()
}
