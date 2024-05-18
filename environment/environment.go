package environment

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type BuildContext struct {
	Version                string
	GitSha                 string
	GitBranch              string
	BuildTime              string
	BuildUser              string
	MaxProtocolVersion     string
	MinProtocolVersion     string
	SemverMaxClientVersion *semver.Constraints
	SemverMinClientVersion *semver.Constraints
	IsDev                  bool
	GoVersion              string
	BuildOS                string
	BuildArch              string
	GO111MODULE            string
	GOARCH                 string
	GOBIN                  string
	GOCACHE                string
	GOENV                  string
	GOEXE                  string
	GOEXPERIMENT           string
	GOFLAGS                string
	GOHOSTARCH             string
	GOHOSTOS               string
	GOINSECURE             string
	GOMODCACHE             string
	GONOPROXY              string
	GONOSUMDB              string
	GOOS                   string
	GOPATH                 string
	GOPRIVATE              string
	GOPROXY                string
	GOROOT                 string
	GOSUMDB                string
	GOTMPDIR               string
	GOTOOLCHAIN            string
	GOTOOLDIR              string
	GOVCS                  string
	GOVERSION              string
	GCCGO                  string
	GOAMD64                string
	AR                     string
	CC                     string
	CXX                    string
	CGO_ENABLED            string
	GOMOD                  string
	GOWORK                 string
	CGO_CFLAGS             string
	CGO_CPPFLAGS           string
	CGO_CXXFLAGS           string
	CGO_FFLAGS             string
	CGO_LDFLAGS            string
	PKG_CONFIG             string
	GOGCCFLAGS             string
}

var Version = "UNSET"            // Version of the application
var GitSha = "UNSET"             // Git SHA of the application
var GitBranch = "UNSET"          // Git branch of the application
var BuildTime = "UNSET"          // Build time of the application
var BuildUser = "UNSET"          // Build user of the application
var MaxProtocolVersion = "3.0.0" // Maximum protocol version supported by the application
var MinProtocolVersion = "2.0.0" // Minimum protocol version supported by the application
var IsDev = "true"               // is the environment a development environment
var GoVersion = "UNSET"          // Go version of the application
var BuildOS = "UNSET"            // Build OS of the application
var BuildArch = "UNSET"          // Build architecture of the application
var GO111MODULE = "UNSET"        // Go Environment Variable
var GOARCH = "UNSET"             // Go Environment Variable
var GOBIN = "UNSET"              // Go Environment Variable
var GOCACHE = "UNSET"            // Go Environment Variable
var GOENV = "UNSET"              // Go Environment Variable
var GOEXE = "UNSET"              // Go Environment Variable
var GOEXPERIMENT = "UNSET"       // Go Environment Variable
var GOFLAGS = "UNSET"            // Go Environment Variable
var GOHOSTARCH = "UNSET"         // Go Environment Variable
var GOHOSTOS = "UNSET"           // Go Environment Variable
var GOINSECURE = "UNSET"         // Go Environment Variable
var GOMODCACHE = "UNSET"         // Go Environment Variable
var GONOPROXY = "UNSET"          // Go Environment Variable
var GONOSUMDB = "UNSET"          // Go Environment Variable
var GOOS = "UNSET"               // Go Environment Variable
var GOPATH = "UNSET"             // Go Environment Variable
var GOPRIVATE = "UNSET"          // Go Environment Variable
var GOPROXY = "UNSET"            // Go Environment Variable
var GOROOT = "UNSET"             // Go Environment Variable
var GOSUMDB = "UNSET"            // Go Environment Variable
var GOTMPDIR = "UNSET"           // Go Environment Variable
var GOTOOLCHAIN = "UNSET"        // Go Environment Variable
var GOTOOLDIR = "UNSET"          // Go Environment Variable
var GOVCS = "UNSET"              // Go Environment Variable
var GOVERSION = "UNSET"          // Go Environment Variable
var GCCGO = "UNSET"              // Go Environment Variable
var GOAMD64 = "UNSET"            // Go Environment Variable
var AR = "UNSET"                 // Go Environment Variable
var CC = "UNSET"                 // Go Environment Variable
var CXX = "UNSET"                // Go Environment Variable
var CGO_ENABLED = "UNSET"        // Go Environment Variable
var GOMOD = "UNSET"              // Go Environment Variable
var GOWORK = "UNSET"             // Go Environment Variable
var CGO_CFLAGS = "UNSET"         // Go Environment Variable
var CGO_CPPFLAGS = "UNSET"       // Go Environment Variable
var CGO_CXXFLAGS = "UNSET"       // Go Environment Variable
var CGO_FFLAGS = "UNSET"         // Go Environment Variable
var CGO_LDFLAGS = "UNSET"        // Go Environment Variable
var PKG_CONFIG = "UNSET"         // Go Environment Variable
var GOGCCFLAGS = "UNSET"         // Go Environment Variable

var Context = BuildContext{
	Version:                "UNSET",
	GitSha:                 "UNSET",
	GitBranch:              "UNSET",
	BuildTime:              "UNSET",
	BuildUser:              "UNSET",
	MaxProtocolVersion:     "UNSET",
	MinProtocolVersion:     "UNSET",
	SemverMaxClientVersion: nil,
	SemverMinClientVersion: nil,
	IsDev:                  true,
	GoVersion:              "UNSET",
	BuildOS:                "UNSET",
	BuildArch:              "UNSET",
	GO111MODULE:            "UNSET",
	GOARCH:                 "UNSET",
	GOBIN:                  "UNSET",
	GOCACHE:                "UNSET",
	GOENV:                  "UNSET",
	GOEXE:                  "UNSET",
	GOEXPERIMENT:           "UNSET",
	GOFLAGS:                "UNSET",
	GOHOSTARCH:             "UNSET",
	GOHOSTOS:               "UNSET",
	GOINSECURE:             "UNSET",
	GOMODCACHE:             "UNSET",
	GONOPROXY:              "UNSET",
	GONOSUMDB:              "UNSET",
	GOOS:                   "UNSET",
	GOPATH:                 "UNSET",
	GOPRIVATE:              "UNSET",
	GOPROXY:                "UNSET",
	GOROOT:                 "UNSET",
	GOSUMDB:                "UNSET",
	GOTMPDIR:               "UNSET",
	GOTOOLCHAIN:            "UNSET",
	GOTOOLDIR:              "UNSET",
	GOVCS:                  "UNSET",
	GOVERSION:              "UNSET",
	GCCGO:                  "UNSET",
	GOAMD64:                "UNSET",
	AR:                     "UNSET",
	CC:                     "UNSET",
	CXX:                    "UNSET",
	CGO_ENABLED:            "UNSET",
	GOMOD:                  "UNSET",
	GOWORK:                 "UNSET",
	CGO_CFLAGS:             "UNSET",
	CGO_CPPFLAGS:           "UNSET",
	CGO_CXXFLAGS:           "UNSET",
	CGO_FFLAGS:             "UNSET",
	CGO_LDFLAGS:            "UNSET",
	PKG_CONFIG:             "UNSET",
	GOGCCFLAGS:             "UNSET",
}

func GetBuildContext() {

	fmt.Printf("isDev: %s\n", IsDev)
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("GitSha: %s\n", GitSha)
	fmt.Printf("GitBranch: %s\n", GitBranch)
	fmt.Printf("BuildTime: %s\n", BuildTime)
	fmt.Printf("BuildUser: %s\n", BuildUser)

	SemverMaxClientVersion, _ := semver.NewConstraint("< " + MaxProtocolVersion)
	SemverMinClientVersion, _ := semver.NewConstraint(">= " + MinProtocolVersion)
	Context = BuildContext{
		Version:                Version,
		GitSha:                 GitSha,
		GitBranch:              GitBranch,
		BuildTime:              BuildTime,
		BuildUser:              BuildUser,
		MaxProtocolVersion:     MaxProtocolVersion,
		MinProtocolVersion:     MinProtocolVersion,
		SemverMaxClientVersion: SemverMaxClientVersion,
		SemverMinClientVersion: SemverMinClientVersion,
		IsDev:                  IsDev == "true",
		BuildOS:                BuildOS,
		BuildArch:              BuildArch,
		GoVersion:              GoVersion,
		GO111MODULE:            GO111MODULE,
		GOARCH:                 GOARCH,
		GOBIN:                  GOBIN,
		GOCACHE:                GOCACHE,
		GOENV:                  GOENV,
		GOEXE:                  GOEXE,
		GOEXPERIMENT:           GOEXPERIMENT,
		GOFLAGS:                GOFLAGS,
		GOHOSTARCH:             GOHOSTARCH,
		GOHOSTOS:               GOHOSTOS,
		GOINSECURE:             GOINSECURE,
		GOMODCACHE:             GOMODCACHE,
		GONOPROXY:              GONOPROXY,
		GONOSUMDB:              GONOSUMDB,
		GOOS:                   GOOS,
		GOPATH:                 GOPATH,
		GOPRIVATE:              GOPRIVATE,
		GOPROXY:                GOPROXY,
		GOROOT:                 GOROOT,
		GOSUMDB:                GOSUMDB,
		GOTMPDIR:               GOTMPDIR,
		GOTOOLCHAIN:            GOTOOLCHAIN,
		GOTOOLDIR:              GOTOOLDIR,
		GOVCS:                  GOVCS,
		GOVERSION:              GOVERSION,
		GCCGO:                  GCCGO,
		GOAMD64:                GOAMD64,
		AR:                     AR,
		CC:                     CC,
		CXX:                    CXX,
		CGO_ENABLED:            CGO_ENABLED,
		GOMOD:                  GOMOD,
		GOWORK:                 GOWORK,
		CGO_CFLAGS:             CGO_CFLAGS,
		CGO_CPPFLAGS:           CGO_CPPFLAGS,
		CGO_CXXFLAGS:           CGO_CXXFLAGS,
		CGO_FFLAGS:             CGO_FFLAGS,
		CGO_LDFLAGS:            CGO_LDFLAGS,
		PKG_CONFIG:             PKG_CONFIG,
		GOGCCFLAGS:             GOGCCFLAGS,
	}
}
