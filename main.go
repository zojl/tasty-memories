package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/fs"
    "log"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

const (
    templatePath = "./templates/index.html"
    memoriesDir  = "./memories"
)

type RawMemory struct {
    Text  string `json:"text"`
    Image string `json:"image,omitempty"`
}

type Memory struct {
    Content []string
    ImageSrc template.HTMLAttr
    Username string
}

type PageParams struct {
    AnalyticsTag template.HTML
    Memory *Memory
}

func main() {
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        log.Fatalf("Error parsing template: %v", err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/memories/") {
            http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
        } else {
            var memory *Memory
            if strings.HasPrefix(r.URL.Path, "/~") {
                urlParts := strings.SplitN(r.URL.Path[2:], "/", 2)
                if len(urlParts) > 1 {
                    http.Redirect(w, r, "/~" + urlParts[0], http.StatusFound)
                }
                memory = loadMemory(urlParts[0])
            }

            if memory == nil && r.URL.Path != "/" {
                http.Redirect(w, r, "/", http.StatusFound)
                return
            }
            
            if memory == nil {
                memory = loadRandomMemory()
            }

            params := PageParams {
                AnalyticsTag: template.HTML(getEnv("ANALYTICS_TAG", "")),
                Memory: memory,
            }
            
            if err := tmpl.Execute(w, params); err != nil {
                http.Error(w, "Error rendering template", http.StatusInternalServerError)
                log.Printf("Error executing template: %v", err)
            }
        }
    })

    http.HandleFunc("/memories", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            files, err := os.ReadDir(memoriesDir)
            if err != nil {
                http.Error(w, "Error reading directory", http.StatusInternalServerError)
                log.Printf("Error reading memories directory: %v", err)
                return
            }

            var filenames []string
            for _, file := range files {
                if !file.IsDir() {
                    filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
                    filenames = append(filenames, filename)
                }
            }

            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(filenames); err != nil {
                http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
                log.Printf("Error encoding JSON: %v", err)
            }
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    serverPort := getEnv("SERVER_PORT", "8080")
    log.Println("Server is running on http://localhost:" + serverPort)
    if err := http.ListenAndServe(":" + serverPort, nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

func loadMemory(name string) *Memory {
    if !isMemoryExist(name) {
        return nil
    }
    
    fileName := filepath.Join(memoriesDir, name + ".json")
    file, err := os.Open(fileName)
    if err != nil {
        return nil
    }
    defer file.Close()

    var rawMemory RawMemory
    var memory Memory
    if err := json.NewDecoder(file).Decode(&rawMemory); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        return nil
    }
    memory.Content = strings.Split(rawMemory.Text, "\n")
    
    if rawMemory.Image != "" {
        memory.ImageSrc = template.HTMLAttr(fmt.Sprintf("src=\"%s\"", rawMemory.Image))
    }
    
    memory.Username = name
    return &memory
}

func isMemoryExist(name string) bool {
    files, err := os.ReadDir(memoriesDir)
    if err != nil {
        log.Printf("Error reading memories directory: %v", err)
        return false
    }
    
    for _, file := range files {
        if !file.IsDir() && file.Name() == name + ".json" {
            return true
        }
    }
    
    return false
}

func loadRandomMemory() *Memory {
    files, err := os.ReadDir(memoriesDir)
    if err != nil {
        log.Printf("Error reading memories directory: %v", err)
        return nil
    }

    var jsonFiles []fs.DirEntry
    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
            jsonFiles = append(jsonFiles, file)
        }
    }

    if len(jsonFiles) == 0 {
        log.Println("No JSON files found in memories directory")
        return nil
    }

    randomFile := jsonFiles[rand.Intn(len(jsonFiles))]
    return loadMemory(strings.TrimSuffix(randomFile.Name(), ".json"))
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}