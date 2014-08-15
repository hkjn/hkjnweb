// Package hkjnweb is the personal websites http://www.hkjn.me /
// http://blog.hkjn.me.
//
// See https://github.com/hkjn/autosite for the framework that enables
// this site.
package hkjnweb

import (
	"appengine"

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

var web = autosite.New(
	"Henrik Jonsson",
	"pages/*.tmpl", // glob for pages
	"www.hkjn.me",  // live domain
	append(baseTemplates, "tmpl/page.tmpl"),
)

var blog = autosite.NewBlog(
	"Henrik Jonsson's blog",
	"blog/*/*/*.tmpl", // glob for blog entries
	"blog.hkjn.me",    // live domain
	append(baseTemplates, "tmpl/blog.tmpl"),
	append(baseTemplates, "tmpl/blog_listing.tmpl"),
)

var redirects = map[string]string{
	"/house": "https://docs.google.com/spreadsheets/d/1WYQErDsJMaicvA21lCrKN89KMV1wn2fZJYP1RCiozIk/edit#gid=0",
	"/where": "http://computare0.appspot.com/where/me@hkjn.me",
}

// init initializes the app.
func init() {
	// TODO: when hkjn.me/foo is accessed, we can tell based on headers
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
