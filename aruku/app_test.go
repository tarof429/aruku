package aruku

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRun(t *testing.T) {

	createTestFootprint()

	var a App
	a.Author = "Taro Fukunaga"
	a.Description = "Install myapp"

	// Install
	var cmds CmdList
	cmds.Description = "Install"

	cmd := Command{
		Name:             "w",
		CommandType:      ExecuteCommandType,
		Args:             []string{"-h"},
		Description:      "Show who's logged in",
		WorkingDirectory: getTestDataDir()}
	cmds.Cmds = append(cmds.Cmds, cmd)

	cmd = Command{

		Name:             "python",
		CommandType:      ExecuteCommandType,
		Args:             []string{"test.py"},
		Description:      "Print platform",
		WorkingDirectory: getTestDataDir()}
	cmds.Cmds = append(cmds.Cmds, cmd)

	cmd = Command{
		Name:         "docker username",
		CommandType:  ReadCommandType,
		Description:  "Enter docker username",
		VariableName: "DOCKER_USERNAME",
	}
	cmds.Cmds = append(cmds.Cmds, cmd)

	a.CmdList = append(a.CmdList, cmds)

	a.Write("testdata")

	writeTestPythonScript("testdata/test.py")

	if a.Author != "Taro Fukunaga" {
		t.Fatalf("Invalid author")
	}

	if a.Description != "Install myapp" {
		t.Fatalf("Invalid description")
	}

	if len(a.CmdList) != 1 {
		t.Fatalf("Did not find 1 command list")
	}

	fmt.Println("Doing install")
	a.SetCmdList("Install")
	for a.HasNextCmd() {
		cmd := a.GetCurrentCmd()
		cmd.Print()

		a.RunCurrentCmd()
		cmd = a.GetCurrentCmd()

		fmt.Printf("Result:\n%v\n", cmd.GetOutput())
		a.PointToNextCmd()
	}

}

func getTestDataDir() string {
	wd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}
	return filepath.Join(wd, "testdata")
}

func writeTestPythonScript(path string) {

	fmt.Println("Writing script to file")

	f, err := os.Create(path)

	if err != nil {
		log.Fatalf("Unable to create test script")
	}

	defer f.Close()

	script := `#!/usr/bin/python
import platform
p = platform.platform()
print(p)`
	n, err := f.WriteString(script)

	if err != nil {
		log.Fatalf("Unable to write string")
	}

	fmt.Printf("Wrote :%v bytes\n", strconv.Itoa(n))
}

func createTestFootprint() {
	wd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	testdataDir := filepath.Join(wd, "testdata")

	_, err = os.Stat(testdataDir)

	if os.IsNotExist(err) {
		mode := int(0755)
		err := os.Mkdir(testdataDir, os.FileMode(mode))

		if err != nil {
			log.Fatalf("Error: unable to create testdata: %v\n", err)
		}
	}
}
