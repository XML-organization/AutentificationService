package handler

import (
	"autentification_service/model"
	"autentification_service/service"
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/XML-organization/common/proto/autentification_service"
	userServicepb "github.com/XML-organization/common/proto/user_service"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
	SecretKeyForJWT     = "v123v1iy2v321sdasada8"
)

type UserHandler struct {
	*pb.UnimplementedAutentificationServiceServer
	UserService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
	}
}

func (handler *UserHandler) AutorizeUser(ctx context.Context, in *pb.AuthorizeUserRequest) (*pb.AuthorizeUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("Unauthorized. Missing metadata")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("Missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		log.Println("Unauthorized. missing authorization header")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		log.Println("Unauthorized. invalid authorization header format")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		log.Println("Unauthorized. unsupported authorization type:", authType)
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	token, claims := handler.UserService.GetClaimsFrowJwt(accessToken)
	if token == nil {
		log.Println("Unauthorized. invalid access token")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("invalid access token")
	}

	//token validation
	if !token.Valid {
		log.Println("Unauthorized. Token is not valid")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("Token is not valid")
	}

	//role validation
	if int(in.Role) != claims.Role {
		log.Println("Unauthorized. Unauthorized")
		return &pb.AuthorizeUserResponse{
			Message: "Unauthorized",
		}, fmt.Errorf("Unauthorize")
	}
	log.Println("Access granted")
	return &pb.AuthorizeUserResponse{
		Message: "Access granted",
	}, nil
}

func (loginHandler *UserHandler) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	user := MapUserDTOFromLoginRequest(in)

	loggedUser, err := loginHandler.UserService.FindByEmail(user.Email)

	if err != nil {
		log.Println("User not found!")
		log.Println(err)
		return &pb.LoginResponse{
			Message: "User not found!",
		}, status.Error(codes.OK, "User not found!")
	}

	if err1 := bcrypt.CompareHashAndPassword(loggedUser.Password, []byte(user.Password)); err1 != nil {
		log.Println("Password is incorrect")
		log.Println(err)
		return &pb.LoginResponse{
			Message: "Password is incorrect!",
		}, status.Error(codes.OK, "Password is incorrect!")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&model.JwtClaims{
			Id:   loggedUser.ID,
			Role: int(loggedUser.Role),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			}})

	token, err := claims.SignedString([]byte(SecretKeyForJWT))

	if err != nil {
		log.Println("Some error ocurred, please try again!")
		log.Println(err)
		return &pb.LoginResponse{
			Message: "Some error ocurred, please try again!",
		}, status.Error(codes.OK, err.Error())
	}

	httpRespHeader := metadata.New(map[string]string{
		"Set-Cookie": "jwt=" + token + "; HttpOnly; SameSite=Strict",
	})

	grpc.SendHeader(ctx, httpRespHeader)

	conn, err := grpc.Dial("user_service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userService := userServicepb.NewUserServiceClient(conn)

	getUserByEmailResponse, err1 := userService.GetUserByEmail(context.TODO(), &userServicepb.GetUserByEmailRequest{Email: loggedUser.Email})

	if err1 != nil {
		log.Println(err1.Error())
		return nil, err1
	}

	println(getUserByEmailResponse.Email)
	id := strings.Split(getUserByEmailResponse.Id, " |")[1]
	println(id)

	return &pb.LoginResponse{
		Id:          id,
		Name:        getUserByEmailResponse.Name,
		Surname:     getUserByEmailResponse.Surname,
		Email:       getUserByEmailResponse.Email,
		Role:        pb.Role(getUserByEmailResponse.Role),
		Country:     getUserByEmailResponse.Country,
		City:        getUserByEmailResponse.City,
		Street:      getUserByEmailResponse.Street,
		Number:      getUserByEmailResponse.Number,
		AccessToken: token,
		Message:     "Login successful!",
	}, err
}

/*func (loginHandler *UserHandler) User(writer http.ResponseWriter, req *http.Request) {

	cookies := req.Cookies()

	if len(cookies) == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unautenticated!"))
		return
	}

	claims := &model.JwtClaims{}

	token, err := jwt.ParseWithClaims(cookies[0].Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKeyForJWT), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("JWT invalid!"))
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad request!"))
		return
	}

	if !token.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	//ObjectID(\"64199dfe142552ce4d25b90b\")
	id := strings.Split(claims.Id.String(), "\"")[1]

	user, err := loginHandler.UserService.FindUser(id)

	json.NewEncoder(writer).Encode(user)
	writer.WriteHeader(http.StatusOK)
}

func (loginHandler *UserHandler) Logout(writer http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Logout success!"))
}*/

func (handler *UserHandler) Registration(ctx context.Context, in *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {

	user := MapUserFromRegistrationRequest(in)

	err := handler.UserService.Create(user)
	if err != nil {
		log.Println(err)
		if err.Error() == "user already exist" {
			err = nil
			return &pb.RegistrationResponse{
				Message: "User already exist!",
			}, err
		}
		return &pb.RegistrationResponse{
			Message: "Error occured, please try again!",
		}, err
	}

	return &pb.RegistrationResponse{
		Message: "Registration successful!",
	}, err
}

func (handler *UserHandler) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (handler *UserHandler) ChangeEmail(ctx context.Context, in *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeEmail not implemented")
}

func (handler *UserHandler) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	userC, err := handler.UserService.FindByEmail(request.Email)
	println("email aut" + request.Email)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err3 := handler.UserService.DeleteUserCredentials(userC)
	if err3 != nil {
		log.Println(err3)
		panic(err3)
	}

	return &pb.DeleteUserResponse{
		Message: "ok",
	}, err
}
