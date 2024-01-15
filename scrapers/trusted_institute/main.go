package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

const (
	BASE_URL          = "https://trustedinstitute.com/"
	PRACTICE_URL      = "https://trustedinstitute.com/practice"
	SUBJECT_URL_REGEX = "practice/."
	QUIZ_URL_REGEX    = "(practice|topic)/."
)

var (
	_binPaths = []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	}
)

func main() {
	u := launcher.New().
		Headless(false).
		Bin(_binPaths[0]).
		MustLaunch()

	browser := rod.New().
		ControlURL(u).
		Timeout(time.Minute).
		MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	links := nav(page, PRACTICE_URL, SUBJECT_URL_REGEX)

	for _, l := range links {
		quizLinks := nav(page, l, QUIZ_URL_REGEX)
		for _, ll := range quizLinks {
			fmt.Println(ll)
		}
	}

	time.Sleep(time.Second * 15)
}

func nav(page *rod.Page, uri, regex string) (uris []string) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" {
		uri, _ = url.JoinPath(BASE_URL, uri)
	}
	page.MustNavigate(uri).MustWaitNavigation()()
	links := page.MustElements("a")

	rx := regexp.MustCompile(regex)
	for _, link := range links {
		href, err := link.Attribute("href")
		if err != nil {
			log.Fatal(err)
		}

		if rx.MatchString(*href) {
			uris = append(uris, *href)
		}
	}

	return
}
