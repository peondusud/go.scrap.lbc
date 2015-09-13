package main

import (
    //"fmt"
    "io/ioutil"
    "io"
    "os"
    "log"
    "sync"
    "strings"
    "net/http"
    //"net/url"
    "bytes"
    "golang.org/x/net/html"
    "golang.org/x/text/transform"
    "golang.org/x/text/encoding/charmap"
    "launchpad.net/xmlpath"

)


func curl( url string, c chan []byte) {

    response, err := http.Get( url )
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

func fix_broken_html( page []byte) string {
    page_reader :=  bytes.NewReader(page)
    root, err := html.Parse(page_reader)

    if err != nil {
        log.Fatal(err)
    }

    var b bytes.Buffer
    html.Render(&b, root)
    fixedHtml := b.String()
    return fixedHtml

}

func decode_utf8(fixedHtml string) io.Reader{
    e := charmap.ISO8859_15
    reader := strings.NewReader(fixedHtml)
    rInUTF8 := transform.NewReader(reader, e.NewDecoder())
    return rInUTF8
}

func remove_noscript(page []byte ) []byte{
    str_noscript := string(page)
    str_noscript = strings.Replace( str_noscript, "<noscript>", "",-1)
    str_noscript = strings.Replace( str_noscript, "</noscript>", "",-1)
    runes := []byte(str_noscript)
    return runes
}

func front_parse( c_front_urls chan string, c_front_page chan []byte, c_doc_urls chan string, wg *sync.WaitGroup ){
    front_page := <-  c_front_page
    //fmt.Printf("%s\n", string(page))
    //path := xmlpath.MustCompile("/html/body/div/div[2]/div/div[3]/div/div[1]/div[1]/h1/text()") //title
    doc_urls_xpath := xmlpath.MustCompile( "/html/body/div[@id=\"page_align\"]/div[@id=\"page_width\"]/div[@id=\"ContainerMain\"]/div[@class=\"content-border list\"]/div[@class=\"content-color\"]/div[@class=\"list-lbc\"]//a/@href") //doc urls
    next_front_urls_xpath := xmlpath.MustCompile( "/html/body/div[@id=\"page_align\"]/div[@id=\"page_width\"]/div[@id=\"ContainerMain\"]/nav/ul[@id=\"paging\"]/li[@class=\"page\"]" ) //next url

    //next_front_urls_xpath := xmlpath.MustCompile( "/html/body/div[1]/div[2]/div/nav[4]/ul/li[24]/a" )
    front_page_noscript := remove_noscript( front_page )
    fix_html := fix_broken_html( front_page_noscript )
    utf8_reader := decode_utf8( fix_html )
    root, err := xmlpath.ParseHTML( utf8_reader )

    if err != nil {
        log.Fatal(err)
    }

    doc_urls := doc_urls_xpath.Iter(root);
    for doc_urls.Next()  {
        doc_url := doc_urls.Node().String()
        c_doc_urls <- doc_url
        //log.Println( "Doc URL:", doc_url)
    }


    prev_next_front_urls := next_front_urls_xpath.Iter(root);
    var node *xmlpath.Node

    for prev_next_front_urls.Next()  {
        node = prev_next_front_urls.Node()
    }

    href_xpath := xmlpath.MustCompile( "a/@href" )
    if next_front_url, ok := href_xpath.String( node ); ok {
        c_front_urls <- next_front_url
        log.Println("Next Front URL:", next_front_url)
        wg.Add(1)
        go front_worker( c_front_urls , c_front_page, c_doc_urls, wg)

        } else {
            log.Println("No Next Front URL")
            log.Println("Front DONE")
            return
        }


    }

    func front_worker( c_front_urls chan string, c_front_page chan []byte, c_doc_urls chan string, wg *sync.WaitGroup){
        defer wg.Done() //call at front_worker func exit
        front_url := <- c_front_urls
        curl( front_url, c_front_page )
        front_parse( c_front_urls, c_front_page , c_doc_urls, wg)
    }

    func front_process( list_urls []string, c_front_urls chan string, c_doc_urls chan string) {
        var wg sync.WaitGroup
        c_front_page := make( chan []byte, 1000 )
        for _,url := range list_urls {
            c_front_urls <- url
            wg.Add(1)
            go front_worker( c_front_urls, c_front_page, c_doc_urls, &wg)
        }
        wg.Wait()
        log.Println("Front Process Done")
    }

    func doc_worker(c_doc_urls chan string, c_documents chan []byte, wg *sync.WaitGroup){

    }

    func doc_process( c_doc_urls chan string, c_documents chan []byte) {
        var wg sync.WaitGroup
        go doc_worker( c_doc_urls, c_documents, &wg)
        log.Println("Doc Process Done")
    }

    func main() {

        list_urls := []string{ "http://www.leboncoin.fr/vins_gastronomie/offres/ile_de_france/",  }
        c_front_urls := make( chan string, 1000)
        c_doc_urls := make( chan string, 100000)

        front_process(list_urls, c_front_urls, c_doc_urls)
        //c_documents := make( chan []byte )
        //doc_process(c_doc_urls, c_documents)

        log.Println("All Done")
    }
