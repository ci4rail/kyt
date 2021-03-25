module github.com/ci4rail/kyt/kyt-dlm-server

go 1.15

require (
	github.com/auth0/go-jwt-middleware v1.0.0 // indirect
	github.com/ci4rail/kyt/kyt-server-common/iothubwrapper v0.0.0-00010101000000-000000000000
	github.com/ci4rail/kyt/kyt-server-common/token v0.0.0-00010101000000-000000000000
	github.com/ci4rail/kyt/kyt-server-common/version v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/urfave/negroni v1.0.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/ci4rail/kyt/kyt-server-common/token => ../kyt-server-common/token

replace github.com/ci4rail/kyt/kyt-server-common/iothubwrapper => ../kyt-server-common/iothubwrapper

replace github.com/ci4rail/kyt/kyt-server-common/version => ../kyt-server-common/version
