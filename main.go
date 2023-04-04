package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ncruces/zenity"
	"golang.org/x/sys/windows/registry"
)

var regKeyName = `software\Go-Custom-Scr-Command`

func main() {
	for _, item := range os.Args {
		println(item)
	}
	if len(os.Args) < 2 {
		return
	}

	command := os.Args[1]

	if strings.Contains(command, "/c") {
		getUserInput()
	}

	if strings.Contains(command, "/S") || strings.Contains(command, "/s") {
		runExecutable()
	}

}

func runExecutable() {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, regKeyName, registry.ALL_ACCESS)
	ChkErr(err)
	defer key.Close()

	filePath, _, _ := key.GetStringValue("TargetPath")
	args, _, _ := key.GetStringValue("Args")
	if filePath == "" {
		return
	} else {
		cmd := exec.Command(filePath, args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		ChkErr(err)
	}

}

func getUserInput() {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, regKeyName, registry.ALL_ACCESS)
	ChkErr(err)
	defer key.Close()

	targetPath, _, _ := key.GetStringValue("TargetPath")
	if targetPath == "" {
		err = key.SetStringValue("TargetPath", `C:\`)
		ChkErr(err)
		targetPath = `C:\`
	}

	selectExecFile := "Select executable"
	enterCmdArgs := "Enter command line arguments"

	actions := []string{selectExecFile, enterCmdArgs}
	selectedAction, err := zenity.List("Select Action", actions, zenity.Title("Custom screen saver settings"))
	ChkErr(err)

	switch selectedAction {
	case selectExecFile:
		targetPath, err = zenity.SelectFile(zenity.Filename(targetPath), zenity.Title("Select Executable"))
		ChkErr(err)
		err = key.SetStringValue("TargetPath", targetPath)
		ChkErr(err)

	case enterCmdArgs:
		args, _, _ := key.GetStringValue("Args")
		args, err := zenity.Entry("Enter command line arguments", zenity.Title("Command Line Arguments"), zenity.EntryText(args))
		ChkErr(err)
		err = key.SetStringValue("Args", args)
		ChkErr(err)
	}

}

func ChkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ChkRegistryValueFound(err error) bool {
	if err != nil {
		if pe, ok := err.(*os.PathError); ok {
			if errno, ok := pe.Err.(syscall.Errno); ok {
				if errno == syscall.ERROR_FILE_NOT_FOUND {
					return false
				}
			}
		} else {
			return false
		}
		return false
	}
	return true
}
