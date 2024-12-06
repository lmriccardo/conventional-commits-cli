package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Filesystem struct {
	Target  string `json:"target"`
	Source  string `json:"source"`
	Fstype  string `json:"fstype"`
	Options string `json:"options"`
}

type FindmntOutput struct {
	Filesystems []Filesystem `json:"filesystems"`
}

// Take all the environment variables and search for the one
// given as input to the function
func GetContainerEnvironmentVariable(varname string) (string, error) {
	variables := os.Environ()
	for _, element := range variables {
		parts := strings.Split(element, "=")
		if strings.Compare(parts[0], varname) == 0 {
			return parts[1], nil
		}
	}

	// If no variable has been found it returns an error
	return "", errors.New("no environment variable matches the provided key")
}

// Read the input cgroup file and returns the container id (in case)
func scanCgroupFile(filepath string) (string, error) {
	// Check that the file exists and can be opened
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err // Check for error and in case returns it
	}

	// First we need to check that the content of the file matches
	// what we expect to find in a typical docker container
	pattern := `^\d+:[a-zA-Z_-]*(=[a-zA-Z_-]+)?:/docker/[a-zA-Z0-9]+$`
	re := regexp.MustCompile(pattern) // Compile the pattern or error
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Check if the current line matches the pattern
		if re.MatchString(line) {
			docker_part := strings.Split(line, ":")[2]
			container_id := strings.Split(docker_part, "/")[2]
			return container_id, nil
		}
	}

	return "", errors.New("no container detected")
}

// Check if the current running environment is inside a docker container
func IsContainerEnvironment() (bool, string, error) {
	// There are some heuristic that can be used to identify
	// if the current running environment is a container or not.
	cgroups_file := "/proc/self/cgroup"
	dockerenv_file := "/.dockerenv"

	// Usually, if ccommits is started using the default command and
	// nothing is changed about cgroups, then the container will use
	// the default one created by docker.
	container_id, err := scanCgroupFile(cgroups_file)
	if err != nil {
		return false, "", err // Returns the error and false
	}

	// Finally, check if the typical .dockerenv file created
	// at the root of the container filesystem exists
	if _, err := os.Stat(dockerenv_file); os.IsNotExist(err) {
		return false, "", err
	}

	return true, container_id, nil
}

// Check if the input folder is mounted into the container
func IsContainerFolderMounted(folder string) *FindmntOutput {
	var findmnt_output bytes.Buffer // Initialize the output command buffer

	// Create the findmnt command, links the output and run
	findmnt := exec.Command("findmnt", "--target", folder, "--json")
	findmnt.Stdout = &findmnt_output
	err := findmnt.Run()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Unmarshal the output from findmnt into the structure
	output := new(FindmntOutput) // Create the output structure
	err = json.Unmarshal(findmnt_output.Bytes(), output)
	if err != nil {
		fmt.Printf("Error parsing the JSON: %s\n", err.Error())
		return nil
	}

	// Check if the target is the same as the input folder. This is a
	// necessary check, since findmnt with current folder returns always
	// something. For folders that are not bind mounted for example
	// there would be something referring to the root folder.
	if len(output.Filesystems) < 1 || output.Filesystems[0].Target != folder {
		return nil
	}

	return output
}
