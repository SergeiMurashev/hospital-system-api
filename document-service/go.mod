module github.com/sergeimurashev/hospital-system-api/document-service

go 1.21

require (
	github.com/elastic/go-elasticsearch/v8 v8.11.1
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.5.1
	github.com/sergeimurashev/hospital-system-api/proto v0.0.0
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)

replace github.com/sergeimurashev/hospital-system-api/proto => ./proto 