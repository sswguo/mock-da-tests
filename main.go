package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type config struct {
	PncRest                 string `yaml:"pnc_rest_url"`
	IndyUrl                 string `yaml:"indy_url"`
	DAGroup                 string `yaml:"da_group"`
	MaxConcurrentGoroutines int    `yaml:"max_concurrent_goroutines"`
}

func loadConfig() config {
	fmt.Println("load config from config.yaml")

	configFile, err := ioutil.ReadFile("config-local.yaml")
	if err != nil {
		fmt.Println("configFile.Get err #%v ", err)
	}

	c := config{}
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

	return string(responseData)
}

func lookupMetadata(url string) string {
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Message:") {
		fmt.Printf(bodyString)
	}

	//fmt.Printf("Bytes Written: %d\n", bodyBytes)

	return "Done"
}

func main() {

	//c := loadConfig()

	buildIds := os.Args[4] //os.Getenv("BUILD_ID") //

	fmt.Println("buildIds: ", buildIds)

	pncRest := os.Args[1] //os.Getenv("PNC_REST") //c.PncRest
	indyUrl := os.Args[2] //os.Getenv("INDY_URL") //c.IndyUrl
	daGroup := os.Args[3] //os.Getenv("DA_GROUP") //c.DAGroup
	goroutines := os.Args[5]

	routines, err := strconv.Atoi(goroutines)
	if err == nil {
		fmt.Println(routines)
	}

	buildIdArray := strings.Split(buildIds, ",")

	var urls []string

	for _, buildId := range buildIdArray {
		url := fmt.Sprintf("%s/builds/%s/logs/align", pncRest, buildId)

		fmt.Println(url)

		alignLog := getAlignLog(url)

		// extract the gav list from alignment log
		var re = regexp.MustCompile(`(?s)REST Client returned.*?\}`)

		var urlsTmp []string

		for _, match := range re.FindAllString(alignLog, -1) {

			gavs := match[len("REST Client returned {") : len(match)-1]

			gavArray := strings.Split(gavs, ",")

			for _, gav := range gavArray {

				s := strings.Split(gav, ":")
				groupId := strings.Trim(s[0], " ")
				artifactId := s[1]

				fmt.Println("GroupID: ", groupId, " ArtifactId: ", artifactId)

				groupIdPath := strings.ReplaceAll(groupId, ".", "/")

				url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyUrl, daGroup, groupIdPath, artifactId)

				urlsTmp = append(urlsTmp, url)
				urls = append(urls, url)

			}

			fmt.Println("Requests:", len(urlsTmp), " for buildId:", buildId)
		}
	}

	fmt.Println("Total requests: ", len(urls), "with routines:", routines)
	//results := make(chan string)

	concurrentGoroutines := make(chan struct{}, routines) //c.MaxConcurrentGoroutines)
	var wg sync.WaitGroup
    
	for i := 0; i < len(urls); i++ {
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("Doing", i)
			start := time.Now()
			lookupMetadata(urls[i])
			elapsed := time.Since(start)
			fmt.Println("Finished #", i, " in ", elapsed)
			<-concurrentGoroutines
		}(i)
	}

	//for i := 0; i < len(urls); i++ {
	//	fmt.Println(<-results)
	//}

	wg.Wait()

}
