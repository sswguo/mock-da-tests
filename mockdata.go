package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type configA struct{
	PncRest string `yaml:"pnc_rest_url"`
    IndyUrl string `yaml:"indy_url"`
	DAGroup string `yaml:"da_group"`
	MaxConcurrentGoroutines int `yaml:"max_concurrent_goroutines"`
}

func loadConfigA() configA {
	fmt.Println("load config from config.yaml")

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("configFile.Get err   #%v ", err)
	}

	c := configA{}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println(c)
	return c
}

func getAlignLog(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "text/plain")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    responseString := string(responseData)

	return responseString
}

func main() {


	c := loadConfigA()

	buildId := os.Args[1]

	fmt.Println(buildId)

	pncRest := c.PncRest

	url := fmt.Sprintf("%s/builds/%s/logs/align", pncRest, buildId)

	fmt.Println(url)

	alignLog := getAlignLog(url)

	fmt.Println(alignLog)

	var re = regexp.MustCompile(`(?s)REST Client returned.*?\}`)
    
    for i, match := range re.FindAllString(alignLog, -1) {
        fmt.Println(match, "found at index", i)

		gavs := match[len("REST Client returned {"):len(match)-1]

		fmt.Println(match[len("REST Client returned {"):len(match)-1])
		
		gavArray := strings.Split(gavs, ",")

		for idx, gav := range gavArray {
			fmt.Println(idx, gav)
			fmt.Println("-", gav)
		}
    }
}