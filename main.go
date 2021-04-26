package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseUrl string = "https://kr.indeed.com/jobs"
/* 
	indeed.com
	Getting max page strategy
	- kr.indeed.com/jobs?q=job&start=HUGE_NUMBER 
	이렇게 query를 하면 제일 마지막 페이지를 보여주기 때문에 max page number를 얻을 수 있다.
*/ 
const HUGE_NUMBER = 1000000
const PAGE_SIZE = 30


// struct for job information. 
type Job struct {
	title 		string
	company 	string
	summary		string
	location	string
}

func main() {
	jobUrl, err := getJobUrl("python", 2)
	jobChan := make(chan [] Job)
	jobname := "python"

	if(err == nil) {
		fmt.Println(jobUrl)
	} else {
		fmt.Println(err)
	}
	fmt.Print("Fetching total pages for job: ", jobname, "... ")
	maxPage := getMaxPage(jobname)
	fmt.Println("done.")
	fmt.Println("Total: ", maxPage, " pages")

	fmt.Println("Start fetching jobs....using go channels.")
	allJobs := make([]Job, 0, PAGE_SIZE*maxPage)
	startTime := time.Now()
	for i:=0 ; i < maxPage ; i++ {
		go getJobs(jobname, i+1, jobChan)
	}
	for i:=0 ; i < maxPage ; i++ {
		job := <- jobChan
		allJobs = append(allJobs, job...)
	}
	fmt.Println("Fetching done. Now saving to file: ", jobname+".csv")
	saveToCsv(allJobs, jobname)
	endTime := time.Now()
	fmt.Println("Saving done. And Total fetching time is", (endTime.Sub(startTime)))
}



func getJobUrl(job string, page int) (string, error) {
	if(page > 0) {
		return fmt.Sprintf("%s?q=%s&start=%d&limit=%d", baseUrl, job, (page-1)*10, PAGE_SIZE), nil
	} else {
		return "", errors.New("page must greater than 0")
	}	
}


func getMaxPage(job string) int {
	url, _ := getJobUrl(job, HUGE_NUMBER)
	res, err := http.Get(url)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	maxPage := 0
	doc.Find(".pagination").Find("a").Each( func(i int, s *goquery.Selection) {
		val, exist := s.Attr("aria-label")
		if exist {
			page, error := strconv.Atoi(val)
			if error == nil {
				if page >= maxPage {
					maxPage = page
				}
			}
		}
	})

	return maxPage
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

func moreTrimSpace(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getJobs(job string, page int, c chan [] Job)  {
	url, _ := getJobUrl(job, page)
	fmt.Println("Job Scrapping [", job, "] ", "Page: ", page)
	var jobs = make([] Job, 0, PAGE_SIZE);

	res, err := http.Get(url)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find(".jobsearch-SerpJobCard").Each( func(i int, s *goquery.Selection) {
		title := moreTrimSpace(s.Find(".title").Text())
		company := moreTrimSpace(s.Find(".company").Text())
		summary := moreTrimSpace(s.Find(".summary").Text())
		location := moreTrimSpace(s.Find(".location").Text())
		jobs = append(jobs, Job{title, company, summary, location})
	})
	c<-jobs
}

func saveToCsv(jobs [] Job, filename string) error {
	lenFilename := len(filename)
	var toSaveFilename string

	// .csv 확장자를 가지고 있는지 확인, 없으면 무조건 마지막에 .csv 붙이기
	if lenFilename < 5 {
		toSaveFilename = filename + ".csv"
	} else {
		lastFour := filename[lenFilename-4:lenFilename]
		if lastFour != ".csv" {
			toSaveFilename = filename + ".csv"
		} else {
			toSaveFilename = filename
		}
	}

	file, err := os.Create("./"+toSaveFilename)
	if err != nil {
		panic(err) // 파일 못 만들어? 그럼 종료.
	}
	writer := csv.NewWriter(bufio.NewWriter(file))
	defer writer.Flush()

	writer.Write([] string{"Title", "Location", "Company", "Summary"})
	
	for _, job := range jobs {
		err := writer.Write([] string{job.title, job.location, job.company, job.summary})
		checkError(err)
	}

	

	return nil
}