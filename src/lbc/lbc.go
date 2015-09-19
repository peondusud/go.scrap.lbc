package main

import (
	"log"
	"runtime"
	//"sync"
)

func main() {

	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	list_urls := []string{
		"http://www.leboncoin.fr/vins_gastronomie/offres/ile_de_france/"}

	c_front_urls := make(chan string, 1000)
	c_doc_urls := make(chan string, 1000*35)
	c_documents := make(chan Lbc_doc, 1000*35)

	go front_process(list_urls, c_front_urls, c_doc_urls)
	doc_process(c_doc_urls, c_documents)

	log.Println("All Done")
}
