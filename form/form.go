package form

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"os"

	"github.com/paymentdata/releaseforms/util"
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

type ChangeItemEmitter <-chan ChangeItem

func (rtd *ReleaseTemplateData) AggregateChanges(rx ChangeItemEmitter) {
	for {
		var (
			change ChangeItem
			more   bool
		)
		if change, more = <-rx; more {
			log.Printf("adding constructed change for prID[%d]", change.ID)
			rtd.Changes = append(rtd.Changes, change)
		} else {
			log.Println("aggregation of changes complete")
			break
		}
	}
}

func (rtd *ReleaseTemplateData) Save() {

	f, err := os.Create(rtd.Product + "-" + rtd.Changes[0].CommitSHA + ".pdf")
	if err != nil {
		panic(err)
	}
	pdfResponse, err := util.GetPDF(rtd.Render())
	if err != nil {
		panic(err)
	}
	_, err = f.Write(pdfResponse)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {

		panic(err)
	}
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
