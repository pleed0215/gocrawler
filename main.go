package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gj "github.com/pleed0215/gocrawler/get_job"
)

func main() {
	// Echo instance
	e := echo.New()
  
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
  
	// Routes
	e.GET("/", handleHome)
	e.POST("/jobs", handleSearch)
  
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
  }
  
  // Handler
  func handleHome(c echo.Context) error {
	return c.File("home.html")
  }

  func handleSearch(c echo.Context) error {
	cleanedTerm := strings.ToLower(gj.MoreTrimSpace(c.FormValue("job")))
	jobs := gj.GetJobs(cleanedTerm)
	filename := cleanedTerm+".csv"
	gj.JobToCsv(jobs, filename)
	defer os.Remove(filename)

	return c.Attachment(filename, filename)
  }