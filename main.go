package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getServiceName() string {
	return getEnv("K_SERVICE", "UNKNOWN")
}

func echoPath(r *http.Request) string {
	return fmt.Sprintf("%s service executing path %s \n", getServiceName(), r.URL.Path)
}

func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, echoPath(r))
	io.WriteString(w, "In root function \n")
}

func find(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, echoPath(r))
	io.WriteString(w, "In find function \n")
}

func serviceFind(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, echoPath(r))
	io.WriteString(w, "In serviceFind function \n")
}

func main() {
	http.HandleFunc("/", root)

	// This function is never called as the URL is always `/<SERVICE>/find`
	http.HandleFunc("/find", find)

	// This function is called as the URL is called as the path is `/<SERVICE>/find`
	path := fmt.Sprintf("/%s/find", getServiceName())
	fmt.Printf("%s is mapped to getServiceFind processor", path)
	http.HandleFunc(path, serviceFind)

	port := fmt.Sprintf(":%s", getEnv("PORT", "8080"))
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal(err)
	}
}
