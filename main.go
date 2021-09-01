package main

import (
	"context"
	"os"

	copy "goGitBack/copy"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v38/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load(".env")
	token := string(os.Getenv("githubPAT"))
	targetOrg := string(os.Getenv("targetOrg"))
	payload := string(os.Getenv("payload"))
	branchName := string(os.Getenv("branchName"))

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, _ := client.Repositories.ListByOrg(ctx, targetOrg, nil)

	count := len(repos)
	for i := 0; i < count; i++ {
		repo := *repos[i].Name
		url := *repos[i].CloneURL
		r, err := git.PlainCloneContext(ctx, repo, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "2",
				Password: token,
			},
			URL: url,
		})
		if err == nil {
			target := "./" + repo + "/."

			copy.Copy(payload, target)
			w, _ := r.Worktree()
			headRef, _ := r.Head()

			ref := plumbing.NewHashReference(plumbing.ReferenceName(branchTarget(branchName)), headRef.Hash())
			r.Storer.SetReference(ref)
			w.Add(payload)
			w.Commit("Added Payload", &git.CommitOptions{})
			r.Push(&git.PushOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "2",
					Password: token,
				},
			})
		}
	}
}

func branchTarget(branchName string)  string {
	return "refs/heads/" + branchName
}