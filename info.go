package fairway

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

var gitInfo = make(map[string]string)

func loadGitInfo() {

	file, err := os.Open("git.properties")

	if err != nil {
		// can't open file or not exists
		out, err := exec.Command("bash", "-c", "git rev-parse --abbrev-ref HEAD; git show -s --format=\"%h%n%ci\"").Output()
		if err != nil {
			logger.Error(err)
		}
		splitted := strings.Split(string(out), "\n")
		if len(splitted) > 3 {
			gitInfo["branch"] = splitted[0]
			gitInfo["commitid"] = splitted[1]
			gitInfo["committime"] = splitted[2]
		}

	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		gitInfo["branch"] = scanner.Text()
		scanner.Scan()
		gitInfo["commitid"] = scanner.Text()
		scanner.Scan()
		gitInfo["committime"] = scanner.Text()
	}

}

func info() ([]byte, error) {
	info := generateInfoData()
	return toJSON(info), nil
}

func generateInfoData() *infoData {
	infoJSON := new(infoData)
	infoJSON.Git = new(git)
	infoJSON.Git.Commit = new(commit)

	infoJSON.Git.Branch = gitInfo["branch"]
	infoJSON.Git.Commit.ID = gitInfo["commitid"]
	infoJSON.Git.Commit.Time = gitInfo["committime"]

	return infoJSON
}

type infoData struct {
	Git *git `json:"git"`
}

type git struct {
	Branch string  `json:"branch"`
	Commit *commit `json:"commit"`
}

type commit struct {
	ID   string `json:"id"`
	Time string `json:"time"`
}
