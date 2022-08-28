package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	Url    string
	Resp   float64
	Total  string
	Output int
	Error  string
	Body   string
}

func main() {
	host := ""
	path := ""
	fmt.Println("Введите host: ")
	fmt.Scan(&host)
	fmt.Println("Введите path: ")
	fmt.Scan(&path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Data := make([]byte, 70000)
	for {
		_, err := file.Read([]byte(Data))
		if err != io.EOF {
			break
		}
	}
	DataStings := strings.Split(string(Data), "\n")
	start := time.Now()
	Gorun(DataStings, host)
	secs := time.Since(start).Seconds()
	fmt.Println(secs)
}
func Gorun(Sites []string, host string) {
	Ch := make(chan Page, 3)
	file, err := os.OpenFile("DataBase.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(Sites); i++ {
		go func(u string, host string) {
			site := host + u
			ping(site, Ch)
			res, err := json.Marshal(<-Ch)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.WriteString(file, string(res)+"\n")
			if err != nil {
				log.Fatal(file, err)
			}
		}(Sites[i], host)
		time.Sleep(125000)
	}
}
func ping(url string, Ch chan Page) {
	start := time.Now()
	respons, err := http.Get(url)
	var code int
	var errror string
	if err != nil {
		errror = fmt.Sprintf("%s\n", err)
		return
	}
	defer respons.Body.Close()
	code = respons.StatusCode
	codetext, err := ioutil.ReadAll(respons.Body)
	if err != nil {
		log.Fatal(err)
	}
	pass := true
	if err != nil {
		pass = false
	}
	secs := time.Since(start).Seconds()
	if code == 404 {

	} else {
		if pass {
			secs, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", secs), 2)
			Ch <- Page{Resp: secs, Url: url, Output: code, Total: "Pass", Error: errror, Body: string(codetext)}
		}
	}
}
