package aruku

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReadWrite(t *testing.T) {

	createTestFootprint()

	var a App
	a.Author = "Taro Fukunaga"
	a.Description = "Install myapp"

	addCmdList := func(description string) {
		var cmds CmdList

		cmds.Description = description

		for i := 0; i < 3; i++ {
			cmd := Command{
				Name:             "ls",
				Args:             []string{"-l", "."},
				Description:      "Simple command " + strconv.Itoa(i),
				WorkingDirectory: "testdata",
			}
			cmds.Cmds = append(cmds.Cmds, cmd)
		}
		a.CmdList = append(a.CmdList, cmds)
	}

	addCmdList("Install")
	addCmdList("Uninstall")
	addCmdList("Upgrade")

	err := a.Write("testdata/data")

	if err != nil {
		t.Fatalf("Unable to write data")
	}

	err = a.Read("testdata/data")

	if err != nil {
		t.Fatalf("Unable to read data")
	}

	if len(a.CmdList) != 3 {
		t.Fatalf("Did not read commmands successfully")
	}

	a.Export("testdata", "export/test.zip")

	a.Import("export/test.zip", "import")

	if a.Author != "Taro Fukunaga" {
		t.Fatalf("Invalid author")
	}

	if a.Description != "Install myapp" {
		t.Fatalf("Invalid description")
	}

	if len(a.CmdList) != 3 {
		t.Fatalf("Did not find 3 commands")
	}

	fmt.Println("Doing install")
	a.SetCmdList("Install")
	for a.HasNextCmd() {
		cmd := a.GetCurrentCmd()
		fmt.Printf("Printing command: %v\n", cmd)

		a.RunCurrentCmd()
		cmd = a.GetCurrentCmd()

		fmt.Printf("Result: %v\n", cmd)
		a.PointToNextCmd()
	}

	fmt.Println("Doing upgrade")
	a.SetCmdList("Upgrade")
	for a.HasNextCmd() {
		cmd := a.GetCurrentCmd()
		fmt.Printf("Printing command: %v\n", cmd)

		a.RunCurrentCmd()
		cmd = a.GetCurrentCmd()

		fmt.Printf("Result: %v\n", cmd)
		a.PointToNextCmd()
	}

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
