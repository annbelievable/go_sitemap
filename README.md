# Sitemap builder

This sitemap builder is implemented in golang. To use the sitemap builder, clone the repository into a folder and run the following command inside the folder:

```
go run main.go
```

It will then prompt you to enter the url. After entering the url, it will take some time to scrape the site. During the process, any fragments (the part with # symbol) or query strings (the part with ? symbol) will be removed from the url. At the end of the process, it will create a file "sitemap.xml".

Included in this repo is an example of the sitemap built from the site "https://golang.org". There are 6924 urls found for the site.