package aruku

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
)

const (
	testDataDir = "testdata"
)

func TestLoad(t *testing.T) {

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

	a.Write(getTestDataDir())

	err := a.Load(getTestDataDir())

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	copyTestFile(path.Join(getTestDataDir(), "aruku.yaml"), path.Join(getTestDataDir(), "aruku-bak.yaml"))
}

func TestLoadEmptyCmdList(t *testing.T) {
	script := `
	{
		"Author": "Taro Fukunaga",
		"Description": "Install myapp",
		"CmdList": []
	}
	`

	writeTestFile(path.Join(testDataDir, "aruku.yaml"), script)

	var a App

	err := a.Load(getTestDataDir())

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}

func TestLoadInvalidYaml(t *testing.T) {
	script := `
	{
		"Author", "Taro Fukunaga",
		"Description", "Install myapp",
		"CmdList", []
	}
	`

	writeTestFile(path.Join(testDataDir, "aruku.yaml"), script)

	var a App

	err := a.Load(getTestDataDir())

	if err != nil {
		fmt.Println(err)
	}
}

func getTestDataDir() string {
	wd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}
	return filepath.Join(wd, "testdata")
}

func copyTestFile(src, dest string) int64 {
	inFile, _ := os.Open(src)

	defer inFile.Close()

	outFile, _ := os.Create(dest)

	defer outFile.Close()

	nBytes, _ := io.Copy(outFile, inFile)

	return nBytes
}

func writeTestFile(path, script string) {

	f, err := os.Create(path)

	if err != nil {
		log.Fatalf("Unable to write file")
	}

	defer f.Close()

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
