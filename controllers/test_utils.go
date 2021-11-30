package controllers

import "context"

type TestRegistryManager struct {
	copyImageStub func(srcImage, dstImage string, srcRegistryCredentials, dstRegistryCredentials *RegistryCredentials)
}

func (tr *TestRegistryManager) CopyImage(ctx context.Context, srcImage, dstImage string, srcRegistryCredentials, dstRegistryCredentials *RegistryCredentials) error {
	tr.copyImageStub(srcImage, dstImage, srcRegistryCredentials, dstRegistryCredentials)
	return nil
}

var SrcImageNames = []string{"library/image1", "quay.io/notcache/image2"}
var DstImageNames = []string{"index.docker.io/user/image1", "index.docker.io/user/image2"}

var DstRegAuth = []byte(`{
    "auths": {
        "index.docker.io": {
            "username": "user",
            "password": "password"
        }
    }
}`)

var DstRegistryCredentials = &RegistryCredentials{
	URL:      "index.docker.io",
	Username: "user",
	Password: "password",
}

var SrcRegAuth1 = []byte(`{
    "auths": {
        "index.docker.io": {
            "username": "user1",
            "password": "password1"
        }
    }
}`)

var SrcRegistryCredentials1 = &RegistryCredentials{
	URL:      "index.docker.io",
	Username: "user1",
	Password: "password1",
}

var SrcRegAuth2 = []byte(`{
    "auths": {
        "quay.io": {
            "username": "user2",
            "password": "password2"
        }
    }
}`)

var SrcRegistryCredentials2 = &RegistryCredentials{
	URL:      "quay.io",
	Username: "user2",
	Password: "password2",
}

var SrcRegistryCredentialList = []*RegistryCredentials{SrcRegistryCredentials1, SrcRegistryCredentials2}
