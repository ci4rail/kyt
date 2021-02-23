module github.com/ci4rail/kyt/kyt-dlm-server

go 1.15

require (
	github.com/amenzhinsky/iothub v0.6.2
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/ci4rail/kyt/kyt-server-common/iothub_wrapper v0.0.0-00010101000000-000000000000
	github.com/ci4rail/kyt/kyt-server-common/token v0.0.0-00010101000000-000000000000
	github.com/ci4rail/kyt/kyt-server-common/version v0.0.0-00010101000000-000000000000 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golangci/golangci-lint v1.36.0
	github.com/lestrrat/go-jwx v0.0.0-20180221005942-b7d4802280ae
	github.com/lestrrat/go-pdebug v0.0.0-20180220043741-569c97477ae8 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/zalando/gin-oauth2 v1.5.2
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

replace github.com/ci4rail/kyt/kyt-server-common/token => ../kyt-server-common/token

replace github.com/ci4rail/kyt/kyt-server-common/iothub_wrapper => ../kyt-server-common/iothub_wrapper

replace github.com/ci4rail/kyt/kyt-server-common/version => ../kyt-server-common/version
