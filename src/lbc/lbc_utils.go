package main

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func curl(url string, c chan []byte) {

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body_bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	c <- body_bytes
}

func fix_broken_html(page []byte) string {
	page_reader := bytes.NewReader(page)
	root, err := html.Parse(page_reader)

	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	html.Render(&b, root)
	fixedHtml := b.String()
	return fixedHtml
}

func decode_utf8(fixedHtml string) string {
	e := charmap.ISO8859_15
	reader := strings.NewReader(fixedHtml)
	rInUTF8 := transform.NewReader(reader, e.NewDecoder())

	return reader_to_str(rInUTF8)
}

func remove_noscript(page string) []byte {
	str_noscript := strings.Replace(page, "<noscript>", "", -1)
	str_noscript = strings.Replace(str_noscript, "<noscript style=\"color:red\">", "", -1)
	str_noscript = strings.Replace(str_noscript, "</noscript>", "", -1)

	runes := []byte(str_noscript)
	return runes
}

func get_only_url(str string) string {
	var Url *url.URL
	Url, err := url.Parse(str)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	parameters := url.Values{}
	Url.RawQuery = parameters.Encode()
	return Url.String()
}

func reader_to_str(readr io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(readr);
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return buf.String()
}
