#
#
#  Generate-Blog
#  Author: Greg Mardon
#  
#  This script generates the html'd mark down for the blog
#
#
###################################################

class Post {
    [string]$ShortTitle
    [string]$LongTitle
    [string]$Path
    [datetime]$PublishDate
}

[string]$PostListHTMLTemplate = @"
<a class="post" href="posts/{{Path}}/">
	<div class="post-title">{{LongTitle}}</div>
	<div class="post-date">{{FormattedPublishDate}}</div>
</a>
"@

function Find-Posts {
    Write-Host "Finding posts in posts directory"

    $posts = @()

    $postDirectories = Get-ChildItem -Directory -Path "posts"
    if ($postDirectories.Length -eq 0) {
        Write-Error "No post directories found!"
        return $posts
    }

    foreach ($postDirectory in $postDirectories) {
        $post = [Post]::new()
        Write-Host "Found post directory" $postDirectory.Name
        $post.Path = $postDirectory.Name

        #attempt to get the date from the directory.
        $postPathComponent = $postDirectory.Name.Split("_")
        if ($postPathComponent.Length -ne 2) {
            Write-Error  -Message ("Error in path component length for post " + $postDirectory.Name)
            return $posts
        }

        $post.PublishDate = [System.DateTime]::Parse($postPathComponent[0])
        $post.ShortTitle = $postPathComponent[1]

        $postMarkdownContent = Get-Content -Path (Join-Path $postDirectory.FullName "post.md")
        $lines = $postMarkdownContent.Split([System.Environment]::NewLine)
        $post.LongTitle = $lines[0].TrimStart("#").Trim()

        $posts += $post
    }

    return $posts
}

Write-Host "Starting blog generation process"

$posts = Find-Posts
Write-Host "Found" $posts.Length "posts to process"

$templateHTML = Get-Content -Path "template.html"

$postListHtml = ""
foreach ($post in $posts) {
    $postItemHtml = $PostListHTMLTemplate.Replace("{{Path}}", $post.Path).Replace("{{ShortTitle}}", $post.ShortTitle).Replace("{{LongTitle}}", $post.LongTitle).Replace("{{FormattedPublishDate}}", $post.PublishDate.ToString("yyyy-MM-dd"))
    $postListHtml += $postItemHtml + [System.Environment]::NewLine

    #Now deal with the mark down.
    $postHtml = $templateHTML
    $postMarkDownHtml = ConvertFrom-Markdown -Path (Join-Path "posts" $post.Path "post.md")
    $postHtml = $postHtml.Replace("{{Body}}", $postMarkDownHtml.Html)
    $postHtml = $postHtml.Replace("main.css", "../../main.css")
    Set-Content -Path (Join-Path "posts" $post.Path "index.html") -Value $postHtml
}

$indexHtml = $templateHTML
$indexHtml = $indexHtml.Replace("{{Body}}", $postListHtml)
Set-Content -Path  "index.html" -Value $indexHtml