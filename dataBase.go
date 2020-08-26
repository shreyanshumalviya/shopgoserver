package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	username := "shreyanshumalviya"
	password := "homnhomnhomn"

	db, err := sql.Open("mysql", username+":"+password+"@tcp(127.0.0.1:3306)/shreyanshumalviya")
	if err != nil {
		log.Print(err.Error())
	}
	table, err := db.Query("CREATE TABLE StockEntry \n( barcode INT, rate FLOAT, price FLOAT,batchNo TEXT, mfd TEXT,expDate TEXT)")

	if err != nil {
		fmt.Println("ad")
		log.Print(err.Error())
	}
	defer db.Close()
	defer table.Close()
}
