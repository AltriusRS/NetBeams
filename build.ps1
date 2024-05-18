Write-Output "Building NetBeams"
Write-Output "Collecting build information"

$version = $(git describe --tags --abbrev=0)
$maxProtocolVersion = "3.0.0"
$minProtocolVersion = "2.0.0"
$gitSha = $(git rev-parse --short HEAD)
$gitBranch = $(git rev-parse --abbrev-ref HEAD)
$buildTime = $(Get-Date -u +'%Y-%m-%dT%H:%M:%SZ')
$buildUser = $(whoami)
$goVersion = $(go env GOVERSION)
$buildOS = $(go env GOOS)
$buildArch = $(go env GOARCH)
$GO111MODULE=$(go env GO111MODULE)
$GOARCH=$(go env GOARCH)
$GOBIN=$(go env GOBIN)
$GOEXPERIMENT=$(go env GOEXPERIMENT)
$GOFLAGS=$(go env GOFLAGS)
$GOHOSTARCH=$(go env GOHOSTARCH)
$GOHOSTOS=$(go env GOHOSTOS)
$GOINSECURE=$(go env GOINSECURE)
$GONOPROXY=$(go env GONOPROXY)
$GONOSUMDB=$(go env GONOSUMDB)
$GOOS=$(go env GOOS)
$GOPRIVATE=$(go env GOPRIVATE)
$GOPROXY=$(go env GOPROXY)
$GOROOT=$(go env GOROOT)
$GOSUMDB=$(go env GOSUMDB)
$GOTOOLCHAIN=$(go env GOTOOLCHAIN)
$GOVCS=$(go env GOVCS)
$GOVERSION=$(go env GOVERSION)
$GCCGO=$(go env GCCGO)
$GOAMD64=$(go env GOAMD64)
$AR=$(go env AR)
$CC=$(go env CC)
$CXX=$(go env CXX)
$CGO_ENABLED=$(go env CGO_ENABLED)
$CGO_CFLAGS=$(go env CGO_CFLAGS)
$CGO_CPPFLAGS=$(go env CGO_CPPFLAGS)
$CGO_CXXFLAGS=$(go env CGO_CXXFLAGS)
$CGO_FFLAGS=$(go env CGO_FFLAGS)
$CGO_LDFLAGS=$(go env CGO_LDFLAGS)
$PKG_CONFIG=$(go env PKG_CONFIG)



Write-Output "Environment information"
Write-Output "Version: $version"
Write-Output "Git SHA: $gitSha"
Write-Output "Git branch: $gitBranch"
Write-Output "Build time: $buildTime"
Write-Output "Build user: $buildUser"
Write-Output "Go version: $goVersion"
Write-Output "Build OS: $buildOS"
Write-Output "Build architecture: $buildArch"
Write-Output "Maximum protocol version: $maxProtocolVersion"
Write-Output "Minimum protocol version: $minProtocolVersion"
Write-Output "Truncating go env variables"

$ldflags = ""
$ldflagsPrefix = "-X 'netbeams/environment"

$ldflags += " $ldflagsPrefix.Version=$version-$gitSha'"
$ldflags += " $ldflagsPrefix.MaxProtocolVersion=$maxProtocolVersion'"
$ldflags += " $ldflagsPrefix.MinProtocolVersion=$minProtocolVersion'"
$ldflags += " $ldflagsPrefix.GitSha=$gitSha'"
$ldflags += " $ldflagsPrefix.GitBranch=$gitBranch'"
$ldflags += " $ldflagsPrefix.BuildTime=$buildTime'"
$ldflags += " $ldflagsPrefix.BuildUser=$buildUser'"
$ldflags += " $ldflagsPrefix.GoVersion=$goVersion'"
$ldflags += " $ldflagsPrefix.BuildOS=$buildOS'"
$ldflags += " $ldflagsPrefix.BuildArch=$buildArch'"
$ldflags+=" $ldflagsPrefix.GO111MODULE=$GO111MODULE'"
$ldflags+=" $ldflagsPrefix.GOARCH=$GOARCH'"
$ldflags+=" $ldflagsPrefix.GOBIN=$GOBIN'"
$ldflags+=" $ldflagsPrefix.GOEXPERIMENT=$GOEXPERIMENT'"
$ldflags+=" $ldflagsPrefix.GOFLAGS=$GOFLAGS'"
$ldflags+=" $ldflagsPrefix.GOHOSTARCH=$GOHOSTARCH'"
$ldflags+=" $ldflagsPrefix.GOHOSTOS=$GOHOSTOS'"
$ldflags+=" $ldflagsPrefix.GOINSECURE=$GOINSECURE'"
$ldflags+=" $ldflagsPrefix.GOMODCACHE=$GOMODCACHE'"
$ldflags+=" $ldflagsPrefix.GONOPROXY=$GONOPROXY'"
$ldflags+=" $ldflagsPrefix.GONOSUMDB=$GONOSUMDB'"
$ldflags+=" $ldflagsPrefix.GOOS=$GOOS'"
$ldflags+=" $ldflagsPrefix.GOPRIVATE=$GOPRIVATE'"
$ldflags+=" $ldflagsPrefix.GOPROXY=$GOPROXY'"
$ldflags+=" $ldflagsPrefix.GOROOT=$GOROOT'"
$ldflags+=" $ldflagsPrefix.GOSUMDB=$GOSUMDB'"
$ldflags+=" $ldflagsPrefix.GOTOOLCHAIN=$GOTOOLCHAIN'"
$ldflags+=" $ldflagsPrefix.GOVCS=$GOVCS'"
$ldflags+=" $ldflagsPrefix.GOVERSION=$GOVERSION'"
$ldflags+=" $ldflagsPrefix.GCCGO=$GCCGO'"
$ldflags+=" $ldflagsPrefix.GOAMD64=$GOAMD64'"
$ldflags+=" $ldflagsPrefix.AR=$AR'"
$ldflags+=" $ldflagsPrefix.CC=$CC'"
$ldflags+=" $ldflagsPrefix.CXX=$CXX'"
$ldflags+=" $ldflagsPrefix.CGO_ENABLED=$CGO_ENABLED'"
$ldflags+=" $ldflagsPrefix.CGO_CFLAGS=$CGO_CFLAGS'"
$ldflags+=" $ldflagsPrefix.CGO_CPPFLAGS=$CGO_CPPFLAGS'"
$ldflags+=" $ldflagsPrefix.CGO_CXXFLAGS=$CGO_CXXFLAGS'"
$ldflags+=" $ldflagsPrefix.CGO_FFLAGS=$CGO_FFLAGS'"
$ldflags+=" $ldflagsPrefix.CGO_LDFLAGS=$CGO_LDFLAGS'"
$ldflags+=" $ldflagsPrefix.PKG_CONFIG=$PKG_CONFIG'"

# Write-Output "LDFLAGS: $ldflags"

Write-Output "Starting build"

go build -ldflags="-s -w $ldflags" -v -o ./artifacts/netbeams.exe . 

Write-Output "Build complete"

Write-Output "Artifact location: ./artifacts/netbeams.exe"