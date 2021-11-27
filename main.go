package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	user := "phobyjun"
	token := os.Getenv("MY_GITHUB_TOKEN")

	ctx := context.Background()
	client := getClientByToken(ctx, token)

	fileContent := []byte("Content of File")
	filePath := "_posts/test.md"
	err := autoCreatePost(ctx, client, user, fileContent, filePath)
	if err != nil {
		fmt.Println(err)
	}
}

func getClientByToken(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client
}

func autoCreatePost(ctx context.Context, client *github.Client, user string, fileContent []byte, filePath string) error {
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("auto create by go-github"),
		Content:   fileContent,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("Junseok Yoon"), Email: github.String("phobyjun@gmail.com")},
	}
	_, _, err := client.Repositories.CreateFile(ctx, user, user+".github.io", filePath, opts)
	if err != nil {
		return err
	}
	return nil
}
