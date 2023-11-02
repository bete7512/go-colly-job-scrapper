package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Job struct {
	JobTitle     string `json:"job_title"`
	Company      string `json:"company"`
	Location     string `json:"location"`
	Salary       string `json:"salary"`
	Description  string `json:"description"`
	Requirements string `json:"requirements"`
	Sector       string `json:"sector"`
	Category     string `json:"category"`
	Experience   string `json:"experience"`
	Position     string `json:"position"`
	JobType      string `json:"job_type"`
	TimeLeft     string `json:"time_left"`
}

func main() {

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Free Jobs API",
		})
	})

	jobs, err := getJobsFromScraper()
	if err != nil {
		fmt.Println("Failed to scrape jobs. Fetching from JSON file.")
		jobs, err = getJobsFromJSON()
		if err != nil {
			return
		}
	}
	router.GET("/jobs", func(c *gin.Context) {
		c.JSON(200, jobs)
	})
	router.Run(":8787")

}

func saveJobsAsJSON(jobs []Job) {
	jsonData, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling jobs to JSON:", err)
		return
	}

	err = os.WriteFile("jobs.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("Jobs data saved as 'jobs.json'")
}

func getJobsFromScraper() ([]Job, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("hahu.jobs"),
	)
	collector.OnError(
		func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		},
	)
	var jobs []Job
	collector.OnHTML("div.grid.grid-cols-1.gap-y-10.mt-5", func(e *colly.HTMLElement) {
		e.ForEach("div.w-full.pt-3.px-4.md\\:px-10.xl\\:px-5", func(_ int, el *colly.HTMLElement) {
			job := Job{}
			job.JobTitle = el.ChildText("h3.font-black.text-lg.text-secondary")
			job.Company = el.ChildText("p.text-left.font-normal.text-base.md\\:text-lg.text-secondary.dark\\:text-secondary-4.line-clamp-2")
			job.Sector = el.ChildText("div.flex.items-center.gap-2.dark\\:text-secondary-4")
			job.Location = el.ChildText("div[title=Location]")
			job.Experience = el.ChildText("div[title='Years of Experience']")
			job.Position = el.ChildText("div[title='Number of Positions']")
			job.JobType = el.ChildText("div[title='Job Type']")
			job.Description = el.ChildText("p.mt-3.font-normal.text-sm.md\\:text-lg.leading-6.md\\:leading-9.text-secondary.dark\\:text-secondary-4.description")
			job.TimeLeft = el.ChildText("span.flex.font-body.text-sm.md\\:text-md.leading-9.font-light.capitalize")
			jobs = append(jobs, job)
		})
		fmt.Println("Scraped", len(jobs), "jobs")
		saveJobsAsJSON(jobs)
	})
	collector.Visit("https://hahu.jobs/jobs")
	return jobs, nil
}

func getJobsFromJSON() ([]Job, error) {
	data, err := os.ReadFile("jobs.json")
	if err != nil {
		return nil, err
	}

	var jobs []Job
	err = json.Unmarshal(data, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
