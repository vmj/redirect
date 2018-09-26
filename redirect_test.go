package main

import (
	"net/url"
	"testing"
)

func TestMakeLocation(t *testing.T) {
	var tests = []struct {
		protocol string
		host string
		port string
		path string
		query string
		original *url.URL
		expected string
	} {
		{"#{protocol}", "#{host}", "#{port}", "#{path}", "#{query}",
			parse("http://localhost:8080/foo?bar=baz#quux", t),
			"http://localhost:8080/foo?bar=baz#quux",
		},
		{"https", "#{host}", "", "#{path}", "#{query}",
		parse("http://localhost:8080", t),
		"https://localhost/",
		},
		{"https", "slackware.fi", "443", "/home", "http=1",
			parse("http://localhost:8080/foo?bar=baz", t),
			"https://slackware.fi/home?http=1",
		},
		{"#{protocol}", "#{host}", "#{port}", "#{path}", "#{query}",
			parse("/just-path", t),
			"/just-path",
		},
		{"#{protocol}", "#{host}", "#{port}", "#{path}", "#{query}",
			parse("//me:pw@localhost/path", t),
			"//me:pw@localhost/path",
		},
		{"#{protocol}", "#{host}", "443", "#{path}", "#{query}",
			parse("/just-path", t),
			"//:443/just-path",
		},
		{"#{protocol}", "slackware.fi", "#{port}", "#{path}", "#{query}",
			parse("/just-path", t),
			"//slackware.fi/just-path",
		},
		{"#{protocol}", "slackware.fi", "#{port}", "#{path}", "#{query}",
			parse("//:80/just-path", t),
			"//slackware.fi:80/just-path",
		},
		{"#{protocol}", "#{host}", "#{port}", "/#{host}/#{port}#{path}", "#{query}",
			parse("http://localhost:8080/foo?bar=baz", t),
			"http://localhost:8080/localhost/8080/foo?bar=baz",
		},
		{"#{protocol}", "#{host}", "#{port}", "#{path}", "host=#{host}&port=#{port}&path=#{path}&#{query}",
			parse("http://localhost:8080/foo?bar=baz", t),
			"http://localhost:8080/foo?host=localhost&port=8080&path=%2Ffoo&bar=baz",
		},
	}
	for _, test := range tests {
		if actual := makeLocation(test.protocol, test.host, test.port, test.path, test.query, test.original); actual != test.expected {
			t.Errorf("makeLocation(%s, %s, %s, %s, %s, %s) = %s, expected %s",
				test.protocol, test.host, test.port, test.path, test.query, test.original, actual, test.expected)
		}
	}
}

func parse(rawurl string, t *testing.T) *url.URL {
	p, e := url.Parse(rawurl)
	if e != nil {
		t.Fatalf("Failed to parse URL: %s", rawurl)
	}
	return p
}