package git

import (
	"os"
	"os/exec"
	"strings"
)

type FileChange struct {
	Action   string
	FileName string
	Diff     string
}

func GetChanges() ([]FileChange, error) {
	// 获取修改的文件列表
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var changes []FileChange
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) < 4 {
			continue
		}
		action := line[:1]

		fileName := strings.TrimSpace(line[3:])
		var diff string
		if action != "D" {
			diff, err = getFileDiff(fileName)
			if err != nil {
				return nil, err
			}
		}
		switch action {
		case "M":
			action = "modified"
		case "A":
			action = "added"
		case "D":
			action = "deleted"
		}
		changes = append(changes, FileChange{
			FileName: fileName,
			Action:   action,
			Diff:     diff,
		})
	}

	return changes, nil
}

func getFileDiff(fileName string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached", fileName)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func AddTrackedFiles() error {
	cmd := exec.Command("git", "add", "-u")
	return cmd.Run()
}
