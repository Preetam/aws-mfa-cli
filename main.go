package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execWrap(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}
	return string(out), nil
}

func getUser(profile string) (string, error) {
	out, err := execWrap("aws",
		"--profile", profile,
		"iam", "get-user",
		"--output", "text")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(out, "\t")[6]), nil
}

func getMFADevice(profile, user string) (string, error) {
	out, err := execWrap("aws",
		"--profile", profile,
		"iam", "list-mfa-devices",
		"--user-name", user,
		"--output", "text")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(out, "\t")[2]), nil
}

func getCredentials(profile, serialNumber, token string) (string, string, string, error) {
	out, err := execWrap("aws",
		"--profile", profile,
		"sts", "get-session-token",
		"--serial-number", serialNumber, "--token-code", token,
		"--output", "text")
	if err != nil {
		return "", "", "", err
	}
	parts := strings.Split(out, "\t")
	accessKeyID := parts[1]
	secretAccessKey := parts[3]
	sessionToken := parts[4]
	return accessKeyID, secretAccessKey, sessionToken, nil
}

func configureCredentials(profile, region, accessKeyID, secretAccessKey, sessionToken string) error {
	opts := [][2]string{
		{"region", region},
		{"aws_access_key_id", accessKeyID},
		{"aws_secret_access_key", secretAccessKey},
		{"aws_session_token", sessionToken},
	}

	for _, opt := range opts {
		_, err := execWrap("aws",
			"--profile", profile,
			"configure", "set",
			opt[0], opt[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	profile := flag.String("profile", "default", "Profile name")
	profileMFA := flag.String("profile-mfa", "profile-mfa", "Temporary MFA profile name")
	token := flag.String("token", "", "MFA token")
	region := flag.String("region", "us-east-1", "AWS Region")
	flag.Parse()

	if *token == "" {
		fmt.Println("Need a token.\n")
		flag.Usage()
		os.Exit(1)
	}

	user, err := getUser(*profile)
	if err != nil {
		fmt.Println("Couldn't get user:", err)
		os.Exit(1)
	}
	fmt.Println("User:", user)

	mfaDevice, err := getMFADevice(*profile, user)
	if err != nil {
		fmt.Println("Couldn't get MFA device:", err)
		os.Exit(1)
	}
	fmt.Println("MFA device:", mfaDevice)

	accessKeyID, secretAccessKey, sessionToken, err := getCredentials(*profile, mfaDevice, *token)
	if err != nil {
		fmt.Println("Couldn't get credentials:", err)
		os.Exit(1)
	}

	fmt.Println("Setting credentials...")
	err = configureCredentials(*profileMFA, *region, accessKeyID, secretAccessKey, sessionToken)
	if err != nil {
		fmt.Println("Couldn't set profile credentials:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
