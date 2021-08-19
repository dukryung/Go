module Go/Projects/SideProject

go 1.15

require (
	cloud.google.com/go/storage v1.16.0
	firebase.google.com/go v3.13.0+incompatible
	firebase.google.com/go/v4 v4.4.0
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/validator/v10 v10.5.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.5 // indirect
	github.com/unrolled/render v1.4.0
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/api v0.54.0
	google.golang.org/appengine v1.6.7
	gopkg.in/yaml.v2 v2.4.0 // indirect
	sideproject.com/route v0.0.0
	sideproject.com/database v0.0.0
	sideproject.com/user v0.0.0
	sideproject.com/project v0.0.0
	sideproject.com/auth v0.0.0
	sideproject.com/profile v0.0.0
	sideproject.com/artist v0.0.0
	sideproject.com/gcloud v0.0.0
	sideproject.com/common v0.0.0

)

replace (
	sideproject.com/route v0.0.0 => ./route
	sideproject.com/database v0.0.0 => ./route/database
	sideproject.com/user v0.0.0 => ./route/user
	sideproject.com/project v0.0.0 => ./route/project
	sideproject.com/auth v0.0.0 => ./route/auth
	sideproject.com/profile v0.0.0 => ./route/profile
	sideproject.com/artist v0.0.0 => ./route/artist
	sideproject.com/gcloud v0.0.0 => ./route/gcloud
	sideproject.com/common v0.0.0 => ./route/common
)