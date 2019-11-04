package version

// Flag contains extra info about the version. It is helpul for tracking
// versions while developing. It should always be empty on the master branch.
// This will be enforced in a continuous integration test.
const Flag = "develop"

var (
	//Version contains the full version string
	Version = "0.3.5"
	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	// GitBranch is set with --ldflags "-X main.gitBranch=$(git symbolic-ref --short HEAD)"
	GitBranch string

	//JSONVersion is set from the run command explicitly.
	//This allows it to be set differently within
	//monetd
	JSONVersion map[string]string
)

func init() {
	Version += "-" + Flag

	// branch is only of interest if it is not the master branch
	if GitBranch != "" && GitBranch != "master" {
		Version += "-" + GitBranch
	}
	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}

	JSONVersion = make(map[string]string)
	JSONVersion["evm-lite"] = Version
}
