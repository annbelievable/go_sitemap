package main

import(
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNs   string   `xml:"xmlns,attr"`
	Urls    []UrlLoc `xml:"url"`
}

type UrlLoc struct {
	Loc string `xml:"loc"`
}

var parsedReqUrl *url.URL
var foundLinks   []string

func main() {
	fmt.Print("Enter url to scrap the site: ")
	requestedUrl := ""
	fmt.Scanln(&requestedUrl)

	parsedReqUrl, err := url.Parse(requestedUrl)

	handlesError("Parsing url", err)
	if parsedReqUrl.Host == "" && parsedReqUrl.Path == "" {
		fmt.Print("Please enter a valid url")
		os.Exit(0)
	}
	if parsedReqUrl.Scheme == "" {
		parsedReqUrl.Scheme = "http"
	}
	startUrl := fmt.Sprintf("%s://%s%s", parsedReqUrl.Scheme, parsedReqUrl.Host, parsedReqUrl.Path)

	foundLinks = append(foundLinks, startUrl)

	for i:=0; i<len(foundLinks); i++ {
		doc := parsePage(foundLinks[i])
		findLinks(doc, parsedReqUrl)
	}

	sort.Strings(foundLinks)
	writeIntoXml()
}

func parsePage(url string) *html.Node {
	resp, err := http.Get(url)
	handlesError("Get request", err)

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	handlesError("Parsing html body", err)

	return doc
}

func findLinks(node *html.Node, parsedReqUrl *url.URL) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && len(attr.Val) > 1 {
				parsedUrl, err := parsedReqUrl.Parse(attr.Val)
				handlesError("Parsing url", err)

				if parsedReqUrl.Host == parsedUrl.Host {
					strippedUrl := fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path  )

					if isNewLink(strippedUrl) {					
						foundLinks = append(foundLinks, strippedUrl)
					}
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findLinks(child, parsedReqUrl)
	}
}

func isNewLink(newUrl string) bool {
   for _, link := range foundLinks {
	  if link == newUrl {
		 return false
	  }
   }
   return true
}

func writeIntoXml() {
	xmlContent := &UrlSet{XMLNs: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	for _, link := range foundLinks {
		xmlContent.Urls = append(xmlContent.Urls, UrlLoc{link})
	}

	file, err := xml.MarshalIndent(xmlContent, "", "    ")
	handlesError("Marshalling content to xml format", err)
	
	file = append([]byte(xml.Header), file...)
	err = ioutil.WriteFile("sitemap.xml", file, 0644)
	handlesError("Writing content into xml", err)
}

func handlesError(event string, err error) {
	if err != nil{
		fmt.Println("An error occured at: ", event)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}