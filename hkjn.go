// Package hkjnweb is the personal websites http://www.hkjn.me /
// http://blog.hkjn.me.
//
// See https://github.com/hkjn/autosite for the framework that enables
// this site.
package hkjnweb

import (
	"fmt"
	"net/http"

	"appengine"

	"html/template"

	"hkjn.me/autosite"
)

var (
	baseTemplates = []string{
		"tmpl/base.tmpl",
		"tmpl/base_header.tmpl",
		"tmpl/head.tmpl",
		"tmpl/style.tmpl",
		"tmpl/fonts.tmpl",
		"tmpl/js.tmpl",
	}
	live = !appengine.IsDevAppServer()
)

// aeLogger returns a pages.Logger from a request.
func aeLogger(r *http.Request) autosite.Logger {
	return appengine.NewContext(r)
}

var (
	webDomain  = "www.hkjn.me"
	blogDomain = "blog.hkjn.me"
	web        = autosite.New(
		"Henrik Jonsson",
		"pages/*.tmpl", // glob for pages
		webDomain,      // live domain
		append(baseTemplates, "tmpl/page.tmpl"),
		aeLogger,
		live,
		tmplFuncs(webDomain),
	)

	blog = autosite.NewBlog(
		"Henrik Jonsson's blog",
		"blog/*/*/*.tmpl", // glob for blog entries
		blogDomain,        // live domain
		append(baseTemplates, "tmpl/blog.tmpl"),
		append(baseTemplates, "tmpl/blog_listing.tmpl"),
		aeLogger,
		live,
		tmplFuncs(blogDomain),
	)

	redirects = map[string]string{
		"/where": "http://computare0.appspot.com/where/me@hkjn.me",
	}
)

// tmplFuncs returns extra template functions.
func tmplFuncs(domain string) template.FuncMap {
	return template.FuncMap{
		"live": func() bool { return live },
		"domain": func() string {
			if live {
				return domain
			}
			return ""
		},
	}
}

// init initializes the app.
func init() {
	if live {
		web.ChangeURI("/webindex", "/")
		http.HandleFunc("hkjn.me/", nakedIndexHandler)
	} else {
		blog.ChangeURI("/", "/blogindex")
		http.HandleFunc("/nakedindex", nakedIndexHandler)
	}
	for uri, newUri := range redirects {
		web.AddRedirect(uri, newUri)
	}

	web.Register()
	blog.Register()
}

func nakedIndexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("nakedIndexHandler for URI %s\n", r.RequestURI)
	q := r.URL.Query()
	if q.Get("go-get") != "" {
		repoRoot := fmt.Sprintf("https://github.com/hkjn%s", r.URL.Path)
		importPrefix := fmt.Sprintf("hkjn.me%s", r.URL.Path)
		c.Debugf("got go-get param, assuming it's `go get` and pointing it to %s\n", repoRoot)
		fmt.Fprintf(w, `<meta name="go-import" content="%s git %s">\n`, importPrefix, repoRoot)
	} else {
		nextURL := ""
		if r.URL.Path == "/" {
			nextURL = fmt.Sprintf("http://%s", webDomain)
			c.Debugf("regular visitor to naked domain, assuming they're here for www and redirecting to %s..\n", nextURL)
		} else {
			nextURL = fmt.Sprintf("https://godoc.org/hkjn.me/%s", r.URL.Path)
			c.Debugf("regular visitor to naked domain, assuming they're here for go packages and redirecting to %s..\n", nextURL)
		}
		http.Redirect(w, r, nextURL, http.StatusSeeOther)
	}
}
