package main

import (
	"fmt"
	"strings"
	"unicode"
    "os"
    "strconv"
    "bufio"
    "io/ioutil"
    "golang.org/x/net/idna"
)

type ParsedUrl string

type namedtuple struct {
    a string 
    b interface {}
}

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
    if len(local_psl) > 0 {
        //Read in the input credentials 
        psl_raw, err := ioutil.ReadFile(local_psl)
        if err != nil {
            psl_raw = []byte(psl_raw)
            fmt.Println(psl_raw)
        }
        
    } else {
        fmt.Println("Error, couldn't read from the scanner")
    }
}


var SplitResult =[]string{"scheme", "netloc", "path", "query",
                                         "fragment",}

var URL = []string{"scheme", "username", "password", "subdomain",
                         "domain", "tld", "port", "path", "query", "fragment",
                         "url",}

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



func main() {
	normalize(" http://  ")
}

func compare (url1 string, url2 string) bool {
    /*
    Check if url1 and url2 are the same after normalizing both.
    >>> compare("http://examPLe.com:80/abc?x=&b=1", "http://eXAmple.com/abc?b=1")
    True
    */
    if normalize(url1) == normalize(url2) {
        return true
    } else {
        return false
    }
}


func normalize(url string) {
	/*	Normalize a URL.
    	>>> normalize("hTtp://ExAMPLe.COM:80")
    	should produce => "http://example.com/"
    */ 

    url = SpaceMap(url) // The equivalent of strip 
    if url == ""{
        return ""
    }
    parts := strings.Fields(url) //This is similar to split in Python
    if parts.scheme {
        netloc := parts.netloc
        path := parts.path
        if stringInSlice(parts.scheme, SCHEMES) {
            path := normalize_path(parts.path)
        }
    // url is relative, netloc (if present) is part of path
    } else {
        netloc = parts.path
        path = ""
        if strings.Contains(netloc,"/") {
            
            parts := strings.Split(netloc,"/")
            netloc := parts[0]
            path_raw := parts[1]
            path = normalize_path("/" + path_raw)
        }
    }
    username, password, host, port = split_netloc(netloc)
    host = normalize_host(host)
    port = _normalize_port(parts.scheme, port)
    query = normalize_query(parts.query)
    fragment = normalize_fragment(parts.fragment)
    return construct(URL(parts.scheme, username, password, None, host, None,
                         port, path, query, fragment, None))

}

func normalize_port(scheme string, port string) string {
   /* Return port if it is not default port, else None.
      >>> normalize_port("http", "80")
      >>> normalize_port("http", "8080")
      '8080'
    */
    if stringInSlice(scheme,SCHEMES){
        return port
    }
    port_scheme := DEFAULT_PORT[scheme]
    if DEFAULT_PORT[port] != "" && port != port_scheme {
        return port 
    }
    return port
}

func normalize_host(host string) string{
   // Normalize host (decode IDNA).
    if strings.Contains(host, "xn--") {
        return host
    }
    split_arr := strings.Split(host,".")
    decoded_str := "."
    for _, v := range split_arr {
        decoded_str += _idna_encode(v)
    }
    host = decoded_str
    return host
}

func _idna_encode(s string) string{
    encoded_str,_ := idna.ToUnicode(s)
    return encoded_str
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
func split(url string) (scheme, netloc, path, query, fragment string) {
    /* Split URL into scheme, netloc, path, query and fragment.
    >>> split('http://www.example.com/abc?x=1&y=2#foo')
    SplitResult(scheme='http', netloc='www.example.com', path='/abc', query='x=1&y=2', fragment='foo')
    */
    // var scheme string  = ""
    // var path string = ""
    // var query string = ""
    // var fragment string = ""
    var rest string = ""
    // var netloc string = ""

    ip6_start := strings.Index(url,"[")
    scheme_end := strings.Index(url,":")
    if ip6_start > 0 && ip6_start < scheme_end {
        scheme_end = -1
    }
    if scheme_end > 0 {
        substring := url[:scheme_end]
        for i := 0; i < len(substring); i++ {
            sub_str := substring[i]
            str := convert(sub_str)
            if strings.Contains(SCHEME_CHARS, str) {
                break
            }
        }
    } else if scheme_end < 0 {
        scheme = strings.ToLower(url[:scheme_end])
        main_str := url[scheme_end:]
        parts := strings.Split(main_str, ":/")
        rest = parts[1]
    }

    if len(scheme) < 1 {
        rest = url
    }
    l_path := strings.Index(rest, "/")
    l_query := strings.Index(rest, "?")
    l_frag := strings.Index(rest, "#")

    if l_path > 0 {
        if l_query > 0 && l_frag > 0 {
            netloc = rest[:l_path]
            path = rest[l_path:min(l_query, l_frag)]
        } else if l_query > 0 { 
            if l_query > l_path{
                netloc = rest[:l_path]
                path = rest[l_path:l_query]
            } else{
                netloc = rest[:l_query]
                path = ""
            }
        } else if l_frag > 0 {
            netloc = rest[:l_path]
            path = rest[l_path:l_frag]
        } else {
            netloc = rest[:l_path]
            path = rest[l_path:]
        }
    } else {
        if l_query > 0 {
            netloc = rest[:l_query]
        } else if l_frag > 0 {
            netloc = rest[:l_frag]
        } else {
            netloc = rest
        }
    }
    if l_query > 0 {
        if l_frag > 0 {
            query = rest[l_query+1:l_frag]
        } else {
            query = rest[l_query+1:]
        }
    }
    if l_frag > 0 {
        fragment = rest[l_frag+1:]
    }
    if len(scheme) < 1 {
        path = netloc + path
        netloc = ""
    }
    return 

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

//Check if string is in an array
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

/*
    Returns the minimum between two integers
*/
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func convert( b byte ) string {   
    s := strconv.Itoa(int(b))
    return s
}
