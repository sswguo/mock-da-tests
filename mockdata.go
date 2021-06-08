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
	"sync"
	"time"
	"io"
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
		fmt.Println("configFile.Get err #%v ", err)
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

func lookupMetadata(gav string, url string) string {
	fmt.Println(url)
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

	// dump the metadata file locally for verifying
	tempArray := strings.Split(gav, "=")
	file := strings.ReplaceAll(tempArray[0], ":", "-")
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

func main() {

	c := loadConfigA()

	buildId := os.Args[1]

	fmt.Println(buildId)

	pncRest := c.PncRest
	indyUrl := c.IndyUrl
	daGroup := c.DAGroup

	url := fmt.Sprintf("%s/builds/%s/logs/align", pncRest, buildId)

	fmt.Println(url)

	alignLog := getAlignLog(url)

	fmt.Println(alignLog)

	// extract the gav list from alignment log
	var re = regexp.MustCompile(`(?s)REST Client returned.*?\}`)

	jobs := 0
	var urls [1000]string
	var gavA [1000]string

	for i, match := range re.FindAllString(alignLog, -1) {
		fmt.Println(match, "found at index", i)

		gavs := match[len("REST Client returned {"):len(match)-1]
		
		gavArray := strings.Split(gavs, ",")

		for idx, gav := range gavArray {
			fmt.Println(idx, gav)

			s := strings.Split(gav, ":")
			groupId := strings.Trim(s[0], " ")
			artifactId := s[1]

			fmt.Println("GroupID: ", groupId, " ArtifactId: ", artifactId)

			groupIdPath := strings.ReplaceAll(groupId, ".", "/")

			url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyUrl, daGroup, groupIdPath, artifactId)

			urls[jobs] = url
			gavA[jobs] = gav
			jobs = jobs+1
		}

		fmt.Println("Total jobs:", jobs, " for buildId:", buildId)
	}

	results := make(chan string)

	concurrentGoroutines := make(chan struct{}, c.MaxConcurrentGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < jobs; i++ {
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("Doing", i)
			start := time.Now()
			lookupMetadata(gavA[i], urls[i])
			elapsed := time.Since(start)
			fmt.Println("Finished #", i, " in ", elapsed)
			<-concurrentGoroutines
		}(i)
	}

	for i := 0; i < jobs; i++ {
		fmt.Println(<-results)
	}

	wg.Wait()

}