package util

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/icholy/replace"
	"golang.org/x/text/transform"
	"mvdan.cc/xurls/v2"
)

var (
	urlRegex = xurls.Strict()

	noScriptEnter = replace.String("<noscript>", "")
	noScriptLeave = replace.String("</noscript>", "")

	noScriptReplacer = transform.Chain(noScriptEnter, noScriptLeave)
)

// FindCanonicalURL returns the canonical URL of an article if possible,
// else the original input is returned
func FindCanonicalURL(url string, secondTry bool) (out string) {
	out = url

	client := http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// If we were redirected, we now get a different URL
			out = req.URL.String()
			return nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", GetUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.7,en;q=0.3")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return
	}

	// t.co doesn't use real redirects. And the <meta http-equiv="refresh" tag is inside a <noscript>
	// tag which the Go html parser doesn't parse (because it's in script mode and thus parses that node as text...)
	// So here is the stupid but working solution: replacing any "<noscript>" "</noscript>" text
	bodyReader := transform.NewReader(bufio.NewReader(resp.Body), noScriptReplacer)

	// Still try to find the tag
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	canon := doc.Find("[rel=canonical]").First()
	if canon.Length() != 0 {
		return canon.AttrOr("href", url)
	}

	equiv := doc.Find("[http-equiv]").First()
	if equiv.Length() != 0 {
		u := urlRegex.FindString(equiv.AttrOr("content", ""))
		if u != "" && u != url && !secondTry {
			return FindCanonicalURL(u, true)
		}
	}

	return
}

type ArticleStore struct {
	// SeenArticles maps an article URL to the last time it was seen
	SeenArticles map[string]time.Time

	filename     string
	keepDuration time.Duration
}

func NewArticleStore(filename string, keepDuration time.Duration) (*ArticleStore, error) {
	var a = &ArticleStore{
		SeenArticles: make(map[string]time.Time),
		filename:     filename,
		keepDuration: keepDuration,
	}
	err := LoadJSON(filename, a)
	if !os.IsNotExist(err) {
		return a, err
	}
	return a, nil
}

func (a *ArticleStore) save() error {
	for key, date := range a.SeenArticles {
		if time.Since(date) > a.keepDuration {
			delete(a.SeenArticles, key)
		}
	}
	return SaveJSON(a.filename, a)
}

func (a *ArticleStore) MarkSeen(url, canonicalArticleURL string) error {
	now := time.Now()
	a.SeenArticles[url] = now
	a.SeenArticles[canonicalArticleURL] = now

	return a.save()
}

func (a *ArticleStore) HasSeen(url string) (seen bool, canonical string, err error) {
	if t, ok := a.SeenArticles[url]; ok && time.Since(t) < a.keepDuration {
		return true, "", nil
	}

	can := FindCanonicalURL(url, false)
	if can == url {
		return false, "", nil
	}

	if t, ok := a.SeenArticles[can]; ok && time.Since(t) < a.keepDuration {
		return true, can, nil
	}

	return false, can, nil
}

func (a *ArticleStore) ShouldIgnoreTweet(tweet *twitter.Tweet) (ignore bool) {

	// Get the text *with* URLs
	var textWithURLs = tweet.SimpleText
	if textWithURLs == "" {
		textWithURLs = tweet.FullText
	}

	// Find all URLs
	urls := urlRegex.FindAllString(textWithURLs, -1)

	// Now check if any of these URLs is ignored
	for _, u := range urls {
		seen, can, err := a.HasSeen(u)
		if err == nil {
			continue
		}
		if seen {
			return true
		}

		err = a.MarkSeen(u, can)
		if err != nil {
			log.Printf("[ArticleStore] Saving article: %s\n", err.Error())
		}
	}

	return false
}
