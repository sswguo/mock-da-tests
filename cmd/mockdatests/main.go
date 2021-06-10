package main

import (
	"sync"
	"time"
	"regexp"
	"os"
	"fmt"
	"strings"
	datest "github.com/Commonjava/indy/mockdatests/datests"
)


func main() {

	c := datest.LoadConfig()

	buildId := os.Args[1]

	fmt.Println("buildId: ", buildId)

	pncRest := c.PncRest
	indyUrl := c.IndyUrl
	daGroup := c.DAGroup

	url := fmt.Sprintf("%s/builds/%s/logs/align", pncRest, buildId)

	fmt.Println(url)

	alignLog := datest.GetAlignLog(url)

	fmt.Println(alignLog)

	// extract the gav list from alignment log
	var re = regexp.MustCompile(`(?s)REST Client returned.*?\}`)

	jobs := 0
	var urls [1000]string
	var gavA [1000]string

	for _, match := range re.FindAllString(alignLog, -1) {

		gavs := match[len("REST Client returned {"):len(match)-1]
		
		gavArray := strings.Split(gavs, ",")

		for _, gav := range gavArray {

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
			datest.LookupMetadata(gavA[i], urls[i])
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