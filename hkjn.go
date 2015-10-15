// Package hkjnweb is the personal website hkjn.me.
//
// See https://github.com/hkjn/autosite for the framework that enables
// this site.
package hkjnweb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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
	isProd = false
	Logger = autosite.Glogger{}
)

// getLogger returns the autosite.Logger to use.
func getLogger(r *http.Request) autosite.Logger {
	return Logger
}

var (
	webDomain = "www.hkjn.me"
	web       = autosite.New(
		"Henrik Jonsson",
		"pages/*.tmpl", // glob for pages
		webDomain,      // live domain
		append(baseTemplates, "tmpl/page.tmpl"),
		getLogger,
		isProd,
		template.FuncMap{},
	)

	goImportTmpl = `<head>
    <meta http-equiv="refresh" content="0; URL='%s'">
    <meta name="go-import" content="%s git %s">
  </head>`

	redirects = map[string]string{
		"/where": "http://computare0.appspot.com/where/me@hkjn.me",
	}
)

// Register registers the handlers.
func Register(prod bool) {
	isProd = prod
	// We remap /webindex (from pages/webindex.tmpl, no special-cases in
	// the autosite package) to serve index.

	// TODO(henrik): Re-enable this once autosite package doesn't
	// interfere with / and /keybase.txt handlers below.

	// web.ChangeURI("/webindex", "/")
	http.HandleFunc("/keybase.txt", keybaseHandler)
	if isProd {
		log.Println("We're in prod, remapping some paths\n")
		http.HandleFunc("/", nakedIndexHandler)
	} else {
		log.Println("We're not in prod, remapping some paths\n")
		http.HandleFunc("/nakedindex", nakedIndexHandler)
	}
	for uri, newUri := range redirects {
		web.AddRedirect(uri, newUri)
	}

	web.Register()
}

// nakedIndexHandler serves requests to hkjn.me/
func nakedIndexHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)
	l.Infof("nakedIndexHandler for URI %s\n", r.RequestURI)
	if r.URL.Path == "/" {
		url := "/webindex"
		l.Debugf("visitor to / of naked domain, redirecting to %q..\n")
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// Our response tells the `go get` tool where to find
		// `hkjn.me/[package]`.
		parts := strings.Split(r.URL.Path, "/")
		godocUrl := fmt.Sprintf("https://godoc.org/hkjn.me%s", r.URL.Path)
		repoRoot := fmt.Sprintf("https://github.com/hkjn/%s", parts[1])
		importPrefix := fmt.Sprintf("hkjn.me/%s", parts[1])
		fmt.Fprintf(w, goImportTmpl, godocUrl, importPrefix, repoRoot)
	}
}

var keybaseVerifyText = `
==================================================================
https://keybase.io/hkjn
--------------------------------------------------------------------

I hereby claim:

  * I am an admin of http://www.hkjn.me
  * I am hkjn (https://keybase.io/hkjn) on keybase.
  * I have a public key with fingerprint D618 7A03 A40A 3D56 62F5  4B46 03EF BF83 9A5F DC15

To do so, I am signing this object:

{
    "body": {
        "key": {
            "fingerprint": "d6187a03a40a3d5662f54b4603efbf839a5fdc15",
            "host": "keybase.io",
            "key_id": "03efbf839a5fdc15",
            "kid": "0101c778a84ebd54581e354296aa8171d79a23a3bcce4dbf99f880dfeafc2e0030030a",
            "uid": "d1a36099363d76e027441cb534d6ba19",
            "username": "hkjn"
        },
        "service": {
            "hostname": "www.hkjn.me",
            "protocol": "http:"
        },
        "type": "web_service_binding",
        "version": 1
    },
    "ctime": 1429520487,
    "expire_in": 157680000,
    "prev": "4e9b7b9ff1b50e84436eb96100528ef25bf277aca7169d5a4330b2f5bd0b9ba4",
    "seqno": 4,
    "tag": "signature"
}

which yields the signature:

-----BEGIN PGP MESSAGE-----
Version: GnuPG v2

owFtUntQVFUYZ3kYD0lRVDKYhqsFNQT3/SBxjcjMmFEYGq3M5Zx7z11uq3vX3WUB
mZWc0olIpxQQMGLMFZKhWSie0kCORUAzIBI6haat2AOzZIAdiqTuJfuvM2fmzPd9
v9/vO79zvncjgwJCDUG/975RY7q6yvD1+Q0w4EX6c7UYg6pUhKUWYxa0eMiK1Yzs
NrtidWKpmMQSPAdwCtA4oCSGZUmZoSHN4hSSocxTAmBkSSQYLAnLUx06Q5OBwIGS
FVXLaYFJkbTs/+At/xYInBA5jgc8jaDE0AxPIIqhSYEFgCc4QuIEQFKAgqKIaAnK
giDzPC7JCMgiiXCc0jfQ5PIX5SQCUCwuCBRLSRyLcJKjaUKEDEVLLASEoAMdyG4F
e5GGzrO8ZsXcSZiWcSki0u3rLu5XCwoKknVEshYlYTa76lRFdY9OczptqTrPWWRb
BCJoui9hgopV0l5QI7iQ3aGoViyV0JCiU9E1Cc0YQ+I0zyVhqNCm2JFJ0REMx/K4
tvQ2yKVJ0kiAHBRkmYAMjniaplgEBZbAcYbkkUwyUCY5DoiAI1hBYgBNUTjU/gZK
OBQgoDHd1D6riqXS2jWBWZN0KGYrcObbEeYOLw2MDQ4whAYsCQnUhyAgPGz5f6Ox
OnvZ33svHP/Asy2uePVj/exDM9G+c5PTt/9qGsw9NFuSVtCYs3n+07rymV8MXqph
rqWyom+lGUusjzsyeHnTwQ0nhIUtU7fIE+q9+k2vu8Nmz300FwPvlOc2lu35Nelm
SKP/8Yoctb2jutYm+f1/RJUrJ1u+TauybPnqmdOWY8YuC21q6B15+MaVl41vBceM
LdT9GZGYHpmRkHCrNNC4YHQ0p7vK/N3vdwx18S8czBRPRdyMSombuF7iKb2gLLVX
95EzZs4zFasaW7O8Q0G1w76JgmVhWTl936xtPftAyLX8s08E5pze5socSPmtvrdJ
aQ3N7jhTtXKhjVvDgC5ff13X+azy2HVqrvvizmezxjPXq90D9jPfe8eMlU3ckh1X
9v98d+DwqynTGVtrhtfecAyCivA7ZSOy7Bx9ZPyu7/i1pU/91JcSJb132/C84Z0s
zyfHWkp2rZhr+9FW0t4+2rzrVPWStu1uXLme5vJOfsbXfFEc720E61YcGHAk7Gbr
QHLhKn+bb7Tnw/pHf+BcO3Y7Nz/Y3NlZVXZk/87DQbVF3yG7GTw3wk/tk0q+vBzs
jqlkE6+exDOmn1Yj52X/fHvP8PjR2dnEhjWO7ob18mT/vbcrhy5u7UHLKydemulo
/7i3U5rwxl/Kto85C73p+a9EiBt9TwYf9XjrD0XHtx0wed7M265sjL70Dw==
=c4ep
-----END PGP MESSAGE-----

And finally, I am proving ownership of this host by posting or
appending to this document.

View my publicly-auditable identity here: https://keybase.io/hkjn

==================================================================`

func keybaseHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)
	l.Infof("keybaseHandler for URI %s\n", r.RequestURI)
	fmt.Fprintf(w, keybaseVerifyText)
}
