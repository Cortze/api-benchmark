package requester

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cortze/api-benchmark/pkg/config"
	"github.com/cortze/api-benchmark/pkg/utils"
)

var (
	moduleName = "Benchmark-Insides"
	log        = logrus.WithField(
		"module", moduleName,
	)
	timeout = 20
)

type Benchmark struct {
	ctx           context.Context
	BenchmarkName string
	// benchmark metadata
	StartTime     time.Time
	FinishTime    time.Time
	TotalDuration time.Duration

	// Query metadata
	HostEndpoint string
	QueryBase    string // TODO: should have it's own type
	AvgReqTime   time.Duration

	BestScore       time.Duration
	BestScoreQuery  string
	WorstScore      time.Duration
	WorstScoreQuery string

	TotQueries   int
	TimeOut      int
	SuccessReq   int
	FailReq      int
	SuccessRatio float64

	// Extra Code to make it work
	requestResponseChan chan *Request // should have same length as totQueries
	// TODO: Put them in a csv as soon as they arrive?
	requestList []Request // keep in memory all the requests done so far
	queryList   []string
	conf        *config.Config
}

func NewBenchmark(ctx context.Context, conf *config.Config) *Benchmark {
	return &Benchmark{
		ctx:                 ctx,
		BenchmarkName:       conf.BenchmarkName,
		HostEndpoint:        conf.HostEndpoint,
		QueryBase:           conf.Query,
		requestResponseChan: make(chan *Request, conf.NumQueries),
		requestList:         make([]Request, 0),
		queryList:           make([]string, 0),
		conf:                conf,
	}
}

// Benchmark Methods

func (b *Benchmark) ImportQueryListFromFile(path string) error {
	qfile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer qfile.Close()

	scanner := bufio.NewScanner(qfile)

	for scanner.Scan() {
		if scanner.Text() != "\n" {
			b.queryList = append(b.queryList, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (b *Benchmark) ComposeQueryList() error {
	for i := 0; i < b.conf.NumQueries*b.conf.ConcurrentReq; i++ {
		// Compose the ranges later used for composing the queries
		ranges := make([]*utils.Range, 0)
		for _, rang := range b.conf.RangeValues {
			r, err := utils.NewRangeFromString(rang)
			if err != nil {
				return err
			}
			ranges = append(ranges, r)
		}

		// compose the queries and add them to the list
		query := b.QueryBase
		for idx, rep := range b.conf.Replaces {
			query = strings.Replace(query, rep, ranges[idx].GetRandomNumberStr(), -1)
		}
		b.queryList = append(b.queryList, query)
	}
	return nil
}

func (b *Benchmark) ExportQueryList(path string) error {
	// open file
	qfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer qfile.Close()
	// write each of the queries in a line
	for _, query := range b.queryList {
		_, err := qfile.WriteString(query + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Benchmark) Run() {

	ctx, cancel := context.WithCancel(b.ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	// control variables to make sure we don't end up the benchmarkbefore the requester or the consumer are done
	var wg sync.WaitGroup

	// request consumer loop
	wg.Add(1)
	go func() {
		log.Info("consumer Go Routine initialized")
		// set up the control counter to don't keep the routine up forever
		ctrlCnt := 0

		// averate and total counters
		totDuration := time.Duration(0)

		for {
			select {
			case req := <-b.requestResponseChan:
				// check if the req is not empty
				if req == nil {
					log.Error("Empty req has been received on the consumer Go routine")
					continue
				}

				// counters
				b.TotQueries++

				if req.Status == "200" {
					b.SuccessReq++
				} else {
					b.FailReq++
				}

				// analyze the time
				// empty base?
				if b.BestScore == time.Duration(0) {
					b.BestScore = req.ResponseTime
					b.BestScoreQuery = req.Query
				}
				if b.WorstScore == time.Duration(0) {
					b.WorstScore = req.ResponseTime
					b.WorstScoreQuery = req.Query
				}
				// get best and worst times
				if req.ResponseTime > b.WorstScore {
					b.WorstScore = req.ResponseTime
					b.WorstScoreQuery = req.Query
				}
				if req.ResponseTime < b.BestScore {
					b.BestScore = req.ResponseTime
					b.BestScoreQuery = req.Query
				}
				// add req time to the total
				totDuration += req.ResponseTime

				// calculate the aggregations
				b.AvgReqTime = totDuration / time.Duration(b.TotQueries)
				b.SuccessRatio = (float64(b.SuccessReq) * 100) / float64(b.TotQueries)

				// add the ReqStatus vlaue to the list
				b.requestList = append(b.requestList, *req)

				// increase the control counter to finish the consumer
				ctrlCnt++
				if ctrlCnt >= b.conf.NumQueries*b.conf.ConcurrentReq {
					log.Infof("total number of queries has reached it max: %d/%d", ctrlCnt, b.conf.NumQueries)
					log.Info("closing consumer routine")
					wg.Done()
					return
				}

			case <-ctx.Done():
				log.Warn("context has been canceled, closing consumer routine")
				wg.Done()
				return
			}
		}

	}()

	// track the initial time of the benchmark
	tinit := time.Now()
	b.StartTime = tinit

	// main request loop
	wg.Add(1)
	go func() {
		log.Info("requester Go Routine initialized")

		// generate a http.Client to manage timeouts
		httpCli := http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
		// check if there are enough queries as the total of

		// requester main loop
		i := 0 // for the number of requests (of X number of simultaneous )
		q := 0 // for the number queries that has been readed (also index)
		for i = 0; i < b.conf.NumQueries; i++ {
			// check if the ctx has died
			if ctx.Err() != nil {
				log.Error("context has died")
				break
			}

			// check if there are enough queries as te ones requested in the config file
			// repeat from the begining otherwise
			var queries []string
			for j := 0; j < b.conf.ConcurrentReq; j++ {
				if q < len(b.queryList) {
					// repeat the queries from the top
					// compose the query to request form the base and from the random numbers
					queries = append(queries, b.HostEndpoint+b.queryList[q])
				} else {
					// compose the query to request form the base and from the random numbers
					queries = append(queries, b.HostEndpoint+b.queryList[q-len(b.queryList)])
				}
			}

			log.Debug("nex queries to request \n", queries)
			// launch the concurrent calls
			var wgr sync.WaitGroup
			for j := 0; j < b.conf.ConcurrentReq; j++ {
				wgr.Add(1)
				queryidx := j
				go func() {
					// request the query
					tnow := time.Now()
					resp, err := httpCli.Get(queries[queryidx])
					log.Debug(queryidx, "response to query:", queries[queryidx], "\n", resp)
					reqTime := time.Since(tnow)

					var reqStatus *Request
					//
					if resp != nil {
						// compose the Request obj
						reqStatus = NewRequest(queries[queryidx], fmt.Sprintf("%d", resp.StatusCode), tnow, reqTime, "NONE")
					} else {
						// compose the Request obj
						reqStatus = NewRequest(queries[queryidx], fmt.Sprintf("NONE"), tnow, reqTime, "NONE")
					}
					if err != nil {
						reqStatus.Error = err.Error()
					}
					log.Debug(queryidx, "tracked status: \n", *reqStatus)
					// add it to the request channel for the consumer
					b.requestResponseChan <- reqStatus

					wgr.Done()
				}()
			}
			wgr.Wait()
		}
		log.Infof("%d number of queries requeted", i)
		wg.Done()
	}()

	log.Infof("Benchmark %s has been launched at %s", b.BenchmarkName, tinit.String())

	// wait untill both go routines are done
	wg.Wait()
	b.FinishTime = time.Now()
	b.TotalDuration = time.Since(tinit)

	log.Info("generating the summary")
	fmt.Println()

	// print result of the benchmark
	bSummary := b.PrintSummary()
	fmt.Println(bSummary)

	// export the results of the summary
	err := b.ExportResults()
	if err != nil {
		log.Error(err.Error())
	}
}

func (b *Benchmark) ExportResults() error {
	log.Info("Exporting results")
	// Generate results folder
	tnow := time.Now()
	year, month, day := tnow.Date()
	folderName := "results/" + b.BenchmarkName + "-" + fmt.Sprintf("%d", day) + "_" + month.String() + "_" + fmt.Sprintf("%d", year) + "_" + fmt.Sprintf("%d", tnow.Hour()) + ":" + fmt.Sprintf("%d", tnow.Minute())
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		return err
	}
	log.Infof("Folder %s created", folderName)

	// make csv for the summary
	csvFile, err := os.Create(folderName + "/query_results.csv")
	if err != nil {
		return err
	}
	csvFile.WriteString(RequestStatusCsvColumnNames())
	for _, req := range b.requestList {
		csvFile.WriteString(req.CsvLine())
	}
	log.Infof("CSV %s created", folderName+"/query_results.csv")
	csvFile.Close()

	// make csv for the requests done
	bFile, err := os.Create(folderName + "/benchmark_summary.txt")
	if err != nil {
		return err
	}
	summary := b.PrintSummary()
	bFile.WriteString(summary)
	log.Infof("Summary created at %s ", folderName+"/benchmark_summary.txt")
	bFile.Close()

	return nil
}

func (b *Benchmark) PrintSummary() string {

	summary := "Benchmark:\t" + b.BenchmarkName + "\n"
	summary += "Start Time:\t " + b.StartTime.String() + "\n"
	summary += "Finish Time:\t " + b.FinishTime.String() + "\n"
	summary += "Query Base:\t " + b.QueryBase + "\n"
	summary += "Query Timeout (secs):\t " + fmt.Sprintf("%d", timeout) + "\n"
	summary += "Average Resp Time:\t " + b.AvgReqTime.String() + "\n"
	summary += "Best Score:\t " + b.BestScore.String() + "\n"
	summary += b.BestScoreQuery + "\n"
	summary += "Worst Score:\t " + b.WorstScore.String() + "\n"
	summary += b.WorstScoreQuery + "\n"
	summary += "Total Queries Done:\t " + fmt.Sprintf("%d", b.TotQueries) + "\n"
	//summary += "Total Timeouts:\t " + fmt.Sprintf("%d", b.TimeOut) + "\n"
	summary += "Successful Requests:\t " + fmt.Sprintf("%d", b.SuccessReq) + "\n"
	summary += "Failed Requests:\t " + fmt.Sprintf("%d", b.FailReq) + "\n"
	summary += "Success Ratio:\t " + fmt.Sprintf("%.2f", b.SuccessRatio) + "% \n"

	return summary
}
