package main

import (
	"bufio"
	"os"
)

// Main Entrypoint
func main() {

	// Create the configuration
	config := CreateConfig()

	// Call the backup function
	BackupFactory(config)
}

// This function will check for configurations being passed in
// via arguments or environment variables and creates an instance of
// Config from the detected values
func CreateConfig() Config {

	// Create empty variables
	var _Directories []string
	var _TargetBucket string
	var _AwsProfile string

	if len(os.Args) > 1 {
		_TargetBucket = os.Args[1]
	} else {
		_TargetBucket = os.Getenv("TARGET_BUCKET")
	}

	if len(_TargetBucket) == 0 {
		exitErrorf("Target Bucket cannot be nil")
	}

	if len(os.Args) > 2 {
		_Directories = GetDirectories(os.Args[2])
	} else {
		_Directories = GetDirectories(os.Getenv("DIRECTORY_LIST"))
	}

	if len(_Directories) == 0 {
		exitErrorf("Directories list cannot be nil")
	}

	if len(os.Args) > 3 {
		_AwsProfile = os.Args[3]
	} else {
		_AwsProfile = os.Getenv("AWS_PROFILE")
	}

	config := Config{
		TargetBucket: _TargetBucket,
		Directories:  _Directories,
		AwsProfile:   _AwsProfile,
	}

	return config
}

// This function returns a slice of strings containing the
// contents of the items.lst file (file containing a list )
// of directories / file names we want to back up
//
// **NOTE**
// The method of passing in data to be backed up will be
// improved in future versions
func GetDirectories(Filename string) []string {
	file, err := os.Open(Filename)
	if err != nil {
		exitErrorf("Unable to read directory list %v", err)
	}

	defer file.Close()

	var directories []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		directories = append(directories, scanner.Text())
	}

	if len(directories) != 0 {
		return directories
	} else {
		exitErrorf("Directories list cannot be empty")
	}
	return nil
}
