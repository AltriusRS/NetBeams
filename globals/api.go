package globals

import "fmt"

const BaseAuthAPIURL = "https://auth.beammp.com"
const BaseAPIURL = "https://backend.beammp.com"

var UserAgent = fmt.Sprintf("NetBeams/%s", Version)
