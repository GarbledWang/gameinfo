package main

import (
    "fmt"
    "strconv"
    "github.com/PuerkitoBio/goquery"
    _ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type GameInfo struct{
    name string
    gametype string
    author string
    pubtime string
    version string
    recommend string
    desc string
}

var page = 1

func main(){
    var db,err = sql.Open("sqlite3","./game")
  	if err!= nil{
		fmt.Println(err)
		return
	}
    defer db.Close()
    for{
        if parse(db,page) {
            page++
        }else{
            break
        }
    }
}

func parse(db *sql.DB,page int) bool{
    doc,err := goquery.NewDocument("http://games.tgbus.com/default.aspx?page="+strconv.Itoa(page)+"&elite=0&keyword=&type=0&tag=0&nid=40");
    if err!=nil {
        fmt.Println(err)
        return false
    }
    node := doc.Find(".ml")
    if node.Length() > 0 {
        info := new(GameInfo)
        doc.Find(".ml").Each(func (i int, s *goquery.Selection)  {
            title := s.Find("a").First().Text()
            info.name = title
            fmt.Println(title)
            s.Find(".ml-c2 dd").Each(func(i int, is *goquery.Selection){
                content := is.Text()
                switch i {
                    case 0:
                        info.gametype = content
                    case 1:
                        info.author = content
                    case 2:
                        info.pubtime = content
                    case 3:
                        info.version = content
                    case 4:
                        info.recommend = content
                    case 5:
                        info.desc = content
                }
            })
            stmt,err := db.Prepare("INSERT INTO gameinfo(name,gametype,author,pubtime,version,recommend,desc) values(?,?,?,?,?,?,?)")
            if err!=nil {
            	fmt.Println(err)
			}
            stmt.Exec(info.name,info.gametype,info.author,info.pubtime,info.version,info.recommend,info.desc)
            fmt.Println(info)
        })
        return true   
    }else{
        return false
    }
}
/*

CREATE TABLE 'gameinfo'(
    name VARCHAR(100) NULL,
    gametype VARCHAR(64) NULL,
    author VARCHAR(64) NULL,
    pubtime VARCHAR(64) NULL,
    version VARCHAR(64) NULL,
    recommend VARCHAR(64) NULL,
    desc VARCHAR(64) NULL
)*/
