package importer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"form3/business"
)

const EmployeesUrl = "https://gist.githubusercontent.com/hcliff/5bf865e0c5e849c726a2478908860303/raw/2ece0a10c62ef52fc9edd2fac1b18a81a24fa13d/employees.json"
const GiftsUrl = "https://gist.githubusercontent.com/hcliff/5bf865e0c5e849c726a2478908860303/raw/2ece0a10c62ef52fc9edd2fac1b18a81a24fa13d/gifts.json"

func Fetch() (categories []*business.Category, employees []*business.Employee, gifts []*business.Gift, err error) {
	employeesBody, err := getContent(EmployeesUrl)
	if err != nil {
		return
	}
	jsonEmployees := make([]*Employee, 0)
	if err = json.Unmarshal(employeesBody, &jsonEmployees); err != nil {
		return
	}

	giftsBody, err := getContent(GiftsUrl)
	if err != nil {
		return
	}
	jsonGifts := make([]*Gift, 0)
	if err = json.Unmarshal(giftsBody, &jsonGifts); err != nil {
		return
	}

	converter := NewConverter(jsonEmployees, jsonGifts)

	return converter.Categories, converter.Employees, converter.Gifts, nil
}

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %w", err)
	}

	return data, nil
}
