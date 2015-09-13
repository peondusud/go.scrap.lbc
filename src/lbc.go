package main

import (
    "fmt"
    "io/ioutil"
    "io"
    "os"
    "log"
    "strings"
    "net/http"
    "golang.org/x/net/html"
    "sync"
    //"golang.org/x/text/encoding"
    "golang.org/x/text/transform"
    "golang.org/x/text/encoding/charmap"
    "launchpad.net/xmlpath"
    "bytes"
    //"net/url"
)


func curl( url string , c chan []byte)  {


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

func fix_broken_html( page []byte) string{
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

func parse( c chan []byte ){
    page := <- c
    //fmt.Printf("%s\n", string(page))
    //path := xmlpath.MustCompile("/html/body/div/div[2]/div/div[3]/div/div[1]/div[1]/h1/text()") //title
    path := xmlpath.MustCompile("/html/body/div[@id=\"page_align\"]/div[@id=\"page_width\"]/div[@id=\"ContainerMain\"]/div[@class=\"content-border list\"]/div[@class=\"content-color\"]/div[@class=\"list-lbc\"]//a/@href") //doc urls
    //path := xmlpath.MustCompile('/html/body/div[@id="page_align"]/div[@id="page_width"]/div[@id="ContainerMain"]/nav/ul[@id="paging"]//li[@class="page"]/a/@href' //next url

    str_noscript := string(page)
    str_noscript = strings.Replace( str_noscript, "<noscript>", "",-1)
    str_noscript = strings.Replace( str_noscript, "</noscript>", "",-1)

    runes := []byte(str_noscript)
    fix_html := fix_broken_html( runes )


    utf8_reader := decode_utf8(fix_html)

    root, err := xmlpath.ParseHTML(utf8_reader) //FIXME UTF-8
    if err != nil {
            log.Fatal(err)
    }

    //if value, ok := path.String(root); ok { fmt.Println("title:", value)  }
    doc_urls := path.Iter(root);
     for  doc_urls.Next()  {
            //fmt.Println("doc urls:", doc_urls.Node().String())
    }
}

func front_worker (url string, c chan []byte, wg *sync.WaitGroup){

    go curl( url , c)
    //fmt.Printf("%s\n", string(body_bytes))
    parse(c)
    defer wg.Done()

}


func main() {
    var wg sync.WaitGroup
    c := make(chan []byte)
    list := []string{   "http://www.leboncoin.fr/vins_gastronomie/offres/ile_de_france/",  }
    for _,url := range list {
        wg.Add(1)
        go front_worker(url , c, &wg)
    }
    wg.Wait()
    fmt.Printf("done")
}
