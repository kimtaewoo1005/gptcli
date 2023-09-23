package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY environment variable is not set.")
		return
	}

	fmt.Println("Please enter your multi-line content. Press Ctrl-D (Unix) or Ctrl-Z (Windows) when done:")

	reader := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// EOF signals end of input.
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading input:", err)
			return
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	content := strings.Join(lines, "\n")

	// Create the JSON payload
	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    []map[string]string{{"role": "user", "content": content}},
		"temperature": 0.7,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Unmarshal the JSON response into a map
	var responseMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &responseMap); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Print the JSON response as is
	fmt.Println("JSON Response As Is:")
	fmt.Println(string(body))

	// Extract the "content" field
	choices := responseMap["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	contentReturned := message["content"].(string)

	// Pretty-print the "content" field
	prettyContent, err := json.MarshalIndent(contentReturned, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println("\nContent Only in Pretty Format:")
	fmt.Println(string(prettyContent))
}

