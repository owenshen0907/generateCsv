// generateCsv project doc.go

/*
generateCsv document
*/
package main

import (
	"database/sql"

	"flag"
	"fmt"
	"io/ioutil"
	"os"
	//	"strconv"
	//"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/larspensjo/config"
)

var (
	dbuser   string
	dbpsw    string
	dbhost   string
	dbport   string
	dbschema string
)

var db *sql.DB

func init() {
	readLogin()
	//	dbuser = "root"
	//	dbpsw = "Ylch#&%Cdd829"
	//	dbhost = "192.168.10.240"
	//	dbport = "3316"
	//	dbschema = "ylchfl_sit"
	dbuser = TOPIC["dbuser"]
	dbpsw = TOPIC["dbpsw"]
	dbhost = TOPIC["dbhost"]
	dbport = TOPIC["dbport"]
	dbschema = TOPIC["dbschema"]

	dblogin := dbuser + ":" + dbpsw + "@tcp(" + dbhost + ":" + dbport + ")/" + dbschema + "?charset=utf8"
	fmt.Println("加载数据库信息：" + dblogin)
	//db, _ = sql.Open("mysql", "root:owenshen123@tcp(127.0.0.1:3306)/test?charset=utf8")
	db, _ = sql.Open("mysql", dblogin)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

func readInfo(fileName string) string {
	dstfile, err := os.Open(fileName)
	erro(err)
	body, _ := ioutil.ReadAll(dstfile)
	content := string(body)
	defer dstfile.Close()
	return content
}

func readLogin() {
	var (
		configFile = flag.String("configfile", "config.ini", "General configuration file")
	)
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		fmt.Println("read ini error")
		return
	}
	if cfg.HasSection("exe") {
		section, err := cfg.SectionOptions("exe")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("exe", v)
				if err == nil {

					TOPIC[v] = options
				}
			}
		}
	}
}
func operateDB(sqlstring string) {

	defer db.Close()
	rows, err := db.Query(sqlstring)
	erro(err)

	columns, err := rows.Columns()
	DBdata = append(DBdata, columns)
	//fmt.Println(columns)//head
	erro(err)

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		var tmpArr []string
		err = rows.Scan(scanArgs...)
		erro(err)

		// Now do something with the data.
		// Here we just print each column as a string.
		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			//			fmt.Println(string(col))
			if col == nil {
				tmpArr = append(tmpArr, "NULL")
			} else {
				tmpArr = append(tmpArr, string(col))
			}

		}
		DBdata = append(DBdata, tmpArr)
		//		fmt.Println(tmpArr)
	}
	defer rows.Close()
}
func erro(err error) {
	if err != nil {
		fmt.Println("出错了", err)
	}
}
