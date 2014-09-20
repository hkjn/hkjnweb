// Package hkjnweb is the personal websites http://www.hkjn.me /
// http://blog.hkjn.me.
//
// See https://github.com/hkjn/autosite for the framework that enables
// this site.
package hkjnweb

import (
	"appengine"
	"net/http"

	"github.com/hkjn/autosite"
)

var baseTemplates = []string{
	"tmpl/base.tmpl",
	"tmpl/base_header.tmpl",
	"tmpl/head.tmpl",
	"tmpl/style.tmpl",
	"tmpl/fonts.tmpl",
	"tmpl/js.tmpl",
}

// aeLogger returns a pages.Logger from a request.
func aeLogger(r *http.Request) autosite.Logger {
	return appengine.NewContext(r)
}

var web = autosite.New(
	"Henrik Jonsson",
	"pages/*.tmpl", // glob for pages
	"www.hkjn.me",  // live domain
	append(baseTemplates, "tmpl/page.tmpl"),
	aeLogger,
	!appengine.IsDevAppServer(),
)

var blog = autosite.NewBlog(
	"Henrik Jonsson's blog",
	"blog/*/*/*.tmpl", // glob for blog entries
	"blog.hkjn.me",    // live domain
	append(baseTemplates, "tmpl/blog.tmpl"),
	append(baseTemplates, "tmpl/blog_listing.tmpl"),
	aeLogger,
	!appengine.IsDevAppServer(),
)

var redirects = map[string]string{
	"/house": "https://docs.google.com/spreadsheets/d/1WYQErDsJMaicvA21lCrKN89KMV1wn2fZJYP1RCiozIk",
	"/car":   "https://docs.google.com/spreadsheets/d/1Dn-2xGAtNVJ_yWW_qeYU5wMzItAOM3NeYROyWpu4zLs",
	"/where": "http://computare0.appspot.com/where/me@hkjn.me",
}

// init initializes the app.
func init() {
	// TODO: when hkjn.me/foo is accessed, we could tell based on headers
	// that the redirect came from there, and should further redirect
	// www.hkjn.me to www.hkjn.me/foo.
	if appengine.IsDevAppServer() {
		blog.ChangeURI("/", "/blogindex")
	} else {
		web.ChangeURI("/webindex", "/")
	}
	for uri, newUri := range redirects {
		web.AddRedirect(uri, newUri)
	}

	web.Register()
	blog.Register()
}
