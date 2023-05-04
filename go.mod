module autentification_service

go 1.20

replace github.com/XML-organization/common => ../common

require (
	github.com/XML-organization/common v1.0.1-0.20230504130318-3a270381458d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.3.0
	golang.org/x/crypto v0.8.0
	golang.org/x/net v0.9.0
	google.golang.org/grpc v1.54.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
