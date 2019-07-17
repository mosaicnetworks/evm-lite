package version

//Maj is Major Version Number
const Maj = "0"

//Min is Minor Version Number
const Min = "3"

//Fix is the Patch Version
const Fix = "0"

var (
	//Version contains the full version string
	Version = "0.3.0"

	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	// GitBranch is set with --ldflags "-X main.gitBranch=$(git symbolic-ref --short HEAD)"
	GitBranch string
)

func init() {
	// branch is only of interest if it is not the master branch
	if GitBranch != "" && GitBranch != "master" {
		Version += "-" + GitBranch
	}
	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}
}
