package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Kelime_1 string             `json:"kelime_1,omitempty" bson:"kelime_1,omitempty"`
	Kelime_2 string             `json:"kelime_2,omitempty" bson:"kelime_2,omitempty"`
	Kelime_3 string             `json:"kelime_3,omitempty" bson:"kelime_3,omitempty"`
	Kelime_4 string             `json:"kelime_4,omitempty" bson:"kelime_4,omitempty"`
	Kelime_5 string             `json:"kelime_5,omitempty" bson:"kelime_5,omitempty"`
	Kelime_6 string             `json:"kelime_6,omitempty" bson:"kelime_6,omitempty"`
	Kelime_7 string             `json:"kelime_7,omitempty" bson:"kelime_7,omitempty"`
	Kelime_8 string             `json:"kelime_8,omitempty" bson:"kelime_8,omitempty"`
}

var (
	client   *mongo.Client
	mongoURL string = "mongodb://mongo:mongo@localhost:27017/?authSource=admin"
)

func CreateTitleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("DB").Collection("collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func GetTitleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("DB").Collection("collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func GetOneTitleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("DB").Collection("collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(person)
}

func GetTitleEndpointbyName(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	name := params["name"]
	fmt.Println(name)
	var people []Person
	collection := client.Database("DB").Collection("collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{{"$or",
		bson.A{
			bson.D{{"kelime_1", name}},
			bson.D{{"kelime_2", name}},
			bson.D{{"kelime_3", name}},
			bson.D{{"kelime_4", name}},
			bson.D{{"kelime_5", name}},
			bson.D{{"kelime_6", name}},
			bson.D{{"kelime_7", name}},
			bson.D{{"kelime_8", name}},
			bson.D{{"kelime_9", name}},
		}},
	})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func DeleteOneTitleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("DB").Collection("collection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.DeleteOne(ctx, Person{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
func main() {
	fmt.Println("Starting App..")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	defer client.Disconnect(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/title", CreateTitleEndpoint).Methods("POST")
	router.HandleFunc("/titles", GetTitleEndpoint).Methods("GET")
	router.HandleFunc("/titles/{name}", GetTitleEndpointbyName).Methods("GET")
	router.HandleFunc("/title/{id}", GetOneTitleEndpoint).Methods("GET")
	router.HandleFunc("/title/{id}", DeleteOneTitleEndpoint).Methods("DELETE")
	http.ListenAndServe(":12345", router)
}
