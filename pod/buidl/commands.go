package main

var commands = map[string][]string{
	"build": {
		"go generate ./version/.",
		"go build  ./pod/.",
	},
	"install": {
		"go generate ./version/.",
		"go install  ./pod/.",
	},
	"justinstall": {
		"go install  ./pod/.",
	},
	"release": {
		"go generate ./version/.",
		"go install  -ldflags '-w -s' ./pod/.",
	},
	"gui": {
		"go generate ./version/.",
		"go run  ./pod/. gui",
	},
	"node": {
		"go generate ./version/.",
		"go run  ./pod/. node",
	},
	"wallet": {
		"go generate ./version/.",
		"go run  ./pod/.",
	},
	"kopach": {
		"go generate ./version/.",
		"go run  ./pod/.",
	},
	"headless": {
		"go generate ./version/.",
		"go install  -tags headless ./pod/.",
	},
	"docker": {
		"go generate ./version/.",
		"go install  -tags headless ./pod/.",
	},
	"appstores": {
		"go generate ./version/.",
		"go install  -tags nominers ./pod/.",
	},
	"tests": {
		"go generate ./version/.",
		"go test ./...",
	},
	"builder": {
		"go generate ./version/.",
		"go install  ./pod/buidl/.",
	},
	"generate": {
		"go generate ./...",
	},
	"patch": {
		"go run ./version/update/. patch",
	},
}
