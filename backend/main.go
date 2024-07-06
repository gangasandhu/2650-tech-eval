package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the To-Do List!")
}

func handleRequests() {
    r := mux.NewRouter()
    r.HandleFunc("/", homePage).Methods("GET")
    r.HandleFunc("/todos", getTodos).Methods("GET")
    r.HandleFunc("/todo", createTodo).Methods("POST")
    r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
    r.HandleFunc("/todo/{id}/status", updateTodoStatus).Methods("PUT")
    r.HandleFunc("/todo/{id}/title", updateTodoTitle).Methods("PUT")

    // Set up CORS
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    connectToDB()
    handleRequests()
}
