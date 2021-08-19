module user

go 1.15

require (
	github.com/gin-gonic/gin v1.7.4
	google.golang.org/appengine v1.6.7
	sideproject.com/auth v0.0.0
)

replace (
sideproject.com/auth v0.0.0 => ../auth
)
