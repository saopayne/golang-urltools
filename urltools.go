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

type TypeError interface {
    Errors() string
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

var _hextochr = make(map[string]string, 256)

func main() {
    for i := 0; i < 256; i++ {
        _hextochr(i) = "%02x" % i      
    }
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


func normalize(url string) string {
	/*	Normalize a URL.
    	>>> normalize("hTtp://ExAMPLe.COM:80")
    	should produce => "http://example.com/"
    */ 

    url = SpaceMap(url) // The equivalent of strip 
    if url == ""{
        return ""
    }
    parts := strings.Fields(url) //This is similar to split in Python
    if parts[0] != "" {
        netloc := parts[1]
        path := parts[2]
        if stringInSlice(parts[0], SCHEMES) {
            path := normalize_path(parts[2])
        }
    // url is relative, netloc (if present) is part of path
    } else {
        netloc := parts[2]
        path := ""
        if strings.Contains(netloc,"/") {
            parts_tmp := strings.Split(netloc,"/")
            netloc := parts_tmp[0]
            path_raw := parts_tmp[1]
            path := normalize_path("/" + path_raw)
        }
    }
    username, password, host, port := split_netloc(netloc)
    host = normalize_host(host)
    port = normalize_port(parts.scheme, port)
    query := normalize_query(parts.query)
    fragment := normalize_fragment(parts.fragment)
    return construct(URL(parts.scheme, username, password, None, host, None,
                         port, path, query, fragment, None))

}

func normalize_path(path string){
    /*  Normalize path: collapse etc.
        >>> normalize_path("/a/b///c")
        "/a/b/c"
    */
    path_list := []string{"//", "/", "",}
    if stringInSlice(path, path_list){
        return "/"
    }
    npath := normpath(unquote(path, QUOTE_EXCEPTIONS["path"]))
    if path[-1] == "/" && npath != "/"{
        npath += "/"
    }
    return npath
}


func normalize_query(query string){
    /*
        Normalize query: sort params by name, remove params without value.
        >>> normalize_query("z=3&y=&x=1")
        "x=1&z=3"
    */
    if query == "" || len(query) <= 2{
        return ""
    }
    nquery = unquote(query, QUOTE_EXCEPTIONS["query"])
    params = nquery.split("&")
    nparams = [] string
    for i := 0; i < params.length; i++ {
        if strings.Contains("=", params[i]){
            s = strings.Split(params[i], "=")
            k, v = s[0], s[1]
            if k!=nil && v!= nil {
                nparams = append(nparams,"%s=%s" % k, v)
            }
        }
    }
    sort(nparams)
    return strings.Join(nparams, "&")
}

func normalize_fragment(fragment string) {
    /* Normalize fragment (unquote with exceptions only).*/
    return unquote(fragment, QUOTE_EXCEPTIONS["fragment"])
}


// _hextochr = dict(('%02x' % i, chr(i)) for i in range(256))
// _hextochr.update(dict(('%02X' % i, chr(i)) for i in range(256)))


func normalize_port(scheme string, port string) string {
   /* Return port if it is not default port, else None.
      >>> normalize_port("http", "80")
      >>> normalize_port("http", "8080")
      '8080'
    */
    if stringInSlice(scheme, SCHEMES){
        return port
    }
    port_scheme := DEFAULT_PORT[scheme]
    if DEFAULT_PORT[port] != "" && port != port_scheme {
        return port 
    }
    return port
}

func normalize_host(host string) string {
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

func _idna_encode(s string) string {
    encoded_str,_ := idna.ToUnicode(s)
    return encoded_str
}

func unquote(text string, exceptions []string) string {
    /*
        Unquote a text but ignore the exceptions.
        >>> unquote(foo%23bar")
        "foo#bar"
        >>> unquote("foo%23bar", ["#"])
        "foo%23bar"
    */
    if text == nil{
        if text != ""{
            return errors.New("None object cannot be unquoted")
        } else {
            return text
        }
    }
    if !strings.Contains(text, "%") {
        return text
    }
    s := strings.Split(text, "%")
    res := [] string{s[0]}
    for _,h := s[1:] {
        c = _hextochr.get(h[:2])
        if c != nil && !(strings.Contains(exceptions, c)){
            if len(h) > 2{
                res.append(c + h[2:])
            } else {
                res.append(c)
            }
        } else{
            res.append("%" + h)
        }
    }
    return strings.Join(res, "")
}

func parse(url string){
    /*
        Parse a URL.
        >>> parse('http://example.com/foo/')
        URL(scheme='http', ..., domain='example', tld='com', ..., path='/foo/', ...)
    */
    parts = split(url)
    if parts.scheme {
        username, password, host, port = split_netloc(parts.netloc)
        subdomain, domain, tld = split_host(host)
    } else {
        username := "" 
        password := ""
        subdomain := ""
        domain := ""
        tld := ""
        port := ""
    }
    return URL(parts.scheme, username, password, subdomain, domain, tld,
               port, parts.path, parts.query, parts.fragment, url)
}


func extract(url string){
    /*
        Extract as much information from a (relative) URL as possible.
        >>> extract('example.com/abc')
        URL(..., domain='example', tld='com', ..., path='/abc', ...)
    */
    parts = split(url)
    if parts.scheme{
        netloc := parts.netloc
        path := parts.path
    } else {
        netloc := parts.path
        path := ""
        if strings.Contains(netloc, "/") {
            netloc, path_raw = netloc.split("/", 1)
            path = "/" + path_raw
        }
    }
    username, password, host, port = split_netloc(netloc)
    subdomain, domain, tld = split_host(host)
    return URL(parts.scheme, username, password, subdomain, domain, tld,
               port, path, parts.query, parts.fragment, url)

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
    >>> split("http://www.example.com/abc?x=1&y=2#foo")
    SplitResult(scheme="http", netloc="www.example.com", path="/abc", query="x=1&y=2", fragment="foo")
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

func _clean_netloc(netloc string){
    /*
        Remove trailing "." and ":"" and tolower.
        >>> _clean_netloc("eXample.coM:")
        "example.com"
    */
    try:
        return netloc.rstrip(".:").lower()
    except:
        return netloc.rstrip(".:").decode("utf-8").lower().encode("utf-8")
}

func split_netloc(netloc string) (string, string, string, string) {
    /*
        Split netloc into username, password, host and port.
        >>> split_netloc("foo:bar@www.example.com:8080")
        ("foo", "bar", "www.example.com", "8080")
    */
    username := ""
    password := ""
    host := ""
    port := ""
    if strings.Contains(netloc, "@") {
        user_pw, netloc = netloc.split("@", 1)
        if strings.Contains(user_pw, ":"){
            username, password = user_pw.split(":", 1)
        } else {
            username = user_pw
        }
    }
    netloc = _clean_netloc(netloc)
    if strings.Contains(netloc, ":") && netloc[-1] != "]"{
        host, port = netloc.rsplit(":", 1)
    } else {
        host = netloc
    }
    return username, password, host, port
}


func split_host(host string) (string, string, string){
    /*
        Use the Public Suffix List to split host into subdomain, domain and tld.
        >>> split_host("foo.bar.co.uk")
        ("foo", "bar", "co.uk")
    */
    // host is IPv6?
    if strings.Contains(host, "["){
        return "", host, ""
    }
    // host is IPv4?
    for c := host{
        if !strings.Contains(IP_CHARS, c){
            break
        } else {
        return "", host, ""
    }
    }
    // host is a domain name
    domain := ""
    subdomain := "" 
    tld := ""
    parts = host.split(".")
    for i := range(len(parts)){
        tld = ".".join(parts[i:])
        wildcard_tld = "*." + tld
        exception_tld = "!" + tld
        if strings.Contains(PSL, exception_tld) {
            domain = ".".join(parts[:i+1])
            tld = ".".join(parts[i+1:])
            break
        }
        if strings.Contains(PSL, tld) {
            domain = ".".join(parts[:i])
            break
        }
        if strings.Contains(PSL, wildcard_tld) {
            domain = ".".join(parts[:i-1])
            tld = ".".join(parts[i-1:])
            break
        }
    }
    if strings.Contains(domain, ".") {
        subdomain, domain = domain.rsplit(".", 1)
    }

    return subdomain, domain, tld
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
