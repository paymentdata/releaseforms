package listener

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/paymentdata/releaseforms/form"
	"github.com/paymentdata/releaseforms/util"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	path1 = "/webhooks1"

	pdfext = ".pdf"
)

//Listen will construct the *github.Webhook, register the http.HandlerFunc, and ListenAndServe the handler over os.Getenv("GitHookPort")
func Listen() {
	hook1, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))

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
			if push.Deleted {
				fmt.Printf("%s (%s) deleted %s\n", push.Pusher.Name, push.Pusher.Email, push.Ref)
				return
			}
			fmt.Printf("commit %s authored by %v pushed @ %s\n",
				push.HeadCommit.ID[:7],
				push.HeadCommit.Author,
				push.HeadCommit.Timestamp)

			var rtd form.ReleaseTemplateData
			//TODO for CD, this needs updating for ~new changelist
			rtd.Date = push.HeadCommit.Timestamp
			rtd.OWASPImpact = "none"
			rtd.PCIImpact = "minimal"
			rtd.BackOutProc = "revert this change"
			rtd.Product = push.Repository.Name

			f, err := os.Create(rtd.Changes[0].CommitSHA + pdfext)
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

	http.ListenAndServe(os.Getenv("GitHookPort"), nil)
}
