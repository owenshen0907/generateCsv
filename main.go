// generateCsv project main.go
package main

import (
	//	"io/ioutil"
	"encoding/csv"
	"os"
	"strconv"
	//	"strings"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/qiniu/iconv"
	//"github.com/smartwalle/going/email"
)

var DBdata [][]string
var TOPIC = make(map[string]string)

//var contract []string
var sysdate = time.Now().Format("2006-01")
var minute = time.Now().Format("01021504")
var yyyymmdd = time.Now().Format("20060102")
var tilog = time.Now().Format("2006/01/02 03:04:05 PM")
var sqlString string

func main() {
	os.Mkdir("report/"+sysdate, 0777)
	//generate()
	os.IsExist(os.Mkdir("log", os.ModePerm))
	logFile, _ := os.OpenFile("log/"+sysdate+".txt", os.O_RDWR|os.O_CREATE, 0666)
	SEEK_END, _ := logFile.Seek(0, os.SEEK_END)
	//找到日志的偏移量
	_, _ = logFile.WriteAt([]byte("\r\n"), SEEK_END)
	logFile.WriteString(tilog + "\r\n")
	attachName := generate(logFile)

	defer logFile.Close()
	cd, err := iconv.Open("gbk", "utf-8")
	erro(err)
	defer cd.Close()
	attachName = cd.ConvString(attachName)

	logFile.WriteString("Generate the attachment in/" + attachName + "\r\n")
}
func generate(logFile *os.File) string {
	sqlString = readInfo("sql/test.sql")
	operateDB(sqlString)
	cd, err := iconv.Open("gbk", "utf-8")
	erro(err)
	cd1, err1 := iconv.Open("utf-8", "gbk")
	erro(err1)
	defer cd.Close()
	//32列数据
	logFile.WriteString("  Today have been generated :" + strconv.Itoa(len(DBdata)) + "row data.\r\n")
	attachNamePrefix := cd1.ConvString(TOPIC["attachNamePrefix"])
	//attachNamePrefix = cd.ConvString(attachNamePrefix)
	filename := TOPIC["attachP"] + sysdate + "/" + attachNamePrefix + yyyymmdd + TOPIC["attachNameStffix"]
	fmt.Println("附件名为：", filename)
	//filename = cd.ConvString(filename)
	//fmt.Println("附件名为：", filename)

	//	f, err := os.Create("111.csv")

	for i, _ := range DBdata {

		//	fmt.Println(len(DBdata[0]))
		//fmt.Println(len(DBdata))
		var tmp []string
		fmt.Println(DBdata[i])
		for j := 0; j < len(DBdata[0]); j++ {

			DBdata[i][j] = cd.ConvString(DBdata[i][j])
			if i != 0 {
				DBdata[i][0] = "'" + DBdata[i][0]
			}
			tmp = append(tmp, DBdata[i][j])
			//tmp = tmp + DBdata[i][j]

		}
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0777)
		erro(err)
		defer f.Close()
		w := csv.NewWriter(f)
		//w.Write([]string{tmp})

		w.Write(tmp)
		w.Flush()
		tmp = nil

		//writeToCsv("report/"+sysdate+"/"+TOPIC["attachName"]+yyyymmdd+".csv", DBdata[i][0], DBdata[i][1], DBdata[i][2], DBdata[i][3], DBdata[i][4], DBdata[i][5], DBdata[i][6], DBdata[i][7], DBdata[i][8], DBdata[i][9], DBdata[i][10], DBdata[i][11], DBdata[i][12], DBdata[i][13], DBdata[i][14], DBdata[i][15], DBdata[i][16])
	}
	DBdata = nil
	return filename
}
