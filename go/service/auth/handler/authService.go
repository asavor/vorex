package handler

import (
	pb "auth/proto"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
)

type AuthServiceServer struct {
	RedisDB *redis.Client
	pb.UnimplementedAuthServiceServer
}

type userModel struct {
	Password []byte `redis:"password"`
	UserID   string `redis:"userID"`
}

func (s *AuthServiceServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var userStore userModel

	if err := s.RedisDB.HGetAll(ctx, in.GetEmail()).Scan(&userStore); err != nil {
		fmt.Println(err)
		return &pb.LoginResponse{
			Success:  false,
			Error:    "Invalid Username Or Password",
			JWTToken: "",
		}, nil
	}

	err := bcrypt.CompareHashAndPassword(userStore.Password, []byte(in.GetPassword()))
	if err != nil {
		return &pb.LoginResponse{
			Success:  false,
			Error:    "Invalid Username Or Password",
			JWTToken: "",
		}, nil
	}

	prvKey, err := ioutil.ReadFile("../auth/secret/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userID":    userStore.UserID,
		"ExpiresAt": 1000,
	})

	SignedToken, err := token.SignedString(signKey)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &pb.LoginResponse{
		Success:  true,
		Error:    "",
		JWTToken: SignedToken,
	}, nil

}

func (s *AuthServiceServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	_, err := s.RedisDB.Get(ctx, in.GetEmail()).Result()
	if err != redis.Nil {
		return &pb.RegisterResponse{
			Success:  false,
			Error:    "Email already exist",
			JWTToken: "",
		}, nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	userID := uuid.New().String()

	if _, err := s.RedisDB.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, in.GetEmail(), "password", hashedPassword)
		rdb.HSet(ctx, in.GetEmail(), "userID", userID)
		return nil
	}); err != nil {
		panic(err)
	}
	prvKey, err := ioutil.ReadFile("../auth/secret/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userID":    userID,
		"ExpiresAt": 1000,
	})

	SignedToken, err := token.SignedString(signKey)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &pb.RegisterResponse{
		Success:  true,
		Error:    "",
		JWTToken: SignedToken,
	}, nil

}
