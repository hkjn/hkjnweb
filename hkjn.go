// Package hkjnweb is the personal website http://www.hkjn.me / http://blog.hkjn.me.
package hkjnweb

import (
	"appengine"

	"github.com/hkjn/autosite"
)

var baseTemplates = []string{
	"tmpl/base.tmpl",
	"tmpl/base_header.tmpl",
	"tmpl/style.tmpl",
	"tmpl/fonts.tmpl",
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

// init initializes the app.
func init() {
	if appengine.IsDevAppServer() {
		blog.ChangeURI("/", "/blogindex")
	} else {
		web.ChangeURI("/webindex", "/")
	}

	web.Register()
	blog.Register()
}
