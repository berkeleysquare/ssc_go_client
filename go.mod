module github.com/SpectraLogic/ssc_go_client

go 1.17

replace github.com/SpectraLogic/ssc_go_client/openapi => ./openapi/

require (
	github.com/SpectraLogic/ds3_go_sdk v5.4.0+incompatible
	github.com/SpectraLogic/ssc_go_client/openapi v0.0.0-00010101000000-000000000000
	github.com/antihax/optional v1.0.0
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/aws/aws-sdk-go v1.26.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
