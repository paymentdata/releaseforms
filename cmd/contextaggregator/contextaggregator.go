package main

import (
	"context"
	"encoding/gob"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/paymentdata/releaseforms/form"
	"github.com/paymentdata/releaseforms/people"
	"github.com/paymentdata/releaseforms/util"
	"golang.org/x/oauth2"

	_ "github.com/joho/godotenv/autoload"
)

var re = regexp.MustCompile(`#[0-9]*`)

var uptoproposal = regexp.MustCompile(`(?s)\*\*Purpose\*\*.*\*\*Proposal`)
var uptobugdescription = regexp.MustCompile(`(?s)\*\*Describe the bug\*\*.*\*\*To`)

var (
	productrepo = os.Getenv("REPO")
	org         = os.Getenv("ORG")
)

func main() {

	var (
		ctx = context.Background()

		client *github.Client

		err error

		tmpnum int
		prIDs  = make(chan int, 0)
		done   = make(chan struct{})
	)

	if pat := os.Getenv("PAT"); len(pat) > 0 {
		client = github.NewClient(
			oauth2.NewClient(ctx, oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: os.Getenv("PAT")},
			)),
		)
	} else {
		client = github.NewClient(nil)
	}

	//prID ingestion gopher
	gd := gob.NewDecoder(os.Stdin)
	go func() {
		log.Println("firing off prID gopher")
		for {
			if err = gd.Decode(&tmpnum); err != nil {
				if err == io.EOF {
					close(prIDs)
					break
				} else {
					panic(err)
				}
			}
			log.Printf("received prID[%d]", tmpnum)
			prIDs <- tmpnum
		}
	}()

	var rtd form.ReleaseTemplateData
	rtd.Date = time.Now().String()
	rtd.Product = "somerepo"
	rtd.BackOutProc = "git revert"
	rtd.PCIImpact = "none"
	rtd.OWASPImpact = "none"

	//github context retriever gopher
	go func(rf *form.ReleaseTemplateData) {
		log.Println("firing off github gopher")
		for {
			var (
				prID int
				more bool
			)
			prID, more = <-prIDs
			if more {
				log.Printf("github gopher processing change item for prID[%d]", prID)
				rf.Changes = append(rf.Changes, ConstructChangeItem(ctx, prID, client))
			} else {
				log.Println("closing done chan")
				close(done)
				break
			}
		}
	}(&rtd)

	//wait for pipeline gophers to complete their jobs
	<-done

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

func ConstructChangeItem(ctx context.Context, pullRequestID int, c *github.Client) form.ChangeItem {
	var (
		change form.ChangeItem
		pr     *github.PullRequest

		err error
	)

	change.ID = pullRequestID
	pr, _, err = c.PullRequests.Get(ctx, org, productrepo, pullRequestID)
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
	reviews, _, err := c.PullRequests.ListReviews(ctx, org, productrepo, pullRequestID, nil)
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
	if name, ok := people.PeopleMap[username]; ok {
		return name
	}
	u, _, err := c.Users.Get(ctx, username)
	if err != nil {
		panic(err)
	}
	return *u.Name
}
