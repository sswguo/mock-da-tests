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

func deleteMetadata(url string) string{
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer ")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return "Done"
}

func fetchMetadata(gav string, url string) string {
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

	//fmt.Println("XML:")
	//io.Copy(os.Stdout, resp.Body)

	// Create new file
	file := strings.ReplaceAll(gav, ":", "-")
	tmp, err := os.Create("results/" + file + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

	bytesWritten, err := io.Copy(tmp, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)

	return "Done"
}

func doRun(c config, ds dataset){
	indyUrl := c.IndyUrl
	daGroup := c.DAGroup

	results := make(chan string)

	jobs := 0
	var urls [10000]string
	var gavs [10000]string

	for idx, element := range ds.Artifacts {
		fmt.Println(idx, element.Name)
		for idx_, gav := range element.Gavs {
		    gavs[jobs] = gav
			fmt.Println(idx_, gav)
			s := strings.Split(gav, ":")
			groupId := strings.ReplaceAll(s[0], ".", "/")
			artifactId := s[1]
			url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyUrl, daGroup, groupId, artifactId)
			urls[jobs] = url
			jobs = jobs+1
		}
	}

	fmt.Println("Total jobs:")
    fmt.Println(jobs)

	concurrentGoroutines := make(chan struct{}, c.MaxConcurrentGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < jobs; i++ {
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("doing", i)
			start := time.Now()
			fetchMetadata(gavs[i], urls[i])
			//deleteMetadata(urls[i])
			elapsed := time.Since(start)
			fmt.Println("finished", i)
			fmt.Println("took", elapsed)
			<-concurrentGoroutines
		}(i)
	}

	for i := 0; i < jobs; i++ {
		fmt.Println(<-results)
	}

	wg.Wait()
}

func main() {
	c := loadConfig()
	ds := loadData()
	doRun(c, ds)
}
