package main

import (
	"fmt"
	"launchpad.net/xmlpath"
	"log"
	"sync"
	"strings"
	//"text/template"
	"time"
	//"regexp"
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
		//for i := range c_doc_urls {
		for i := 1; i <= 100; i++ {
			//log.Println("RANGE value=%s", i)
			wg.Add(1)
			go doc_worker(c_doc_urls, c_documents, &wg)
		}
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
	log.Println("Doc Process Done")
}

func doc_worker(c_doc_urls chan string, c_documents chan Lbc_doc, wg *sync.WaitGroup) {
	defer wg.Done()
	doc_url := <-c_doc_urls
	//log.Println("+--+ DOC URL:", doc_url)
	c_doc_page := make(chan []byte, 200*35)
	curl(doc_url, c_doc_page)
	doc_parse(c_doc_page, c_documents, wg)
}

func doc_parse(c_doc_page chan []byte, c_documents chan Lbc_doc, wg *sync.WaitGroup) {
	doc_page := <-c_doc_page
	//fmt.Printf("%s\n", string(doc_page))

	utf8_reader := decode_utf8( string(doc_page) )
	doc_page_noscript := remove_noscript( utf8_reader )
	fix_html := fix_broken_html(doc_page_noscript)

	//r, _ := regexp.Compile("<meta .+?=.+?=\"\".+>")
	//r, _ := regexp.Compile("<meta .+?")
	//str_noscript := r.ReplaceAllString(fix_html, "")

	//root, err := xmlpath.ParseHTML( strings.NewReader(fix_html) )
	_, err := xmlpath.ParseHTML( strings.NewReader(fix_html) )

	if err != nil {
		fmt.Println( "!!!!!!!!!!!! BUG DOC page", err)
		//fmt.Println( "!!!!!!!!!!!! BUG DOC page utf8", string(doc_page_noscript))
		log.Println( "!!!!!!!!!!!! BUG DOC page", err)
		return
	}
	/*
	title_xpath := xmlpath.MustCompile("/html/body/div/div[2]/div/div[3]/div/div[1]/div[1]/h1/text()") //doc urls

	if doc_title, ok := title_xpath.String(root); ok {
		log.Println("##### DOC Title:", doc_title)
	}
	*/

}
