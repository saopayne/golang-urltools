package main

import (
	"fmt"
	"strings"
	"unicode"
)

type URL string 
type ParsedUrl string


var SCHEMES = []string {"http", "https", "ftp", "sftp", "file", "gopher", "imap", "mms",
           "news", "nntp", "telnet", "prospero", "rsync", "rtsp", "rtspu",
           "svn", "git", "ws", "wss"}
var SCHEME_CHARS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var IP_CHARS string = "0123456789.:"
var DEFAULT_PORT = map[string]string {
    "http": "80",
    "https": "443",
    "ws": "80",
    "wss": "443",
    "ftp": "21",
    "sftp": "22",
    "ldap": "389",
}
var QUOTE_EXCEPTIONS = map[string]string {
    "path": " /?+#",
    "query": " &=+#",
    "fragment": " +#",
}

func main() {
	normalize(" http://  ")
}

func normalize(url string) {
	/*	Normalize a URL.
    	>>> normalize('hTtp://ExAMPLe.COM:80')
    	should produce => 'http://example.com/'
    */ 
    url = SpaceMap(url)
    fmt.Println(url)

}

// This is equivalent to strip() which removes all trailing white spaces
func SpaceMap(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, str)
}