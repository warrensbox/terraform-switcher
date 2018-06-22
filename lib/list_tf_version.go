package lib

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type tfVersionList struct {
	tflist []string
}

//GetTFList:  Get the list of available terraform version given the hashicorp url
func GetTFList(hashiURL string) ([]string, error) {

	/* Get list of terraform versions from hashicorp releases */
	resp, errURL := http.Get(hashiURL)
	if errURL != nil {
		log.Printf("Error getting url: %v", errURL)
		return nil, errURL
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Printf("Error reading body: %v", errBody)
		return nil, errBody
	}

	bodyString := string(body)
	result := strings.Split(bodyString, "\n")

	var tfVersionList tfVersionList

	for i := range result {
		//getting versions from body; should return match /X.X.X/
		r, _ := regexp.Compile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)

		if r.MatchString(result[i]) {
			str := r.FindString(result[i])
			trimstr := strings.Trim(str, "/") //remove "/" from /X.X.X/
			tfVersionList.tflist = append(tfVersionList.tflist, trimstr)
		}
	}

	return tfVersionList.tflist, nil

}
