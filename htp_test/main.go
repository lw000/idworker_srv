// htp_test project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type UrlInfo struct {
	Url  string `json:url`
	Data string `json:data`
}

type Config struct {
	Count  int    `json:count`
	Method string `json:method`
	Url    string `json:url`
	Data   string `json:data`
}

var (
	wg *sync.WaitGroup
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfig() bool {
	data, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func httpPost(cout int, url string, data string) {
	for i := 0; i < cout; i++ {
		go func(i int, w *sync.WaitGroup) {
			w.Add(1)
			defer w.Done()
			resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
			if err != nil {
				fmt.Println(i, ":", err.Error())
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(i, ":", err.Error())
				return
			}

			if len(body) > 0 {

			}

			fmt.Printf("%04d %s\n", i, string(body))
		}(i, wg)
		time.Sleep(time.Millisecond * 1)
	}
}
func httpGet(cout int, url string) {
	for i := 0; i < cout; i++ {
		go func(i int, w *sync.WaitGroup) {
			w.Add(1)
			defer w.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(i, err.Error())
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(i, err.Error())
				return
			}

			if len(body) > 0 {

			}
			fmt.Printf("%04d %s\n", i, string(body))
		}(i, wg)
		// time.Sleep(time.Millisecond * 1)
	}
}
func main() {
	cnf := NewConfig()
	if !cnf.LoadConfig() {
		return
	}

	wg = &sync.WaitGroup{}

	start := time.Now()
	if cnf.Method == "get" {
		httpGet(cnf.Count, cnf.Url)
		// for _, v := range cnf.Url {
		// 	httpGet(cnf.Count, v.Url)
		// }
	} else {
		httpPost(cnf.Count, cnf.Url, cnf.Data)
		// for _, v := range cnf.Url {
		// 	httpPost(cnf.Count, v.Url, v.Data)
		// }
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("tms:", end.Sub(start))
}
