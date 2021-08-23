# Building a Static Blog from Markdown
*2021-08-23*

<img src="meta.gif" width="300"/>

## Summary

When I wanted to start a blog, I wanted something to be 1) simple, 2) easy to maintain, and 3) cheap to run (free). 
I decided that the best way to accomplish this was to make a *static web site* (ie: no back end) that had content sourced by Markdown files. Using a static web site, it could be hosted on Azure using either Static Web Apps, or just a storage account.  (I ended up putting this on a Azure Web App, but more on that later.)

So what is Markdown? Markdown is (From Wikipedia) [a lightweight markup language for creating formatted text using a plain-text editor. Markdown is widely used in blogging, instant messaging, online forums, collaborative software, documentation pages, and readme files.](https://en.wikipedia.org/wiki/Markdown) Long story short, Markdown is a great option for creating content as it's easy to read, easy to write, and widely used in the tech community.  Additionally, it requires exactly zero ongoing maintenance.  They're just files! 

But these files need to be saved somewhere!  I decided that using github would be a perfect place for the blog source files (obviously).  Additionally, I am using github actions to build and deploy the blog, because I like free. 

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
    > Note: the format of the post directories is {Date}_{ShortTitle}.  This format is important as its used in the generation script below. 
4. The html generation script.  This is a powershell script using the new ```ConvertFrom-Markdown``` commandlet introduced in powershell 6. This scripts responsibility is to iterate through each of the directories in the post directory, and build the index.html of the main page **and** each of the post index.html files.  These are the actual files to be served to the browser.  This script is called by my CI/CD process.   The index.html files [are ignored](https://git-scm.com/docs/gitignore) in the git repository, as they are generated as part of the CI/CD process.

### Generation Powershell Script
The generation script is probably the most important piece of this. Its job is to generate the needed index.html files to be served to the browser. 

**Pseudo code**

The Pseudo code for the script is as follows:
```
Foreach directory under "posts"
    Split the directory by an underscore
    Parse the first part of the directory as a date
    Gather all other information about the post (title, short title, etc)
    Add to the list of posts
    Convert the post markdown file to html, save it as index.html

Foreach post found
    Create an index.html using the template.html file
    inject the post list in the index.html
```

## Hosting Location
This blog ended up being hosted on an **Azure Web App**.  I had a shared App Service that hosts a bunch of random websites, so I ended up using that. So effectively this is free, as I already was paying for a App Service Plan. 

I did try to use Azure Static Web Sites, but the DNS provider I used for my domains didn't support ANAME aliases that Azure Static Web Sites required for custom domains.  

I could use an Azure Storage Account to host static content, but I really wanted TLS which is currently not supported for custom domain names on storage account static content.  If I didn't care about TLS, I would have chosen this option.  

## CI/CD

The CI/CD process has two steps:
1. Build the blog from mark down
2. Deploy the blog to the Azure Web App

### The Workflow

<img src="2021-08-23 11_19_11-Updated readme · gregorosaurus_flyingcanuck@5ce5053 and 10 more pages - Work - M.png">

Because I'm using github as the repository for the blog, I'm using github actions for CI/CD.  The workflow that is used to build and deploy the blog is located [here](https://github.com/gregorosaurus/flyingcanuck/blob/main/.github/workflows/DeployWorkflow.yaml).  
But generally has these steps: 

1. Checkout the repo
2. Build the blog using the generate-blog.ps1 powershell script
3. Deploy the blog to Azure

> For information on how to setup a CI/CD pipeline with azure, reference the following article. https://docs.microsoft.com/en-us/azure/app-service/deploy-github-actions?tabs=applevel  
> **With a static web site, you don't need to include the build steps, just the deploy steps.**

## Summary
So in the end, I have a blog backed by a set of Markdown files, saved in github, that is automatically pushed to an Azure App Service using github actions. Because there is no required back end, pages load in milliseconds, which is definitely a nice added benefit.  

What I'll end up using this blog for? Probably just things I feel that should be shared, things I want to remember, or just pictures of giraffes, who knows.  Anyway, if you got this far, thanks for reading. 

Actually maybe the giraffe thing isn't such a bad idea:
<img src="giraffe.png" width="400"/>