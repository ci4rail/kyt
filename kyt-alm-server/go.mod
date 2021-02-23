module github.com/ci4rail/kyt/kyt-alm-server

go 1.15

require (
	github.com/amenzhinsky/iothub v0.6.2
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/ci4rail/kyt/kyt-server-common/iothub_wrapper v0.0.0-00010101000000-000000000000
	github.com/ci4rail/kyt/kyt-server-common/token v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/golangci/golangci-lint v1.36.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a
)

replace github.com/ci4rail/kyt/kyt-server-common/token => ../kyt-server-common/token
replace github.com/ci4rail/kyt/kyt-server-common/iothub_wrapper => ../kyt-server-common/iothub_wrapper
