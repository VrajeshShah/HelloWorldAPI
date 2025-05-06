package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	redisutils "github.com/VrajeshShah/HelloWorldAPI/utils"
)

type JsonRequest struct {
	Name string
}
type JsonResponse struct {
	Message string
}
type KeyValuePair struct {
	Key   string
	Value string
}

func GetEnvVariable(envName string, defaultValue string) string {
	rValue, present := os.LookupEnv(envName)
	if !present {
		rValue = defaultValue
	}
	return rValue
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rValue := JsonResponse{Message: "Hello World"}
	json.NewEncoder(w).Encode(rValue)
}

func HelloPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var request JsonRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Name == "" {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rValue := JsonResponse{Message: fmt.Sprintf("Hello %s", request.Name)}
	json.NewEncoder(w).Encode(rValue)
}

func HelloPostHandlerWithAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	authToken := strings.TrimSpace(r.Header.Get("x-auth-token"))
	if authToken == "" {
		http.Error(w, "x-auth-token empty or not found", http.StatusUnauthorized)
		return
	}
	if authToken != "ABC" {
		http.Error(w, "Invalid AuthToken", http.StatusUnauthorized)
		return
	}
	var request JsonRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Name == "" {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rValue := JsonResponse{Message: fmt.Sprintf("Hello %s", request.Name)}
	json.NewEncoder(w).Encode(rValue)
}

func RedisGetHandler(w http.ResponseWriter, r *http.Request) {
	var request KeyValuePair
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Key == "" {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	redisObject := redisutils.InitRedis()
	rValue := JsonResponse{Message: redisObject.Get(request.Key)}
	json.NewEncoder(w).Encode(rValue)
}
func RedisPutHandler(w http.ResponseWriter, r *http.Request) {
	var request KeyValuePair
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Key == "" || request.Value == "" {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	redisObject := redisutils.InitRedis()
	errorValue := redisObject.Set(request.Key, request.Value)
	msg := ""
	if errorValue != nil {
		msg = errorValue.Error()
	}
	rValue := JsonResponse{Message: msg}
	json.NewEncoder(w).Encode(rValue)
}

func main() {
	http.HandleFunc("/helloGet", HelloHandler)
	http.HandleFunc("/helloPost", HelloPostHandler)
	http.HandleFunc("/helloPostAuth", HelloPostHandlerWithAuth)
	http.HandleFunc("/Redis/Get", RedisGetHandler)
	http.HandleFunc("/Redis/Set", RedisPutHandler)
	fmt.Println("Server Staring on Port 8080")
	http.ListenAndServe(":8080", nil)
}
