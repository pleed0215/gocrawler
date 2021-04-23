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
	getPages()
}

func getJobUrl(job string, page int) (string, error) {
	if(page > 0) {
		return fmt.Sprintf("%s?q=%s&start=%d", baseUrl, job, (page-1)*10), nil
	} else {
		return "", errors.New("page must greater than 0")
	}
	
}

func getPages() int {
	url, _ := getJobUrl("python", 1)
	res, err := http.Get(url)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find(".pagination").Each( func(i int, s *goquery.Selection) {
		s.Find("a").Each( func(i int, s *goquery.Selection) {
			fmt.Println(s.Text())
		})
	})

	return 0
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("fail to connect. ", "Code: ", res.StatusCode, "Status: ", res.Status)
	}
}