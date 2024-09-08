package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Servers struct {
	Servers map[string]interface{} `json:"servers"`
}

func main() {
	// Read the JSON file
	file, err := os.ReadFile("servers.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var servers Servers
	if err := json.Unmarshal(file, &servers); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Create a temporary file to store the server list
	tmpfile, err := os.CreateTemp("", "fzf-ssh")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	// Write servers to the temporary file
	writeServers(tmpfile, servers.Servers)

	tmpfile.Close()

	// Run fzf to select a server
	cmd := exec.Command("fzf", "--height", "10", "--reverse")
	cmd.Stdin, _ = os.Open(tmpfile.Name())
	selectedServer, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running fzf:", err)
		return
	}

	// Trim the selected server
	server := strings.TrimSpace(string(selectedServer))

	// List available SSH keys
	sshKeyFiles, err := os.ReadDir("authorized_keys/")
	if err != nil {
		fmt.Println("Error reading SSH keys directory:", err)
		return
	}

	// Create a temporary file to store the list of SSH keys
	keyTmpfile, err := os.CreateTemp("", "fzf-ssh-keys")
	if err != nil {
		fmt.Println("Error creating temp file for SSH keys:", err)
		return
	}
	defer os.Remove(keyTmpfile.Name())

	// Write SSH keys to the temporary file
	for _, file := range sshKeyFiles {
		if !file.IsDir() && !strings.HasSuffix(file.Name(), ".pub") && !strings.HasPrefix(file.Name(), "known_hosts") {
			keyTmpfile.WriteString(file.Name() + "\n")
		}
	}
	keyTmpfile.Close()

	// Run fzf to select an SSH key
	keyCmd := exec.Command("fzf", "--height", "10", "--reverse")
	keyCmd.Stdin, _ = os.Open(keyTmpfile.Name())
	selectedKey, err := keyCmd.Output()
	if err != nil {
		fmt.Println("Error running fzf for SSH keys:", err)
		return
	}

	// Trim the selected key
	sshKeyFile := strings.TrimSpace(string(selectedKey))

	// Split the server string into hostname and port
	serverParts := strings.Split(server, " ")
	if len(serverParts) < 3 {
		fmt.Println("Invalid server format")
		return
	}
	hostname := serverParts[1]
	port := serverParts[2]

	// Remove old host key
	removeHostKeyCmd := exec.Command("ssh-keygen", "-R", fmt.Sprintf("[%s]:%s", hostname, port))
	removeHostKeyCmd.Run()

	// Connect to the selected server using the selected SSH key file
	sshCmd := exec.Command("ssh", "-i", "authorized_keys/"+sshKeyFile, "-p", port, "user@"+hostname)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	sshCmd.Stdin = os.Stdin

	if err := sshCmd.Run(); err != nil {
		fmt.Println("Error connecting to server:", err)
	}
}

func writeServers(file *os.File, servers map[string]interface{}) {
	for _, value := range servers {
		switch v := value.(type) {
		case map[string]interface{}:
			writeServers(file, v)
		case []interface{}:
			for _, server := range v {
				file.WriteString(server.(string) + "\n")
			}
		}
	}
}
