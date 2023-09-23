package main

import (
	"context"
	"github.com/auyeongweibin/auth-server/auth"
	"github.com/auyeongweibin/auth-server/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myRegistrationServer struct {
	auth.UnimplementedRegistrationServer
}

func (s myRegistrationServer) CheckUserEnrolled(c context.Context, user *auth.UserEnrolledRequest) (*auth.UserEnrolledResponse, error) {
	usersCollection := configs.GetCollection(configs.DB, "users")
	var result bson.M
	err := usersCollection.FindOne(
		context.TODO(),
		bson.D{
			{"email", user.Email},
			{"first_name", user.FirstName},
			{"last_name", user.LastName},
			{"birthdate", user.Birthdate},
		},
	).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return &auth.UserEnrolledResponse{
				Message: "Sorry you are not enrolled in our database, please make sure that you have an account with one of our partners!",
			}, nil
		}
		log.Fatal(err)
	}
	if result["status"] == "pending" {
		return &auth.UserEnrolledResponse{
			Message: "You are enrolled in in our database, please proceed to register an account with us!",
		}, nil
	}
	return &auth.UserEnrolledResponse{
		Message: "You are enrolled in in our database, please proceed to register an account with us!",
	}, nil
}

func (s myRegistrationServer) Register(c context.Context, userCredentials *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	usersCollection := configs.GetCollection(configs.DB, "users")
	var result bson.M
	err := usersCollection.FindOne(
		context.TODO(),
		bson.D{
			{"email", userCredentials.Email},
		},
	).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return &auth.RegistrationResponse{
				Email: userCredentials.Email, 
				Message: "Sorry you are not enrolled in our database, please make sure that you have an account with one of our partners!",
			}, nil
		}
		log.Fatal(err)
	}
	if result["status"] != "pending" {
		return &auth.RegistrationResponse{
			Email: userCredentials.Email, 
			Message: "You have already registered, please proceed to login instead!",
		}, nil
	}

	// TODO: Send email verification
	
	return &auth.RegistrationResponse{
		Email: userCredentials.Email, 
		Message: "Registration successful, email verification sent!",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	registrationService := &myRegistrationServer{}
	auth.RegisterRegistrationServer(serverRegistrar, registrationService)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}