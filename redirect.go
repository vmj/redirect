package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main()  {
	addr := flag.String("addr", ":4040", "Listening address")
	temp := flag.Bool("temp", false, "Respond with 302 instead of 301")
	protocol := flag.String("protocol", "#{protocol}",
		"'http' or 'https', or '#{protocol}' to retain the original")
	host := flag.String("host", "#{host}",
		"Domain to which to redirect, or '#{host}' to retain the original")
	port := flag.String("port", "#{port}",
		"Port to which to redirect, or '#{port}' to retain the original")
	path := flag.String("path", "#{path}",
		"The absolute path, starting with /; can include '#{host}', '#{port}', and '#{path}' to retain those from original")
	query := flag.String("query", "#{query}",
		"Query params; can include '#{protocol}', '#{host}', '#{port}', '#{path}', and '#{query}' to retain those from original")

	flag.Parse()

	if *protocol == "#{protocol}" && *host == "#{host}" && *port == "#{port}" && *path == "#{path}" {
		log.Fatal("Must modify at least one of following components to avoid a redirect loop: protocol, host, port, or path.")
	}

	if *protocol != "http" && *protocol != "https" && *protocol != "#{protocol}" {
		log.Fatal("Protocol must be one of: http, https, or #{protocol}")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		orig := url.URL{
			// FIXME: User
			Host: r.Host,
			Path: r.URL.Path,
			RawQuery: r.URL.RawQuery,
			Fragment: r.URL.Fragment,
		}
		if r.TLS != nil {
			orig.Scheme = "https"
		} else {
			orig.Scheme = "http"
		}
		w.Header().Set("Location", makeLocation(*protocol, *host, *port, *path, *query, &orig))
		w.Header().Set("Content-Length", "0")
		if *temp {
			w.WriteHeader(302)
		} else {
			w.WriteHeader(301)
		}
		log.Printf("%s %s %s (%s)\n", r.RemoteAddr, r.Method, r.URL.Path, r.UserAgent())
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func makeLocation(protocol, host, port, path, query string, orig *url.URL) string {
	location := url.URL{}

	if protocol == "#{protocol}" {
		location.Scheme = orig.Scheme // may be absent
	} else {
		location.Scheme = protocol // http or https
	}

	location.User = orig.User

	if host == "#{host}" {
		if port == "#{port}" {
			location.Host = cleanupHost(location.Scheme, orig.Hostname(), orig.Port())
		} else {
			location.Host = cleanupHost(location.Scheme, orig.Hostname(), port)
		}
	} else {
		if port == "#{port}" {
			location.Host = cleanupHost(location.Scheme, host, orig.Port())
		} else {
			location.Host = cleanupHost(location.Scheme, host, port)
		}
	}

	if path == "#{path}" {
		location.Path = cleanupPath(orig.Path)
	} else {
		location.Path = strings.Replace(strings.Replace(strings.Replace(path,
			"#{host}", url.PathEscape(orig.Hostname()), -1),
			"#{port}", url.PathEscape(orig.Port()), -1),
			"#{path}", cleanupPath(orig.Path), -1)
	}

	if query == "#{query}" {
		location.RawQuery = orig.RawQuery
	} else {
		location.RawQuery = strings.Replace(strings.Replace(strings.Replace(strings.Replace(query,
			"#{host}", url.QueryEscape(orig.Hostname()), -1),
			"#{port}", url.QueryEscape(orig.Port()), -1),
			"#{path}", url.QueryEscape(orig.Path), -1),
			"#{query}", orig.RawQuery, -1)
	}

	location.Fragment = orig.Fragment

	return location.String()
}

func cleanupHost(scheme, host, port string) string {
	if port == "" {
		return host
	}
	if scheme == "http" && port == "80" {
		return host
	}
	if scheme == "https" && port == "443" {
		return host
	}
	return fmt.Sprintf("%s:%s", host, port)
}

func cleanupPath(path string) string {
	if path == "" {
		return "/"
	}
	return path
}