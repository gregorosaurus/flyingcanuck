package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

var htmlPostTemplate = `
<a class="post" href="posts/{{.Path}}/">
	<div class="post-title">{{.LongTitle}}</div>
	<div class="post-date">{{.PublishDate}}</div>
</a>`

func main() {
	log.Println("Starting blog generation process")

	posts := findPosts()

	var postsBuffer = new(bytes.Buffer)
	postsTemplate, err := template.New("posts").Parse(htmlPostTemplate)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating posts HTML")

	for _, post := range posts {
		postPath := filepath.Join("..", "posts", post.Path)
		postHtmlPath := filepath.Join(postPath, "index.html")
		postMarkDownPath := filepath.Join(postPath, "post.md")
		postMarkdown, err := ioutil.ReadFile(postMarkDownPath)
		if err != nil {
			log.Fatalf("Unable to read markdown file: %s", err)
		}
		htmlContent := markdown.ToHTML(postMarkdown, nil, nil)
		err = ioutil.WriteFile(postHtmlPath, htmlContent, os.ModePerm)
		if err != nil {
			log.Fatalf("Unable to write post file. ")
		}

		err = postsTemplate.Execute(postsBuffer, post)
		if err != nil {
			log.Fatalf("Unable to execute the template: %s", err)
		}
	}

	mainHtmlData, err := ioutil.ReadFile(filepath.Join("..", "index.template.html"))
	if err != nil {
		log.Fatalf("Unable to find main template html. %s", err)
	}

	log.Println("Generating main index html")

	mainHtml := string(mainHtmlData)
	mainHtml = strings.ReplaceAll(mainHtml, "{{posts}}", postsBuffer.String())
	mainHtmlPath := filepath.Join("..", "index.html")

	err = ioutil.WriteFile(mainHtmlPath, []byte(mainHtml), os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to write index file %s", err)
	}

	log.Println("Process completed.")
}

func findPosts() []*Post {
	log.Println("Finding posts to render")

	var posts = make([]*Post, 0)

	items, err := ioutil.ReadDir(filepath.Join("..", "posts"))
	if err != nil {
		log.Fatalf("Unable to find posts under the posts folder.  Are you in the right directory? %s", err)
	}
	for _, item := range items {
		if item.IsDir() {
			log.Printf("Found post directory: %s", item.Name())
			var post = new(Post)
			post.Path = item.Name()
			//attempt to get the date from the post directory
			postFolderComponent := strings.Split(item.Name(), "_")
			if len(postFolderComponent) != 2 {
				log.Fatalf("Unable to properly analyze the post folder %s - is a non post folder in the post directory?", item.Name())
			}

			post.PublishDate, err = time.Parse("2006-01-02", postFolderComponent[0])
			if err != nil {
				log.Fatalf("Error converting folder to publish date for '%s': %s", item.Name(), err)
			}

			post.ShortTitle = postFolderComponent[1]

			//read first line to get the long title.
			postMarkDownFile, err := os.Open(filepath.Join("..", "posts", item.Name(), "post.md"))
			if err != nil {
				log.Fatalf("Unable to open post.md for %s, %s", item.Name(), err)
			}

			fileScanner := bufio.NewScanner(postMarkDownFile)
			fileScanner.Split(bufio.ScanLines)
			fileScanner.Scan()
			firstLine := fileScanner.Text()
			post.LongTitle = strings.TrimSpace(strings.TrimLeft(firstLine, "#"))

			postMarkDownFile.Close()

			posts = append(posts, post)
		}
	}
	return posts
}
