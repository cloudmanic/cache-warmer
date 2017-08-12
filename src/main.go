package main

import (
  "fmt"
  "log"
  "time"
  "errors"
  "net/http"
  "io/ioutil"
  "encoding/xml"
  "encoding/json"
)

const httpPort = "8080"
const accessKey = "<ChangeMe>"
const workerCount = 10

var (
  jobs = make(chan string, 100000)  
)

type Query struct {
  XMLName xml.Name `xml:"urlset"`
  Locs []Loc `xml:"url>loc"`
}

type Loc string

//
// Main.
//
func main () {
  startWorkers()
  startWebserver()  
}

//
// Start web server
//
func startWebserver() {
  
  // Register some handlers:
  mux := http.NewServeMux()
  
  // Http Routes
  mux.HandleFunc("/", doSitemap)    

  // Start server on port
  s := &http.Server{
    Addr: ":" + httpPort,
    Handler: mux,
    ReadTimeout:  2 * time.Second,
    WriteTimeout: 2 * time.Second,			
  }
  
  fmt.Println("Starting web server http://localhost:" + httpPort)
  
  log.Fatal(s.ListenAndServe()) 
     
}

//
// Start wokers
//
func startWorkers() {
    
  for w := 1; w <= workerCount; w++ {
    go worker(w, jobs)
  }  
  
}

//
// Worker
//
func worker(id int, urls <-chan string) {
  
  for row := range urls {
    fmt.Println("worker", id, "fectching url : ", row)
    
    err := downloadPage(row) 
    
    if err != nil {
      fmt.Println(err)
    }
  }  
  
}

//
// Grab the sitemap and warm the cache.
//
func doSitemap(w http.ResponseWriter, r *http.Request) {

  // Make sure this is a post request.
	if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
	} 
	
  // Decode json passed in
  decoder := json.NewDecoder(r.Body)
  
  type Post struct {
    Token string
    Sitemap string
  }
  
  var p Post 
  
  err := decoder.Decode(&p)
  
  defer r.Body.Close()  
  
  if err != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong. Sorry for the trouble.\"}"))   
    return 
  }
  
  // Check access token
  if accessKey != p.Token {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Invalid access token\"}"))   
    return    
  }

  fmt.Println("New sitemap : ", p.Sitemap)

  // Download the sitemap
  query, err := downloadSiteMap(p.Sitemap)
  
  if err != nil {
    panic(err)
  }

  // Loop through the different pages
  for _, row := range query.Locs {
    jobs <- string(row)
  }

  // Return success json.
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("{\"status\":1}"))   
  
}

//
// Access one page in the sitemap so we warm the cache
//
func downloadPage(url string) error {
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup request
  req, _ := http.NewRequest("GET", url, nil)
 
  res, err := client.Do(req)
      
  if err != nil {
    return err  
  }        
  
  // Close Body
  defer res.Body.Close() 
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return errors.New(fmt.Sprint("Sitemap did not return 200, It returned ", res.StatusCode, "on url : ", url)) 
  }    
  
  return nil  
  
}

//
// Download the sitemap
//
func downloadSiteMap(sitemap string) (Query, error) {
  
  var q Query
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup request
  req, _ := http.NewRequest("GET", sitemap, nil)
  req.Header.Set("Accept", "application/xml") 
 
  res, err := client.Do(req)
      
  if err != nil {
    return q, err  
  }        
  
  // Close Body
  defer res.Body.Close()    
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return q, errors.New(fmt.Sprint("Sitemap did not return 200, It returned ", res.StatusCode)) 
  }    
     
  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body) 

  // Unmarshal xml
  xml.Unmarshal(body, &q)
  
  return q, nil
    
}

/* End File */