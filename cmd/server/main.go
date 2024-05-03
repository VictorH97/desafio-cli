package main

import (
	"html/template"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type RunFunc func(cmd *cobra.Command, args []string)

type Report struct {
	TotalTime          time.Duration
	NumberOfRequests   int
	NumberOfRequestsOk int
	Requests           map[int]int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Stress test",
	Short: "Perform a stress test on a website url",
	Long:  ``,
	Run:   executeStressTest(),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Url of the website")
	rootCmd.Flags().IntP("requests", "r", 0, "Number of requests")
	rootCmd.Flags().IntP("concurrency", "c", 0, "Number of simultaneous calls")
	rootCmd.MarkFlagsRequiredTogether("url", "requests", "concurrency")
}

func executeStressTest() RunFunc {
	return func(cmd *cobra.Command, args []string) {
		startExecutionTime := time.Now()

		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(requests)

		report := Report{
			Requests: make(map[int]int),
		}

		for i := 0; i < int(math.Ceil(float64(requests)/float64(concurrency))); i++ {
			for j := 0; j < int(math.Abs(float64(requests*i-concurrency))); j++ {
				go func() {
					resp, err := http.Get(url)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()

					report.NumberOfRequests++

					if resp.StatusCode == http.StatusOK {
						report.NumberOfRequestsOk++
					} else {
						report.Requests[resp.StatusCode]++
					}

					waitGroup.Done()
				}()
			}

			time.Sleep(time.Second)
		}

		waitGroup.Wait()

		report.TotalTime = time.Since(startExecutionTime)

		generateReport(report)
	}
}

func generateReport(report Report) {
	temp := template.Must(template.New("template.html").ParseFiles("cmd/server/template.html"))

	err := temp.Execute(os.Stdout, report)
	if err != nil {
		panic("Generating report error: " + err.Error())
	}
}
