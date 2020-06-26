package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type PullRequest struct {
	Owner string
	Repo string
	Number int
} 

func (pr *PullRequest) Comment(output string) {
	c := fmt.Sprintf("```diff\n%s```", output)
	service := github.PullRequestsService{}
	comment := github.PullRequestComment{
		Body: &c,
	}

	_, resp, err := service.CreateComment(context.Background(), pr.Owner, pr.Repo, pr.Number, &comment)

	fmt.Println(comment, resp, err)
}
