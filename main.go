package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

func main() {
	user := "phobyjun"
	token := os.Getenv("MY_GITHUB_TOKEN")

	ctx := context.Background()
	client := getClientByToken(ctx, token)

	fileContent := readMarkDownToByte("./TIL/test1.md")
	filePath := "_posts/test1.md"
	err := autoCreatePost(ctx, client, user, fileContent, filePath)
	checkErr(err)

	targetDir := "./TIL"
	files := getCurrentFiles(targetDir)
	for _, file := range files {
		fmt.Println(file)
	}
}

func getCurrentFiles(targetDir string) []string {
	files, err := ioutil.ReadDir(targetDir)
	checkErr(err)
	var f []string
	for _, file := range files {
		f = append(f, file.Name())
	}
	return f
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

func readMarkDownToByte(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	checkErr(err)
	return content
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
