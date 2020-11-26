package aruku

import (
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
				Args:             []string{"a", "b", "c"},
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
