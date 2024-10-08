package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type MetricsResponse map[string]int

func parseMetricsContent(content string) (MetricsResponse, error) {
	var err error
	result := MetricsResponse{}

	// Validate the Content
	var isValid bool
	isValid, err = regexp.MatchString(`^#`, content)
	isValid = isValid && err == nil && content != ""
	if !isValid {
		return nil, errors.New("invalid metrics content")
	}

	// Regex for Metrics Format,
	// selects all content that does not start with # in multiline mode
	re := regexp.MustCompile(`(?m)^[^#].*`)
	rows := re.FindAllString(content, -1)

	// Add the strings to the map
	for _, row := range rows {
		var value int
		keyValue := strings.Split(row, " ")
		value, err = strconv.Atoi(keyValue[1])
		result[keyValue[0]] = value
	}

	// Handle any error results
	if err != nil {
		return nil, err
	}

	// Give the user what they asked for
	return result, nil
}

// GetMetrics parses the /metrics endpoint from algod into a map
func GetMetrics(server string, token string) (MetricsResponse, error) {

	// Create Request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/metrics", server), nil)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("unable to create a new http request")
	}

	// Add Token
	req.Header.Add("X-Algo-API-Token", token)

	// Execute Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the data
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseMetricsContent(string(content))
}
