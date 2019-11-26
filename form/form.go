package form

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

var funcMap = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"inc": func(i int) int {
		return i + 1
	},
}

//ReleaseTemplateData is the encapsulating struct for software release forms.
type ReleaseTemplateData struct {
	Date        string       `json:"Date"`
	Product     string       `json:"Product"`
	Changes     []ChangeItem `json:"ChangeItems"`
	BackOutProc string       `json:"BackOutProc"`
	PCIImpact   string       `json:"PCIImpact"`
	OWASPImpact string       `json:"OWASPImpact"`
}

//ChangeItem represents a unit of change, where and engineering solution addresses a biz need, either a feature or a bug.
type ChangeItem struct {
	IssueID, //Issue ID
	ID int // PR ID

	Title, // PR Title
	CommitSHA, // Merge Commit
	RequestedBy, // Source of Issue
	SummaryOfChangesNeeded, // Issue Body as Need and PR Body as Solution
	SummaryOfChangesImplemented, // Issue Body as Need and PR Body as Solution
	Notes, // Issue Body
	Developer, // PR User
	TestedBy, // PR User
	CodeReviewAndTesting, // Approving Reviewers + Developer
	CodeReviewAndTestingNotes, // Reviewer comments?
	ApprovedBy string // Approving Reviewers
}

//Render is a receiver which returns the ReleaseTemplateData as a []byte payload.
//Used for transporting the bytes over the network currently.
func (rtd *ReleaseTemplateData) Render() []byte {
	var (
		tpl bytes.Buffer
		e   error
	)
	t, e := template.New("thing").Funcs(funcMap).Parse(ReleaseTemplate)
	if e != nil {
		fmt.Printf("Err %v", e)
	}
	e = t.Execute(&tpl, rtd)
	if e != nil {
		log.Printf("TEMPLATE ERR: %v\n", e)
	}
	return tpl.Bytes()
}
