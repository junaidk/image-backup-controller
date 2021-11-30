package controllers

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type authConfiguration struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
}

type authConfigurations struct {
	Configs map[string]authConfiguration `json:"configs"`
}

type dockerConfig struct {
	Auth     string `json:"auth"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// code snippet taken from https://github.com/fsouza/go-dockerclient/blob/main/auth.go
func parseDockerConfig(dockerCfg []byte) (*authConfigurations, error) {
	confsWrapper := struct {
		Auths map[string]dockerConfig `json:"auths"`
	}{}
	if err := json.Unmarshal(dockerCfg, &confsWrapper); err != nil {
		return nil, fmt.Errorf("failed to parse docker config: %v", err)
	}

	c := &authConfigurations{
		Configs: make(map[string]authConfiguration),
	}
	for reg, conf := range confsWrapper.Auths {
		authConfig := authConfiguration{
			Username:      conf.Username, //userpass[0],
			Password:      conf.Password, //userpass[1],
			ServerAddress: reg,
		}

		c.Configs[reg] = authConfig
	}

	return c, nil
}

func getRegistryCredentialsFromDockerCfg(dockerCfg []byte) (*RegistryCredentials, error) {
	var auth *authConfigurations
	auth, err := parseDockerConfig(dockerCfg)
	if err != nil {
		return nil, err
	}

	var regCreds RegistryCredentials
	for k, v := range auth.Configs {
		regCreds = RegistryCredentials{
			URL:      k,
			Username: v.Username,
			Password: v.Password,
		}
	}

	return &regCreds, nil
}

func getDockerConfigSecret(username, password, registryUrl string) (*corev1.Secret, error) {

	authConfig := authConfiguration{
		Username: username,
		Password: password,
	}

	authConfigs := authConfigurations{
		Configs: map[string]authConfiguration{
			registryUrl: authConfig,
		},
	}

	confsWrapper := struct {
		Auths map[string]authConfiguration `json:"auths"`
	}{
		Auths: authConfigs.Configs,
	}

	encodedDockerConfig, err := json.Marshal(confsWrapper)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "destination-registry-creds"},
		Immutable:  nil,
		Data: map[string][]byte{
			".dockerconfigjson": encodedDockerConfig,
		},
		StringData: nil,
		Type:       corev1.SecretTypeDockerConfigJson,
	}
	return secret, nil
}
