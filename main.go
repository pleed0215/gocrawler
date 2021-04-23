package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseUrl string = "https://kr.indeed.com/jobs"

func main() {
	jobUrl, err := getJobUrl("python", 2)
	if(err == nil) {
		fmt.Println(jobUrl)
	} else {
		fmt.Println(err)
	}
}

func getJobUrl(job string, page int) (string, error) {
	if(page > 0) {
		return fmt.Sprintf("%s?q=%s&start=%d", baseUrl, job, (page-1)*10), nil
	} else {
		return "", errors.New("page must greater than 0")
	}
	
}

func getPages() int {
	res, err := http.Get(baseUrl)
	checkError(err)
	checkCode(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc)
	return 0
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(code int) {
	if code != 200 {
		log.Fatalln("fail to connect")
	}
}