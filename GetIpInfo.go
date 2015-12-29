package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Result Example: {"code":0,"data":{"country":"\u4e2d\u56fd","country_id":"CN","area":"\u534e\u5357","area_id":"800000","region":"\u5e7f\u4e1c\u7701","region_id":"440000","city":"\u6df1\u5733\u5e02","city_id":"440300","county":"","county_id":"-1","isp":"\u7535\u4fe1","isp_id":"100017","ip":"59.37.125.73"}}
type IpResult struct {
	Code int
	Data interface{}
}

func GetIpInfo(ip string) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	url := "http://ip.taobao.com//service/getIpInfo.php?ip=" + ip
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	info, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}

	var m IpResult
	err = json.Unmarshal(info, &m)
	if err != nil {
		logger.Fatal(err)
	}
	if m.Code == 0 { // success
		data := m.Data.(map[string]interface{})
		fmt.Printf("%s %s %s %s\n", ip, data["region"], data["city"], data["isp"])
	} else { // fail
		fmt.Printf("IP: %s ReturnCode: %d, ErrString: %s\n", ip, m.Code, m.Data.(string))
	}
}

func main() {
	var qps int
	flag.IntVar(&qps, "p", 10, "parallism number per second.")
	flag.Parse()
	var wg sync.WaitGroup
	rate := time.Second / time.Duration(qps) // ip.taobao.com limit 10 qps
	throttle := time.Tick(rate)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		<-throttle // rate limit our request
		go func(ip string) {
			defer wg.Done()
			GetIpInfo(ip)
		}(strings.TrimSpace(scanner.Text()))
	}

	wg.Wait()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
