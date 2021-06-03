package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type config struct{
    IndyUrl string `yaml:"indy_url"`
	DAGroup string `yaml:"da_group"`
	MaxConcurrentGoroutines int `yaml:"max_concurrent_goroutines"`
}

type entry struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Gavs    []string `yaml:"gavs"`
}

type dataset struct {
	Artifacts []entry `yaml:"artifacts"`
}

func loadData() dataset {

	fmt.Println("load data from dataset.yaml")

	yamlFile, err := ioutil.ReadFile("dataset.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get err   #%v ", err)
	}

	ds := dataset{}
	err = yaml.Unmarshal(yamlFile, &ds)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Println(ds)

	return ds
}

func loadConfig() config {
	fmt.Println("load config from config.yaml")

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("configFile.Get err   #%v ", err)
	}

	c := config{}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Println(c)

	return c
}

func fetchMetadata(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/xml")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("XML:")
	io.Copy(os.Stdout, resp.Body)
	return "Done"
}

func mockRemoteReq() {
	time.Sleep(500 * time.Millisecond)
}

func main() {
	fmt.Println("Rest query metadata ...")

	c := loadConfig()

	indyUrl := c.IndyUrl
	daGroup := c.DAGroup

	ds := loadData()

	results := make(chan string)

	jobs := 0
	var urls [50]string

	for idx, element := range ds.Artifacts {
		fmt.Println(idx, element.Name)
		for idx_, gav := range element.Gavs {
			fmt.Println(idx_, gav)
			s := strings.Split(gav, ":")
			groupId := strings.ReplaceAll(s[0], ".", "/")
			artifactId := s[1]
			url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyUrl, daGroup, groupId, artifactId)
			fmt.Println(url)
            // need to find a better way to collect the url
			urls[jobs] = url
			jobs = jobs+1
		}
	}

	fmt.Println("Total jobs:")
    fmt.Println(jobs)

	concurrentGoroutines := make(chan struct{}, c.MaxConcurrentGoroutines)
	var wg sync.WaitGroup

	for i := 0; i<jobs;i++{
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("doing", i)
			fetchMetadata(urls[i])
			//mockRemoteReq()
			fmt.Println("finished", i)
			<-concurrentGoroutines
		}(i)
	}

	for i := 0; i < jobs; i++ {
		fmt.Println(<-results)
	}

	wg.Wait()

}
