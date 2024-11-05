package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/unixpickle/essentials"
	"golang.org/x/net/html"
)

const IndexFilename = "index.html"

type ArticleDate struct {
	Year  int
	Month int
	Day   int
}

type ArticleInfo struct {
	Dir       string
	ShortName string
	Date      ArticleDate
	Title     string
}

func (a *ArticleInfo) Path() string {
	return fmt.Sprintf("/%04d/%02d/%02d/%s", a.Date.Year, a.Date.Month, a.Date.Day, a.ShortName)
}

func Articles(rootDir string) ([]*ArticleInfo, error) {
	var results []*ArticleInfo

	// Walk the directory to find all index.html files
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == IndexFilename {
			h1Content, err := extractFirstH1(path)
			if err != nil {
				return fmt.Errorf("error extracting H1 from %s: %w", path, err)
			}
			parts := strings.Split(path, string(filepath.Separator))
			dateParts := make([]int, 3)
			for i := 0; i < 3; i++ {
				part, err := strconv.Atoi(parts[len(parts)-5+i])
				if err != nil {
					return fmt.Errorf("failed to parse date from path: %s", path)
				}
				dateParts[i] = part
			}
			results = append(results, &ArticleInfo{
				Dir:       filepath.Dir(path),
				ShortName: parts[len(parts)-2],
				Date: ArticleDate{
					Year:  dateParts[0],
					Month: dateParts[1],
					Day:   dateParts[2],
				},
				Title: h1Content,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	essentials.VoodooSort(results, func(i, j int) bool {
		d1 := results[i].Date
		d2 := results[j].Date
		if d1.Year > d2.Year {
			return true
		} else if d1.Year < d2.Year {
			return false
		}
		if d1.Month > d2.Month {
			return true
		} else if d1.Month < d2.Month {
			return false
		}
		return d1.Day > d2.Day
	})

	return results, nil
}

func extractFirstH1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return "", err
	}

	var h1Text string
	var foundH1 bool

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" && !foundH1 {
			foundH1 = true
			h1Text = extractText(n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if foundH1 {
		return h1Text, nil
	}
	return "", fmt.Errorf("no <h1> found in %s", filePath)
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractText(c))
	}
	return sb.String()
}
