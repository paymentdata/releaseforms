package listener

import (
	"fmt"
	"log"
	"os"

	"github.com/paymentdata/releaseforms/util"

	"github.com/paymentdata/releaseforms/form"

	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	path1 = "/webhooks1"
	path2 = "/webhooks2"

	pdfext = ".pdf"
)

func Listen() {
	hook1, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))
	hook2, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect2...?"))

	http.HandleFunc(path1, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook1.Parse(r, github.ReleaseEvent, github.PullRequestEvent, github.PushEvent)
		if err != nil {
			fmt.Printf("Parse error: %s", err.Error())
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}

		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			fmt.Printf("commit %s authored by %v pushed @ %s\n",
				push.HeadCommit.ID[:7],
				push.HeadCommit.Author,
				push.HeadCommit.Timestamp)

			var rtd form.ReleaseTemplateData
			rtd.Commit = push.HeadCommit.ID[:7]
			rtd.CommitterName = push.HeadCommit.Committer.Email
			rtd.Author = push.HeadCommit.Author.Name
			rtd.Date = push.HeadCommit.Timestamp
			rtd.OWASPImpact = "none"
			rtd.PCIImpact = "minimal"
			rtd.BackOutProc = "revert this change"
			rtd.Product = push.Repository.Name

			f, err := os.Create(rtd.Commit + pdfext)
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
		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	})

	http.HandleFunc(path2, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook2.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	})
	http.ListenAndServe(os.Getenv("GitHookPort"), nil)
}
