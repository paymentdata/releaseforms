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

type prID int
type prIDEmitter <-chan prID

func main() {

	var (
		ctx = context.Background()
		client *github.Client
		
		prIDs prIDEmitter
		changes <-chan form.ChangeItem
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

	var rtd form.ReleaseTemplateData
	rtd.Date = time.Now().String()
	rtd.Product = "somerepo"
	rtd.BackOutProc = "git revert"
	rtd.PCIImpact = "none"
	rtd.OWASPImpact = "none"

	prIDs = ingestPRs(os.Stdin)
	
	changes = prIDs.gatherChangeContexts(ctx, client)

	rtd.AggregateChanges(changes)
	
	rtd.Save()
}

func (pullRequestID prID) ConstructChangeItem(ctx context.Context, c *github.Client) form.ChangeItem {
	var (
		change form.ChangeItem
		pr     *github.PullRequest

		err error
	)

	change.ID = int(pullRequestID)
	pr, _, err = c.PullRequests.Get(ctx, org, productrepo, int(pullRequestID))
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
	reviews, _, err := c.PullRequests.ListReviews(ctx, org, productrepo, int(pullRequestID), nil)
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

//prID ingestion gopher
func ingestPRs(input io.Reader) prIDEmitter {
	var (
		err error
		gd  = gob.NewDecoder(input)

		tmpnum int
		prIDs  = make(chan prID, 0)
	)

	log.Println("firing off prID gopher")
	go func(downstreamPRlistener chan<- prID) {
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
			downstreamPRlistener <- prID(tmpnum)
		}
	}(prIDs)
	return prIDs
}

//github context retriever gopher
func (emitter prIDEmitter) gatherChangeContexts(ctx context.Context, c *github.Client) <-chan form.ChangeItem {
	var (
		changeItems = make(chan form.ChangeItem, 0)
	)
	log.Println("firing off github gopher")
	go func(prIDs <-chan prID) {
		for {
			var (
				id   prID
				more bool
			)
			if id, more = <-prIDs; more {
				log.Printf("github gopher constructing change item for prID[%d]", id)
				changeItems <- id.ConstructChangeItem(ctx, c)
			} else {
				log.Println("closing changeItems chan")
				close(changeItems)
				break
			}
		}
	}(emitter)
	return changeItems
}
