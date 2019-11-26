package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/paymentdata/releaseforms/form"
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

		err    error
		prIDs  []int
		tmpnum int
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
	gd := gob.NewDecoder(os.Stdin)

	for {
		err = gd.Decode(&tmpnum)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		prIDs = append(prIDs, tmpnum)
	}

	var rtd form.ReleaseTemplateData
	rtd.Date = time.Now().String()
	rtd.Product = "somerepo"
	rtd.BackOutProc = "git revert"
	rtd.PCIImpact = "none"
	rtd.OWASPImpact = "none"
	for _, num := range prIDs {
		rtd.Changes = append(rtd.Changes, ConstructChangeItem(ctx, num, client))
	}

	f, err := os.Create(rtd.Product + "-" + rtd.Changes[0].CommitSHA + ".pdf")
	if err != nil {
		panic(err)
	}
	pdfResponse, err := util.GetPDF(rtd.Render())
	if err != nil {
		panic(err)
	}
	n, err := f.Write(pdfResponse)
	if err != nil {
		panic(err)
	}
	log.Printf("Wrote %d bytes!", n)
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func ConstructChangeItem(ctx context.Context, pullRequestID int, c *github.Client) form.ChangeItem {
	var change form.ChangeItem
	change.ID = pullRequestID
	pr, _, err := c.PullRequests.Get(ctx, org, productrepo, pullRequestID)
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
		fmt.Printf("addresses issue[%d]:(%s) opened by %s\n\t", *iss.Number, *iss.Title, GetName(*iss.User.Login, ctx, c))

		if len(summaryofissue) > 0 {
			change.SummaryOfChangesNeeded = string(summaryofissue)
			fmt.Printf("description of issue:\n\t%s", summaryofissue)
		}
	}
	reviews, _, err := c.PullRequests.ListReviews(ctx, org, productrepo, pullRequestID, nil)
	if err != nil {
		panic(err)
	}
	for _, r := range reviews {
		if *r.State == "APPROVED" {
			change.ApprovedBy += "[" + GetName(*r.User.Login, ctx, c) + "]"
			fmt.Printf("\tapproved by %s\n", GetName(*r.User.Login, ctx, c))
		}
	}
	fmt.Printf("changeitem: %+v", change)
	fmt.Printf("\n\n##############################################\n\n")
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
