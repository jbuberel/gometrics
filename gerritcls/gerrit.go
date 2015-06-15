package gerritcls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type GerritCL struct {
	Id       string `json:"id"`
	Project  string `json:"project"`
	ChangeId string `json:"change_id"`
	Subject  string `json:"subject"`
	Status   string `json:"status"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Owner    struct {
		Name    string `json:"owner"`
		Account int `json:"_account_id"`
	} `json:"owner"`
}

var search string = "https://go-review.googlesource.com/changes/?q=status:open&n=2000"

func Capture(dirname *string) {
	resp, err := http.Get(search)
	if err != nil {
		log.Printf("Unable to retrieve %v\n", search)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	bodyText := strings.SplitN(string(body[:]), "\n", 2)[1]
	log.Println(len(bodyText))

	var result []GerritCL
	err = json.Unmarshal([]byte(bodyText), &result)
	if err != nil {
		log.Println("error:", err)
	}

    countsByStatus := make(map[string]int)
    total := 0
	for _, r := range result {
		fmt.Printf("ID: %v - Status: %v - Owner: %v\n", r.ChangeId, r.Status, r.Owner.Name)
		countsByStatus[r.Status] += 1
		total += 1
	}
	
    
    timestamp := time.Now().Format("2006-01-02")
	outfile := fmt.Sprintf("%v/gerrit-open-cls-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

    // NEW, SUBMITTED, MERGED, ABANDONED, DRAFT
	f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", timestamp, total, countsByStatus["NEW"], countsByStatus["SUBMITTED"],
	    countsByStatus["MERGED"], countsByStatus["ABANDONED"], countsByStatus["DRAFT"] ))
}
