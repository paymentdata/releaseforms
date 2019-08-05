package listener

import (
	"fmt"

	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	path1 = "/webhooks1"
	path2 = "/webhooks2"
)

func listen() {
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
			fmt.Printf("commit %s authored by %v pushed @ %s",
				push.Commits[0].ID[:7],
				push.Commits[0].Author,
				push.Commits[0].Timestamp)

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
	http.ListenAndServe(":8080", nil)
}
