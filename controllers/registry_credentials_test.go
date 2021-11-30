package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthFromDockerCfg(t *testing.T) {

	testData := []byte(`{
    "auths": {
        "https://index.docker.io": {
            "username": "junaidk",
            "password": "abcxd",
            "auth": "anVuYWlkazphYmN4ZA=="
        }
    }
}`)
	creds, err := getRegistryCredentialsFromDockerCfg(testData)

	assert.NoError(t, err)

	assert.Equal(t, creds.Username, "junaidk")
	assert.Equal(t, creds.Password, "abcxd")
}

func TestGetDockerCfgSecret(t *testing.T) {

	out, err := getDockerConfigSecret("alpha", "beta", "https://index.docker.io")
	assert.NoError(t, err)
	expected := []byte(`
{"auths":{"https://index.docker.io":{"username":"alpha","password":"beta"}}}`)

	assert.JSONEqf(t, string(expected), string(out.Data[".dockerconfigjson"]), "Expected %s to be %s", string(out.Data[".dockerconfigjson"]), string(expected))
}
