# golang-urltools

Some functions to parse and normalize URLs.

The main focus of this library is to make it possible to work on all segments of an URL. Thus a core feature (which is not provided by stdlib) is to split a domain name correctly by using the Public Suffix List (see below).
