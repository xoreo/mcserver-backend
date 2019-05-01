package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/xoreo/mcserver-backend/common"
	"github.com/xoreo/mcserver-backend/types"
)

var (
	// ErrUnsupportedVersion is thrown when an unsupported server version is given.
	ErrUnsupportedVersion = errors.New("that is not a supported version")

	// ErrServerHasNotBeenInitialized is thrown when a server's metadata exists but the server has not actually been initialized on the local machine.
	ErrServerHasNotBeenInitialized = errors.New("that server has not actually been initialized yet. Initialize it with InitializeServer()")
)

// GenerateStartScript generates the script that launches the server.
func GenerateStartScript(server types.Server) []byte {
	ramstr := strconv.Itoa(server.RAM) // Convert the ram to a string

	path := filepath.Join(server.Path, server.Version)
	script := `#!/bin/bash
cd ` + path + `
java -Xms` + ramstr + `M -Xmx` + ramstr + `M -jar ` + path + `/` + server.Version + `.jar nogui`
	return []byte(script)
}

// InitializeServer initializes a new server onto the local machine.
func InitializeServer(server *types.Server) error {
	var url string
	dServer := *server // Make a copy of the pointer (dereference for convenience)

	// Determine the pre-made server download url
	switch dServer.Version {
	case "1.12":
		url = common.ServerV112
	case "1.8":
		url = common.ServerV18
	case "1.7.2":
		url = common.ServerV172
	case "1.2.1":
		url = common.ServerV121
	default:
		return ErrUnsupportedVersion
	}

	// Download the pre-made server
	zipPath, err := common.DownloadServer(url, dServer.Path, dServer.Version)
	if err != nil {
		return err
	}

	// Unzip the downloaded file
	_, err = common.Unzip(zipPath, dServer.Path)
	if err != nil {
		return err
	}

	// Create start script for the server
	workingPath := filepath.Join(dServer.Path, dServer.Version)
	startScriptPath := filepath.Join(workingPath, "start.sh")
	script := GenerateStartScript(dServer)

	// Install the script
	err = ioutil.WriteFile(startScriptPath, script, 0644)
	if err != nil {
		return err
	}

	server.Initialized = true // Set the server's initialized state to true
	server.StartScript = startScriptPath
	return nil
}

// StartServer starts a server.
func StartServer(server *types.Server) error {
	dServer := *server // Dereference for convenience

	// Make sure that the server has been initialized.
	if !dServer.Initialized || dServer.StartScript == "" {
		return ErrServerHasNotBeenInitialized
	}

	launcher := filepath.Join(dServer.Path, dServer.Version, "start.sh")
	cmd := exec.Command("/bin/sh", launcher)
	fmt.Println(cmd.Output())

	return nil
}

// RestartServer restarts a server.
func RestartServer(server types.Server) {
	exec.Command("/bin/sh", "systemctl restart "+server.Name)
	status := exec.Command("/bin/sh", "systemctl status "+server.Name)
	fmt.Println(status.Output())
}

// StopServer stops a server.
func StopServer(server types.Server) {
	exec.Command("/bin/sh", "systemctl stop "+server.Name)
	status := exec.Command("/bin/sh", "systemctl status "+server.Name)
	fmt.Println(status.Output())
}

// EnterServer launches a shell of the server console.
func EnterServer(server *types.Server) {

}

// EditProperties is used to edit a server property (such as max build height or default gamemode).
func EditProperties(property, newValue string) error {
	return nil
}
