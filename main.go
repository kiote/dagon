package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "net/url"
    "encoding/json"
    "time"
)

type Issue struct {
    Title string `json:"title"`
    CreatedAt time.Time `json:"created_at"`
    URL       string    `json:"html_url"`
}

func main() {
    // Get the personal access token from the environment variable
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        panic("GITHUB_TOKEN environment variable not set")
    }
    
    owner := "kubernetes"
    repo := "kubernetes"
    tag := url.QueryEscape("good first issue")
    
    // Construct the API endpoint URL
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&state=open", owner, repo, tag)
    
    // Create a new HTTP request with the appropriate headers and token
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Authorization", "token " + token)
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    
    // Send the request and get the response
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Request failed with status code %d\n", resp.StatusCode)
        if resp.StatusCode == http.StatusBadRequest {
            fmt.Println("The request was invalid. Check the request parameters.")
        }
        return
    }

    // Read the response body into a byte array
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // Decode the JSON response into a slice of Issue structs
    var issues []Issue
    err = json.Unmarshal(body, &issues)
    if err != nil {
        panic(err)
    }

    // Print the titles of the issues
    for _, issue := range issues {
        fmt.Println(issue.Title)
        fmt.Printf("Created at: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
        fmt.Printf("URL: %s\n\n", issue.URL)
    }    
}


