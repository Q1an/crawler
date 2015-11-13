// file: crawler.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type article struct {
	id          int64
	title       string
	description string
	body        string
	author      string
	date        string
}

var nodes = make([]article, 0)

func Allcollectdata(httpBody io.Reader) (string, string, string, string, string) {
	title := ""
	body := ""
	date := ""
	author := ""
	description := ""
	page := html.NewTokenizer(httpBody)
	titlecheck := false
	descriptioncheck := false
	datacheck := false
	authorcheck := false
	divclasscontentcheck := false
	divcount := 0
	indexcheck := false
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return title, description, date, body, author
		}
		token := page.Token()
		if indexcheck {
			if token.Data == "cnBeta.COM_中文业界资讯站" {
				return "", "", "", "", ""
			}
			indexcheck = false
		} else if titlecheck {
			title = token.Data
			titlecheck = false
		} else if datacheck {
			date = token.Data
			datacheck = false
		} else if authorcheck {
			author = token.Data[16 : len(token.Data)-1]
			authorcheck = false
		} else if divclasscontentcheck {
			if tokenType == html.TextToken {
				body = body + token.Data
			} else if tokenType == html.EndTagToken && token.DataAtom.String() == "div" {
				if divcount == 0 {
					divclasscontentcheck = false
				} else {
					divcount--
				}
			} else if tokenType == html.StartTagToken && token.DataAtom.String() == "div" {
				divcount++
			}
		} else if tokenType == html.StartTagToken && token.DataAtom.String() == "div" && len(token.Attr) > 0 {
			if token.Attr[0].Val == "content" {
				divclasscontentcheck = true
			}
		} else if tokenType == html.StartTagToken && token.DataAtom.String() == "h2" {
			titlecheck = true
		} else if tokenType == html.StartTagToken && token.DataAtom.String() == "title" {
			indexcheck = true
		} else if tokenType == html.StartTagToken && token.DataAtom.String() == "span" && len(token.Attr) > 0 {
			if token.Attr[0].Val == "date" {
				datacheck = true
			} else if token.Attr[0].Val == "author" {
				authorcheck = true
			}
		} else if tokenType == html.SelfClosingTagToken && token.DataAtom.String() == "meta" {
			for _, attr := range token.Attr {
				if attr.Key == "name" {
					if attr.Val == "description" {
						descriptioncheck = true
					}
				} else if attr.Key == "content" {
					if descriptioncheck == true {
						description = attr.Val
						descriptioncheck = false
					}
				}
			}
		}
	}
}

func getdata(link string) (string, string, string, string, string) {
	resp, err := http.Get("http://www.cnbeta.com/articles/" + link + ".htm")
	if err != nil {
		return "", "", "", "", ""
	}
	defer resp.Body.Close()
	title, description, date, body, author := Allcollectdata(resp.Body)
	return title, description, date, body, author
}

func watching() {
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
}

func main() {

	db, err := sql.Open("sqlite3", "./cnbeta222.db")
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO article(id, title, description, body, author, date) values(?,?,?,?,?,?)")
	checkErr(err)

	t := time.Tick(time.Second)
	startvalue := 445511
	endvalue := 1

	go func() {
		for {
			select {
			case <-t:
				watching()
			}
		}
	}()

	nbConcurrentGet := 222
	urls := make(chan string, nbConcurrentGet)
	var wg sync.WaitGroup
	for i := 0; i < nbConcurrentGet; i++ {
		go func() {
			for url := range urls {
				title, description, date, body, author := getdata(url)
				if title != "" {
					urlid, _ := strconv.ParseInt(url, 10, 64)
					//if urlid%1000 == 0 {
					//	fmt.Println("I am now in %d.", urlid)
					//}
					fmt.Println("I am now in ", urlid)
					//defer stmt.Exec(urlid, title, description, body, author, date)
					var temp article
					temp.id = urlid
					temp.title = title
					temp.description = description
					temp.body = body
					temp.author = author
					temp.date = date
					nodes = append(nodes, temp)
					//res = res
				}
				wg.Done()
			}
		}()
	}
	for i := startvalue; i != endvalue; i-- {
		wg.Add(1)
		urls <- fmt.Sprintf("%d", i)
	}
	wg.Wait()
	fmt.Println("Finished grabing")

	for index, node := range nodes {
		fmt.Println(index)
		_, err := stmt.Exec(node.id, node.title, node.description, node.body, node.author, node.date)
		checkErr(err)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
		fmt.Println(err)
		return
	}
}
