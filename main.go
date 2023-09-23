package main

import (
	"context"
	"github.com/auyeongweibin/auth-server/configs"
	"github.com/auyeongweibin/auth-server/otp"
	"github.com/auyeongweibin/auth-server/registration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os"
)

type myRegistrationServer struct {
	registration.UnimplementedRegistrationServer
}

func (s myRegistrationServer) CheckUserEnrolled(c context.Context, user *registration.UserEnrolledRequest) (*registration.UserEnrolledResponse, error) {
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
			return &registration.UserEnrolledResponse{
				Message: "Sorry you are not enrolled in our database, please make sure that you have an account with one of our partners!",
			}, nil
		}
		log.Fatal(err)
	}
	if result["status"] == "pending" {
		return &registration.UserEnrolledResponse{
			Message: "You are enrolled in in our database, please proceed to register an account with us!",
		}, nil
	}
	return &registration.UserEnrolledResponse{
		Message: "You are enrolled in in our database, please proceed to register an account with us!",
	}, nil
}

func (s myRegistrationServer) Register(c context.Context, userCredentials *registration.RegistrationRequest) (*registration.RegistrationResponse, error) {
	usersCollection := configs.GetCollection(configs.DB, "users")
	var result bson.M
	err := usersCollection.FindOne(
		context.TODO(),
		bson.D{
			{"email", userCredentials.Email},
		},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &registration.RegistrationResponse{
				Email: userCredentials.Email, 
				Message: "Sorry you are not enrolled in our database, please make sure that you have an account with one of our partners!",
			}, nil
		}
		log.Fatal(err)
	}
	if result["status"] != "pending" {
		return &registration.RegistrationResponse{
			Email: userCredentials.Email, 
			Message: "You have already registered, please proceed to login instead!",
		}, nil
	}
	
	return &registration.RegistrationResponse{
		Email: userCredentials.Email, 
		Message: "Registration successful, email verification sent!",
	}, nil
}


func GenerateOTP(n int) string {
	var numberRunes = []rune("0123456789")
    b := make([]rune, n)
    for i := range b {
        b[i] = numberRunes[rand.Intn(len(numberRunes))]
    }
    return string(b)
}

type myOTPServer struct {
	otp.UnimplementedOTPServer
}

func (s myOTPServer) GetOTP(c context.Context, user *otp.OTPRequest) (*otp.OTPResponse, error) {
	usersCollection := configs.GetCollection(configs.DB, "auth")
	var result bson.M
	err := usersCollection.FindOne(
		context.TODO(),
		bson.D{
			{"email", user.Email},
		},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &otp.OTPResponse {
				VerificationKey: "NA",
				Message: "Email Not Registered",
			}, nil
		}
		log.Fatal(err)
	}

	// Generate OTP
	otpCode := GenerateOTP(6)

	// TODO: Generate Verification Key
	key := "key"

	// Send Email
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST"))
	to := []string{user.Email}
	msg := []byte("To: "+user.Email+"\r\n" +
		"Subject: OTP for Account Verification\r\n" +
		"\r\n" +
		"OTP: "+ otpCode +".\r\n")
	emailErr := smtp.SendMail(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_HOST_PORT"), auth, os.Getenv("EMAIL_ADDRESS"), to, msg)
	if emailErr != nil {
		log.Fatal(emailErr)
		return &otp.OTPResponse {
			VerificationKey: "NA",
			Message: "Error Sending Email",
		}, nil
	}
	return &otp.OTPResponse {
		VerificationKey: key,
		Message: "Success",
	}, nil
}

func (s myOTPServer) VerifyOTP(c context.Context, user *otp.VerifyOTPRequest) (*otp.VerifyOTPResponse, error) {
	// TODO: Implement Verification Key and OTP Check
	
	return &otp.VerifyOTPResponse{
		Status: "Success",
		Details: "OTP Matched",
		Email: user.Email,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	registrationService := &myRegistrationServer{}
	registration.RegisterRegistrationServer(serverRegistrar, registrationService)
	otpService := &myOTPServer{}
	otp.RegisterOTPServer(serverRegistrar, otpService)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}