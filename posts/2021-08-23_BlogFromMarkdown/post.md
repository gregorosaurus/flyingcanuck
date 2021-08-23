# Building a Static Blog from Markdown
*2021-08-23*

<img src="meta.gif" width="300"/>

## Summary

When I wanted to start a blog, I wanted something to be 1) simple, 2) easy to maintain, and 3) cheap to run (free). 
I decided that one of the best ways to accomplish this was to make a static web site (ie: no back end) that was build by using Markdown and running on Azure using either Static Web Apps, or just a storage account.  (I ended up putting this on a Azure Web App, but more on that later.)

So what is Markdown? Markdown is (From Wikipedia) [a lightweight markup language for creating formatted text using a plain-text editor. Markdown is widely used in blogging, instant messaging, online forums, collaborative software, documentation pages, and readme files.](https://en.wikipedia.org/wiki/Markdown) Long story short, Markdown is a great option for creating content as it's easy to read, easy to write, and widely used in the tech community. 

For my post repository and storage, I am using github.  Additionally, I am using github actions to build and deploy the blog. 

## Overview

The github repo for this blog is located and available here: [https://github.com/gregorosaurus/flyingcanuck/](https://github.com/gregorosaurus/flyingcanuck/)

<img src="2021-08-23 10_16_26-gregorosaurus_flyingcanuck and 4 more pages - Work - Microsoft​ Edge.png" width="300">

As shown above, the blog is has the following components:

1. The template.html file. This is the main HTML template for the entire blog. The template html is rather simple, with the body of each page being injected into:
```html
    <div class="content">
        {{Body}}
    </div>
```
2. The main.css file.  This is the stylesheet for the blog.  Pretty self explanatory.
3. The posts directory. This is the directory in which all posts are saved.  Each post has its own directory within the posts directory. Inside the directory is the main post.md (all posts require this), and any supporting files like images.  The directory tree looks something like this:
- posts
    - 2021-08-23_Post1ShortTitle
        - post.md
        - image.jpg
    - 2021-08-30_Post2Title
        - post.md
        - image.gif
4. The html generation script.  This is a powershell script using the new ```ConvertFrom-Markdown``` commandlet introduced in powershell 6. This scripts responsibility is to iterate through each of the directories in the post directory, and build the index.html of the main page **and** each of the post index.html files.  These are the actual files to be served to the browser.  This script is called by my CI/CD process.   The index.html files [are ignored](https://git-scm.com/docs/gitignore) in the git repository.

### Generation Powershell Script
The generation script is probably the most important piece of this. Its job is to generate the needed index.html files to be served to the browser. 

## CI/CD

