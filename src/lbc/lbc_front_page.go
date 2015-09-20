package main

import (
	"launchpad.net/xmlpath"
	"log"
	"sync"
	"strings"
	//"fmt"
	//"text/template"
)

func front_process(list_urls []string, c_front_urls chan string, c_doc_urls chan string) {
	var wg sync.WaitGroup
	c_front_page := make(chan []byte, 1000)
	for _, url := range list_urls {
		c_front_urls <- url
		wg.Add(1)
		go front_worker(c_front_urls, c_front_page, c_doc_urls, &wg)
	}
	wg.Wait()
	log.Println("Front Process Done")
}

func front_worker(c_front_urls chan string, c_front_page chan []byte, c_doc_urls chan string, wg *sync.WaitGroup) {
	defer wg.Done() //call at front_worker func exit
	front_url := <-c_front_urls
	curl(front_url, c_front_page)
	front_parse(c_front_urls, c_front_page, c_doc_urls, wg)
}

func front_parse(c_front_urls chan string, c_front_page chan []byte, c_doc_urls chan string, wg *sync.WaitGroup) {
	front_page := <-c_front_page
	//fmt.Printf("%s\n", string(front_page))
	//path := xmlpath.MustCompile("/html/body/div/div[2]/div/div[3]/div/div[1]/div[1]/h1/text()") //title
	doc_urls_xpath := xmlpath.MustCompile("/html/body/div[@id=\"page_align\"]/div[@id=\"page_width\"]/div[@id=\"ContainerMain\"]/div[@class=\"content-border list\"]/div[@class=\"content-color\"]/div[@class=\"list-lbc\"]//a/@href") //doc urls
	next_front_urls_xpath := xmlpath.MustCompile("/html/body/div[@id=\"page_align\"]/div[@id=\"page_width\"]/div[@id=\"ContainerMain\"]/nav/ul[@id=\"paging\"]/li[@class=\"page\"]")                                                   //next url
	/*
	front_page_noscript := remove_noscript(front_page)
	fix_html := fix_broken_html(front_page_noscript)
	utf8_reader := decode_utf8(fix_html)
	root, err := xmlpath.ParseHTML(utf8_reader)*/

	utf8_reader := decode_utf8( string(front_page) )
	doc_page_noscript := remove_noscript( utf8_reader )

	fix_html := fix_broken_html(doc_page_noscript)

	//fmt.Println(string(fix_html))

	root, err := xmlpath.ParseHTML( strings.NewReader(fix_html) )

	if err != nil {
		//log.Println("ca rentre")
		log.Fatal("FRONT PAGE",err)
	}

	doc_urls := doc_urls_xpath.Iter(root)
	for doc_urls.Next() {
		doc_url := doc_urls.Node().String()
		c_doc_urls <- doc_url
		//log.Println( "Doc URL:", doc_url)  //<-- DOC URL
	}

	prev_next_front_urls := next_front_urls_xpath.Iter(root)
	var node *xmlpath.Node

	for prev_next_front_urls.Next() {
		node = prev_next_front_urls.Node()
	}

	href_xpath := xmlpath.MustCompile("a/@href")
	if next_front_url, ok := href_xpath.String(node); ok {
		c_front_urls <- next_front_url
		log.Println("Next Front URL:", next_front_url)
		wg.Add(1)
		go front_worker(c_front_urls, c_front_page, c_doc_urls, wg)
	} else {
		log.Println("No Next Front URL")
		log.Println("Front DONE")
		return
	}
}
