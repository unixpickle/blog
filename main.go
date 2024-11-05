package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/unixpickle/essentials"
)

type Server struct {
	AssetDir   string
	ArticleDir string
}

func main() {
	var server Server
	var port int
	flag.StringVar(&server.AssetDir, "assets", "assets", "asset directory")
	flag.StringVar(&server.ArticleDir, "articles", "articles", "article directory")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	fs := http.FileServer(http.Dir(server.AssetDir))
	http.Handle("/favicon.ico", fs)
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || r.URL.Path == "/" {
			if err := server.ServeIndex(w, r); err != nil {
				internalError(w, err)
			}
		} else {
			if err := server.ServeArticle(w, r); err != nil {
				internalError(w, err)
			}
		}
	})

	essentials.Must(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func (s *Server) ServeIndex(w http.ResponseWriter, r *http.Request) error {
	articles, err := Articles(s.ArticleDir)
	if err != nil {
		return err
	}

	indexTemplate, err := os.ReadFile(filepath.Join(s.AssetDir, "index.template"))
	if err != nil {
		return err
	}
	tmplParsed, err := template.New("webpage").Parse(string(indexTemplate))
	if err != nil {
		return err
	}
	return tmplParsed.Execute(w, articles)
}

func (s *Server) ServeArticle(w http.ResponseWriter, r *http.Request) error {
	articles, err := Articles(s.ArticleDir)
	if err != nil {
		return err
	}

	absPath := path.Clean(r.URL.Path)
	for _, article := range articles {
		if strings.HasPrefix(absPath, article.Path()) {
			return s.serveArticle(w, r, article, absPath[len(article.Path()):])
		}
	}

	http.NotFound(w, r)
	return nil
}

func (s *Server) serveArticle(
	w http.ResponseWriter,
	r *http.Request,
	article *ArticleInfo,
	relpath string,
) error {
	if r.URL.Path == article.Path() {
		newURL := *r.URL
		newURL.Path += "/"
		http.Redirect(w, r, r.URL.String()+"/", http.StatusPermanentRedirect)
		return nil
	}

	if strings.HasPrefix(relpath, "/") {
		relpath = relpath[1:]
	}
	if relpath == "" {
		relpath = IndexFilename
	}
	localPath := article.Dir
	for _, comp := range strings.Split(relpath, "/") {
		localPath = filepath.Join(localPath, comp)
	}
	stat, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	if relpath == IndexFilename {
		contents, err := os.ReadFile(localPath)
		if err != nil {
			return err
		}
		headSuffix, err := os.ReadFile(filepath.Join(s.AssetDir, "head_suffix.html"))
		if err != nil {
			return err
		}
		bodyPrefix, err := os.ReadFile(filepath.Join(s.AssetDir, "body_prefix.html"))
		if err != nil {
			return err
		}
		titleTag := "<title>" + html.EscapeString(article.Title) + "</title>"
		dateTag := ("<div class=\"article-date\">Posted on " +
			fmt.Sprintf("%02d/%02d/%04d", article.Date.Month, article.Date.Day, article.Date.Year) +
			"</div>")
		replacements := map[string]string{
			"</head>": titleTag + string(headSuffix) + "</head>",
			"<body>":  "<body>" + string(bodyPrefix),
			"</h1>":   "</h1>" + dateTag,
		}
		for old, new := range replacements {
			contents = []byte(strings.Replace(string(contents), old, new, 1))
		}
		contents = []byte(strings.Replace(string(contents), "\\vec{", "\\overrightarrow{", -1))
		http.ServeContent(w, r, filepath.Base(localPath), stat.ModTime(), bytes.NewReader(contents))
		return nil
	} else {
		f, err := os.Open(localPath)
		if err != nil {
			return err
		}
		defer f.Close()
		http.ServeContent(w, r, filepath.Base(localPath), stat.ModTime(), f)
		return nil
	}
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}
