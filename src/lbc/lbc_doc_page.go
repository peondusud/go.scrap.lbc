package main

import (
	"fmt"
	"launchpad.net/xmlpath"
	"log"
	"sync"
	//"text/template"
	"time"
)

type Lbc_doc struct {
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"cdate"`
}

type Lbc_doc_struct struct {
	c_doc_urls  chan string
	c_documents chan *Lbc_doc
	wg          sync.WaitGroup
}

func doc_process(c_doc_urls chan string, c_documents chan Lbc_doc) {
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		go doc_worker(c_doc_urls, c_documents, &wg)

	}
    wg.Wait()
	log.Println("Doc Process Done")
}

func doc_worker(c_doc_urls chan string, c_documents chan Lbc_doc, wg *sync.WaitGroup) {
	defer wg.Done()
	doc_url := <-c_doc_urls
	log.Println("+--+ DOC URL:", doc_url)
	c_doc_page := make(chan []byte)
	curl(doc_url, c_doc_page)
	doc_parse(c_doc_page, c_documents, wg)
}

func doc_parse(c_doc_page chan []byte, c_documents chan Lbc_doc, wg *sync.WaitGroup) {
	doc_page := <-c_doc_page
	fmt.Printf("%s\n", string(doc_page))

	title_xpath := xmlpath.MustCompile("'/html/body/div/div[2]/div/div[3]/div/div[1]/div[1]/h1/text()'") //doc urls

	doc_page_noscript := remove_noscript(doc_page)
	fix_html := fix_broken_html(doc_page_noscript)
	utf8_reader := decode_utf8(fix_html)
	root, err := xmlpath.ParseHTML(utf8_reader)

	if err != nil {
		log.Fatal(err)
	}

	if doc_title, ok := title_xpath.String(root); ok {
		log.Println("DOC Title:", doc_title)
	}

}
