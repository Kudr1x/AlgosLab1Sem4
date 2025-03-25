package util

import (
	"os"
	"os/exec"
)

func RunPython(jsonFilePath string) {
	cmd := exec.Command("python", "/home/kudrix/GolandProjects/AlgosLab1Sem4/scripts/graph.py", jsonFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
