package controllers

import (
	"errors"
	"os"
	"strings"
)

func GetIgnoreNamespacesEnv() []string {
	var ignoreNamespaceEnvVar = "IGNORE_NAMESPACES"

	var ignoreNamespaces []string
	ns, found := os.LookupEnv(ignoreNamespaceEnvVar)
	if !found {
		ignoreNamespaces = []string{"kube-system"}
	} else {
		ignoreNamespaces = strings.Split(ns, ",")
	}
	return ignoreNamespaces
}

func GetBackUpRegistryURLEnv() (string, error) {
	var backUpRegistryURLEnvVar = "BACKUP_REGISTRY_URL"

	env, found := os.LookupEnv(backUpRegistryURLEnvVar)
	if !found {
		return "", errors.New(backUpRegistryURLEnvVar + " environment variable not found")
	}

	if strings.HasPrefix(env, "http://") {
		env = strings.TrimPrefix(env, "http://")
	} else if strings.HasPrefix(env, "https://") {
		env = strings.TrimPrefix(env, "https://")
	}

	if strings.HasSuffix(env, "/") {
		env = strings.TrimSuffix(env, "/")
	}

	return env, nil
}

func GetBackUpRegistryUserNameEnv() (string, error) {
	var backUpRegistryUserNameEnvVar = "BACKUP_REGISTRY_USERNAME"

	env, found := os.LookupEnv(backUpRegistryUserNameEnvVar)
	if !found {
		return "", errors.New(backUpRegistryUserNameEnvVar + " environment variable not found")
	}
	return env, nil
}

func GetBackUpRegistryPasswordEnv() (string, error) {
	var backUpRegistryPasswordEnvVar = "BACKUP_REGISTRY_PASSWORD"

	env, found := os.LookupEnv(backUpRegistryPasswordEnvVar)
	if !found {
		return "", errors.New(backUpRegistryPasswordEnvVar + " environment variable not found")
	}
	return env, nil
}

func GetPodNameSpaceEnv() string {
	var nameSpaceEnvVar = "MY_POD_NAMESPACE"

	env, found := os.LookupEnv(nameSpaceEnvVar)
	if !found {
		return ""
	}
	return env
}
