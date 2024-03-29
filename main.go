package main

import (
	"bufio"
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// hello world
func main() {
	targetDir := "./TIL"
	fileHistory := ".history"
	posts := getNewPosts(targetDir, fileHistory)

	user := "phobyjun"
	token := os.Getenv("ACCESS_TOKEN")

	ctx := context.Background()
	client := getClientByToken(ctx, token)

	for _, post := range posts {
		s := strings.Split(post, "\\")
		filePath := "_posts/" + s[len(s)-1]
		fileContent := readMarkDownToByte(post)
		autoCreatePost(ctx, client, user, fileContent, filePath)
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

func autoCreatePost(ctx context.Context, client *github.Client, user string, fileContent []byte, filePath string) {
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("auto create by go-github"),
		Content:   fileContent,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("Junseok Yoon"), Email: github.String("phobyjun@gmail.com")},
	}
	_, _, err := client.Repositories.CreateFile(ctx, user, user+".github.io", filePath, opts)
	checkErr(err)
}

func readMarkDownToByte(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	checkErr(err)
	return content
}

func getNewPosts(targetDir, fileHistory string) []string {
	files := getCurrentFiles(targetDir)
	history := getFileHistory(fileHistory)
	posts := difference(files, history)
	syncHistory(files, fileHistory)
	return posts
}

func syncHistory(files []string, fileHistory string) {
	f, err := os.OpenFile(
		fileHistory,
		os.O_WRONLY,
		os.FileMode(0644))
	defer func() {
		err := f.Close()
		checkErr(err)
	}()
	checkErr(err)
	writer := bufio.NewWriter(f)
	for _, file := range files {
		_, err := writer.WriteString(file + "\n")
		checkErr(err)
	}
	wErr := writer.Flush()
	checkErr(wErr)
}

func getFileHistory(fileHistory string) []string {
	var history []string
	file, err := os.Open(fileHistory)
	defer func() {
		err := file.Close()
		checkErr(err)
	}()
	checkErr(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
	return history
}

func getCurrentFiles(targetDir string) []string {
	var files []string
	err := filepath.Walk(targetDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	checkErr(err)
	return files
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
