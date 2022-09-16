package main

import (
	"fmt"
	"strings"
)

func CredentialsFromPointers(envPtr, usrPtr, pwPtr, dbPtr, portPtr, hostPtr *string) (*Credentials, error) {
	if *usrPtr != "" && *pwPtr != "" {
		return NewCredentials(*usrPtr, *pwPtr, *dbPtr, *portPtr, *hostPtr), nil
	} else {
		if *envPtr == "" {
			return nil, fmt.Errorf("neither .env path nor credentials where provided")
		}
		if !strings.HasSuffix(*envPtr, ".env") {
			return nil, fmt.Errorf("%s is not a .env file", *envPtr)
		}
		return LoadCredentialsCustomEnvPath(*envPtr)
	}
}
