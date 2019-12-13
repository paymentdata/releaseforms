package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"

	_ "github.com/joho/godotenv/autoload"
)

type PullRequestID int
type PullRequestIDEmitter <-chan PullRequestID

var (
	client *github.Client
	ctx    = context.Background()

	re = regexp.MustCompile(`#[0-9]*`)

	//below searches are relatively arbitrary, current vals come from dependence on our org issue templates.
	uptoproposal       = regexp.MustCompile(`(?s)\*\*Purpose\*\*.*\*\*Proposal`)
	uptobugdescription = regexp.MustCompile(`(?s)\*\*Describe the bug\*\*.*\*\*To`)

	productrepo = os.Getenv("REPO")
	org         = os.Getenv("ORG")
)

//initialize github client that func main() consumes
func init() {
	if pat := os.Getenv("PAT"); len(pat) > 0 {
		client = github.NewClient(
			oauth2.NewClient(ctx, oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: pat},
			)),
		)
	} else {
		client = github.NewClient(nil)
	}
}

//aggregate context
func main() {

	var (
		PullRequestIDs PullRequestIDEmitter
		changes        ChangeItemEmitter

		rtd = ReleaseTemplateData{
			Date:        time.Now().String(),
			Product:     productrepo,
			BackOutProc: "git revert",
			PCIImpact:   "none",
			OWASPImpact: "none",
		}
	)

	PullRequestIDs = ingestPRs(os.Stdin)

	changes = PullRequestIDs.gatherChangeContexts(ctx, client)

	rtd.AggregateChanges(changes)

	rtd.Save(false /*toPDF flag, saves rendered HTML if false.*/)

}

func (prID PullRequestID) ConstructChangeItem(ctx context.Context, c *github.Client) ChangeItem {
	var (
		change ChangeItem
		pr     *github.PullRequest

		err error
	)

	change.ID = int(prID)
	pr, _, err = c.PullRequests.Get(ctx, org, productrepo, int(prID))
	if err != nil {
		panic(err)
	}
	change.Title = *pr.Title
	change.Developer = GetName(*pr.User.Login, ctx, c)
	change.SummaryOfChangesImplemented = *pr.Body
	change.CommitSHA = *pr.MergeCommitSHA
	if issueContext := re.Find([]byte(*pr.Body)); len(issueContext) > 0 {
		issueID, err := strconv.Atoi(string(issueContext[1:]))
		if err != nil {
			panic(err)
		}
		change.IssueID = issueID
		iss, _, err := c.Issues.Get(ctx, org, productrepo, issueID)
		if err != nil {
			panic(err)
		}
		summaryofissue := uptoproposal.Find([]byte(*iss.Body))
		summaryofissue = []byte(strings.ReplaceAll(string(summaryofissue), "**Purpose**\r\n", ""))
		summaryofissue = []byte(strings.ReplaceAll(string(summaryofissue), "\r\n\r\n**Proposal", ""))
		if len(summaryofissue) == 0 {
			summaryofissue = uptobugdescription.Find([]byte(*iss.Body))
			summaryofissue = []byte(strings.ReplaceAll(string(summaryofissue), "**Describe the bug**\r\n", ""))
			summaryofissue = []byte(strings.ReplaceAll(string(summaryofissue), "\r\n\r\n**To", ""))
		}
		if len(summaryofissue) > 0 {
			change.SummaryOfChangesNeeded = string(summaryofissue)
		}
	}
	reviews, _, err := c.PullRequests.ListReviews(ctx, org, productrepo, int(prID), nil)
	if err != nil {
		panic(err)
	}
	for _, r := range reviews {
		if *r.State == "APPROVED" {
			change.ApprovedBy += "[" + GetName(*r.User.Login, ctx, c) + "]"
		}
	}
	return change
}
func GetName(username string, ctx context.Context, c *github.Client) string {
	if name, ok := PeopleMap[username]; ok {
		return name
	}
	u, _, err := c.Users.Get(ctx, username)
	if err != nil {
		panic(err)
	}
	return *u.Name
}

//IngestPRs spawns a ~short-lived goroutine that consumes a stream of gob-encoded ints, which are sent over the returned PullRequestIDEmitter
func ingestPRs(input io.Reader) PullRequestIDEmitter {
	var (
		gd    = gob.NewDecoder(input)
		prIDs = make(chan PullRequestID)
	)
	go func(downstreamPRlistener chan<- PullRequestID) {
		var tmpnum int
		for {
			if err := gd.Decode(&tmpnum); err != nil {
				if err == io.EOF {
					close(prIDs)
					break
				} else {
					panic(err)
				}
			}
			downstreamPRlistener <- PullRequestID(tmpnum)
		}
	}(prIDs)
	return prIDs
}

//github context retriever gopher
func (prEmitter PullRequestIDEmitter) gatherChangeContexts(ctx context.Context, c *github.Client) ChangeItemEmitter {
	var changeItems = make(chan ChangeItem)

	go func(e PullRequestIDEmitter) {
		var (
			id   PullRequestID
			more bool
		)
		for {
			if id, more = <-e; more {
				changeItems <- id.ConstructChangeItem(ctx, c)
			} else {
				close(changeItems)
				break
			}
		}
	}(prEmitter)

	return changeItems
}

const (
	//PDFConverterEndpoint points to the url which ought to route to the /convert endpoint for an instance of the AthenaPDF microservice.
	PDFConverterEndpoint = "http://athenapdf:8080/convert?auth=%s&ext=html"
	weaverAuthKey        = "arachnys-weaver"
)

//GetPDF is intended to receive a rendered HTML template payload, which is intended to be submitted to the athenapdf service.
//The func responds with the PDF response as []byte, and/or an error if there is one.
func GetPDF(data []byte) ([]byte, error) {
	pdfEndpoint := fmt.Sprintf(PDFConverterEndpoint, weaverAuthKey)
	client := &http.Client{}

	// Create buffer
	buf := new(bytes.Buffer)

	r := bytes.NewReader(data)
	w := multipart.NewWriter(buf)

	// Create file field
	fw, err := w.CreateFormFile("file", "statementsummary")
	if err != nil {
		return nil, err
	}

	// Write file field from file to upload
	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST", pdfEndpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	obuf := new(bytes.Buffer)
	obuf.ReadFrom(res.Body)

	return obuf.Bytes(), err
}

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

func (rtd *ReleaseTemplateData) Save(toPDF bool) {
	var ext string
	if toPDF {
		ext = ".pdf"
	} else {
		ext = ".html"
	}
	f, err := os.Create(rtd.Product + "-" + rtd.Changes[0].CommitSHA + ext)
	if err != nil {
		panic(err)
	}
	if toPDF {
		pdfResponse, err := GetPDF(rtd.Render())
		if err != nil {
			panic(err)
		}
		_, err = f.Write(pdfResponse)
		if err != nil {
			panic(err)
		}
	} else {
		_, err = f.Write(rtd.Render())
		if err != nil {
			panic(err)
		}
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
