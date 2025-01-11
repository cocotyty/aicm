package git

import (
	"os/exec"
	"strings"
)

type FileChange struct {
	FileName string
	Diff     string
}

func GetChanges() ([]FileChange, error) {
	// 获取修改的文件列表
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	var changes []FileChange
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) < 4 {
			continue
		}
		fileName := strings.TrimSpace(line[3:])
		diff, err := getFileDiff(fileName)
		if err != nil {
			return nil, err
		}
		changes = append(changes, FileChange{
			FileName: fileName,
			Diff:     diff,
		})
	}

	return changes, nil
}

func getFileDiff(fileName string) (string, error) {
	out, err := exec.Command("git", "diff", "--cached", fileName).Output()
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
