package main

import (
	"fmt"
	"strings"
	"unicode"
    "os"
    "bufio"
    "io"
    "golang.org/x/net/idna"
)

type ParsedUrl string

var __all__ = []string {"URL", "SplitResult", "parse", "extract", "construct", "normalize",
           "compare", "normalize_host", "normalize_path", "normalize_query",
           "normalize_fragment", "encode", "unquote", "split", "split_netloc",
           "split_host",}


var PSL_URL string = "https://publicsuffix.org/list/effective_tld_names.dat"

func get_public_suffix_list() {
    /*
    Return a set containing all Public Suffixes.
    If the env variable PUBLIC_SUFFIX_LIST does not point to a local copy of the
    public suffix list it is downloaded into memory each time urltools is
    imported.
    */
    local_psl := os.Getenv("PUBLIC_SUFFIX_LIST")
    if local_psl.length > 0 {
        //Read in the input credentials 
        input := bufio.NewScanner(io.Reader(local_psl))
        psl_raw := input.Scan()
        fmt.Println(psl_raw)

    } else {
        fmt.Println("Error, couldn't read from the scanner")
    }
}


var SplitResult = namedtuple("SplitResult", []string{"scheme", "netloc", "path", "query",
                                         "fragment"})
var URL = namedtuple("URL", []string{"scheme", "username", "password", "subdomain",
                         "domain", "tld", "port", "path", "query", "fragment",
                         "url",})



var SCHEMES = []string {"http", "https", "ftp", "sftp", "file", "gopher", "imap", "mms",
           "news", "nntp", "telnet", "prospero", "rsync", "rtsp", "rtspu",
           "svn", "git", "ws", "wss",}
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

type namedtuple struct {
    a , b interface{}
}

func main() {
	normalize(" http://  ")
}

func normalize(url string) {
	/*	Normalize a URL.
    	>>> normalize("hTtp://ExAMPLe.COM:80")
    	should produce => "http://example.com/"
    */ 
    url = SpaceMap(url)
    fmt.Println(url)

}

func normalize_port(scheme string, port string) {
   /* Return port if it is not default port, else None.
    >>> normalize_port('http', '80')
    >>> normalize_port('http', '8080')
    '8080'
    */
    if SCHEMES[scheme]{
        return port
    }
    if DEFAULT_PORT[port] && port != DEFAULT_PORT[scheme]{
        return port
    }
}

func normalize_host(host string) {
   // Normalize host (decode IDNA).
    if strings.Contains(host, "xn--") {
        return host
    }
    split_arr := strings.Split(host,'.')
    decoded_str := '.'
    for _, v := range split_arr {
        decoded_str += _idna_decode(decoded_split)
    }
}

func _idna_encode(s string){
    return idna.ToUnicode(s)
}

func parameterize_url() {
    
}

// Read the file and convert to lines
func File2lines(filePath string) []string {
      f, err := os.Open(filePath)
      if err != nil {
              panic(err)
      }
      defer f.Close()

      var lines []string
      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
              lines = append(lines, scanner.Text())
      }
      if err := scanner.Err(); err != nil {
              fmt.Fprintln(os.Stderr, err)
      }

      return lines
}

/*
 Splits the url into the constituent elements 
*/
func split(url string) {
    /* Split URL into scheme, netloc, path, query and fragment.
    >>> split('http://www.example.com/abc?x=1&y=2#foo')
    SplitResult(scheme='http', netloc='www.example.com', path='/abc', query='x=1&y=2', fragment='foo')
    */
    var scheme string  = ""
    var path string = ""
    var query string = ""
    var fragment string = ""

    ip6_start = url.find('[')
    scheme_end = url.find(':')
    if ip6_start > 0 && ip6_start < scheme_end {
        scheme_end = -1
    }
    if scheme_end > 0 {
        for c in url[:scheme_end]{
            if !SCHEME_CHARS[c]{
                break
            }
        }
        else{
            scheme = url[:scheme_end].lower()
            rest = url[scheme_end:].lstrip(':/')
        }
    }
    if not scheme:
        rest = url
    l_path = rest.find('/')
    l_query = rest.find('?')
    l_frag = rest.find('#')
    if l_path > 0:
        if l_query > 0 and l_frag > 0:
            netloc = rest[:l_path]
            path = rest[l_path:min(l_query, l_frag)]
        elif l_query > 0:
            if l_query > l_path:
                netloc = rest[:l_path]
                path = rest[l_path:l_query]
            else:
                netloc = rest[:l_query]
                path = ''
        elif l_frag > 0:
            netloc = rest[:l_path]
            path = rest[l_path:l_frag]
        else:
            netloc = rest[:l_path]
            path = rest[l_path:]
    else:
        if l_query > 0:
            netloc = rest[:l_query]
        elif l_frag > 0:
            netloc = rest[:l_frag]
        else:
            netloc = rest
    if l_query > 0:
        if l_frag > 0:
            query = rest[l_query+1:l_frag]
        else:
            query = rest[l_query+1:]
    if l_frag > 0:
        fragment = rest[l_frag+1:]
    if not scheme:
        path = netloc + path
        netloc = ''
    return SplitResult(scheme, netloc, path, query, fragment)

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