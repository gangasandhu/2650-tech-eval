package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
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

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	cur, err := todoCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var todo Todo
		err := cur.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	json.NewDecoder(r.Body).Decode(&newTodo)
	newTodo.ID = primitive.NewObjectID()
	newTodo.Done = false
	_, err := todoCollection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objID}
	_, err = todoCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bson.M{"message": "Todo deleted"})
}

func updateTodoStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	json.NewDecoder(r.Body).Decode(&updatedTodo)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"done": updatedTodo.Done}}
	_, err = todoCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTodo)
}

func updateTodoTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	json.NewDecoder(r.Body).Decode(&updatedTodo)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"title": updatedTodo.Title}}
	_, err = todoCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTodo)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectToDB()
	handleRequests()
}
