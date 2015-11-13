package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "./cnbeta222.db")
	checkErr(err)

	//stmt, err := db.Prepare("INSERT INTO article(id, title, description, body, author, date) values(?,?,?,?,?,?)")
	//checkErr(err)

	//res, err := stmt.Exec(1, "test", "description", "body", "author", "2012-12-09")
	//checkErr(err)
	//res = res

	rows, err := db.Query("SELECT * FROM article")
	checkErr(err)

	for rows.Next() {
		var id int
		var title string
		var description string
		var body string
		var author string
		var created time.Time
		err = rows.Scan(&id, &title, &description, &body, &author, &created)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(title)
		//fmt.Println(description)
		//fmt.Println(body)
		//fmt.Println(author)
		//fmt.Println(created)
	}

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
