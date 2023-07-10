package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func HelloKubernetes(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Kubernetes!")
}

func GetEnv(w http.ResponseWriter, r *http.Request) {
	envs := os.Environ()
	io.WriteString(w, fmt.Sprintf("System Env: %+v", envs))
}

func main() {
	http.HandleFunc("/", HelloKubernetes)
	http.ListenAndServe(":8080", nil)
}
