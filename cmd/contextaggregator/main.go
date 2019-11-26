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

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"

	_ "github.com/joho/godotenv/autoload"
)

var re = regexp.MustCompile(`#[0-9]*`)

var uptoproposal = regexp.MustCompile(`(?s)\*\*Purpose\*\*.*\*\*Proposal`)
var uptobugdescription = regexp.MustCompile(`(?s)\*\*Describe the bug\*\*.*\*\*To`)

func main() {
	if len(os.Getenv("PAT")) == 0 {
		log.Fatal("this process is initially intended to service private repositories, and thus requires a populated $PAT env var")
	}

	var (
		ctx = context.Background()

		tc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("PAT")},
		))

		client = github.NewClient(tc)
	)

	var (
		err    error
		prIDs  []int
		tmpnum int
	)

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

	for _, num := range prIDs {
		ConstructChangeItem(ctx, num, client)
	}

}

func ConstructChangeItem(ctx context.Context, pullRequestID int, c *github.Client) {
	pr, _, err := c.PullRequests.Get(ctx, "someorg", "somerepo", pullRequestID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PR[%d]: %s (opened by %s)\n\t", pullRequestID, *pr.Title, *pr.User.Login)
	if len(*pr.Body) > 0 {
		if len(*pr.Body) > 80 {
			fmt.Printf("body: [%s]\n\t", (*pr.Body)[:80])
		} else {
			fmt.Printf("body: [%s]\n\t", *pr.Body)
		}
	}
	if issueContext := re.Find([]byte(*pr.Body)); len(issueContext) > 0 {
		issueID, err := strconv.Atoi(string(issueContext[1:]))
		if err != nil {
			panic(err)
		}
		iss, _, err := c.Issues.Get(ctx, "someorg", "somerepo", issueID)
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
			fmt.Printf("description of issue:\n\t%s", summaryofissue)
		}
	}
	reviews, _, err := c.PullRequests.ListReviews(ctx, "someorg", "somerepo", pullRequestID, nil)
	if err != nil {
		panic(err)
	}
	for _, r := range reviews {
		if *r.State == "APPROVED" {
			fmt.Printf("\tapproved by %s\n", GetName(*r.User.Login, ctx, c))
		}
	}
	fmt.Printf("\n\n##############################################\n\n")
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
