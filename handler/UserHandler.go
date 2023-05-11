package handler

import (
	"autentification_service/model"
	"autentification_service/service"
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

const SecretKeyForJWT = "v123v1iy2v321sdasada8"

type UserHandler struct {
	*pb.UnimplementedAutentificationServiceServer
	UserService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
	}
}

func (loginHandler *UserHandler) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	user := MapUserDTOFromLoginRequest(in)

	loggedUser, err := loginHandler.UserService.FindByEmail(user.Email)

	if err != nil {
		return &pb.LoginResponse{
			Message: "User not found!",
		}, status.Error(codes.OK, "User not found!")
	}

	if err1 := bcrypt.CompareHashAndPassword(loggedUser.Password, []byte(user.Password)); err1 != nil {
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
		println(err1.Error())
		return nil, err1
	}

	println("grpc metoda uspjesno izvrsena")
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
	//obrisi
	println("Pogodjen registration")

	user := MapUserFromRegistrationRequest(in)

	//obrisi
	println("Podaci")
	println(user.Email)

	err := handler.UserService.Create(user)
	if err != nil {
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
