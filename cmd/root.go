package cmd

import (
	"html/template"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

type RunFunc func(cmd *cobra.Command, args []string)

type Report struct {
	TotalTime          time.Duration
	NumberOfRequests   int32
	NumberOfRequestsOk int32
	Requests           map[int32]int32
}

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
	rootCmd.Flags().Int32P("requests", "r", 0, "Number of requests")
	rootCmd.Flags().Int32P("concurrency", "c", 0, "Number of simultaneous calls")
	rootCmd.MarkFlagsRequiredTogether("url", "requests", "concurrency")
}

func executeStressTest() RunFunc {
	return func(cmd *cobra.Command, args []string) {
		startExecutionTime := time.Now()

		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt32("requests")
		concurrency, _ := cmd.Flags().GetInt32("concurrency")

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(int(requests))

		report := Report{
			Requests: make(map[int32]int32),
		}

		for i := 0; i < int(concurrency); i++ {
			go func() {
				mutex := sync.Mutex{}

				for i := 0; ; i++ {
					resp, _ := http.Get(url)

					defer resp.Body.Close()

					atomic.AddInt32(&report.NumberOfRequests, 1)

					if resp.StatusCode == http.StatusOK {
						atomic.AddInt32(&report.NumberOfRequestsOk, 1)
					} else {
						mutex.Lock()
						report.Requests[int32(resp.StatusCode)]++
						mutex.Unlock()
					}

					waitGroup.Done()
				}
			}()
		}

		waitGroup.Wait()

		report.TotalTime = time.Since(startExecutionTime).Round(time.Millisecond)

		generateReport(report)
	}
}

func generateReport(report Report) {
	temp := template.Must(template.New("template.txt").ParseFiles("template.txt"))

	err := temp.Execute(os.Stdout, report)
	if err != nil {
		panic("Generating report error: " + err.Error())
	}
}
