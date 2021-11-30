package controllers

import (
	"context"
	"fmt"
	"os"

	//"time"

	//"github.com/containers/common/pkg/retry"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"

	//"github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

type RegistryManager interface {
	CopyImage(ctx context.Context, srcImage, dstImage string, srcRegistryCredentials, dstCredentials *RegistryCredentials) error
}

type ContainerRegistryManager struct {
}

type RegistryCredentials struct {
	URL      string
	Username string
	Password string
}

func (c *ContainerRegistryManager) CopyImage(ctx context.Context, srcImage, dstImage string, srcRegistryCredentials, dstCredentials *RegistryCredentials) error {

	srcImage = "docker://" + srcImage
	dstImage = "docker://" + dstImage

	srcRef, err := alltransports.ParseImageName(srcImage)
	if err != nil {
		return fmt.Errorf("invalid source name %s: %v", srcImage, err)
	}
	destRef, err := alltransports.ParseImageName(dstImage)
	if err != nil {
		return fmt.Errorf("invalid destination name %s: %v", dstImage, err)
	}

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	if err != nil {
		return fmt.Errorf("failed to get default policy: %v", err)
	}
	policyCtx, err := signature.NewPolicyContext(policy)
	if err != nil {
		return fmt.Errorf("failed to get default policy: %v", err)
	}

	srcCtx := &types.SystemContext{
		OSChoice:      "linux",
		VariantChoice: "amd64",
	}
	if srcRegistryCredentials != nil {
		srcCtx.DockerAuthConfig = &types.DockerAuthConfig{
			Username: srcRegistryCredentials.Username,
			Password: srcRegistryCredentials.Password,
		}
	}

	dstCtx := &types.SystemContext{
		DockerAuthConfig: &types.DockerAuthConfig{
			Username: dstCredentials.Username,
			Password: dstCredentials.Password,
		},
	}

	_, err = copy.Image(ctx, policyCtx, destRef, srcRef, &copy.Options{
		SourceCtx:      srcCtx,
		DestinationCtx: dstCtx,
		ReportWriter:   os.Stdout,
	})

	//return retry.RetryIfNecessary(ctx, func() error {
	//	_, err := copy.Image(ctx, policyCtx, destRef, srcRef, &copy.Options{
	//		SourceCtx: &types.SystemContext{
	//			OSChoice:      "linux",
	//			VariantChoice: "amd64",
	//		},
	//		ReportWriter: os.Stdout,
	//	})
	//	if err != nil {
	//		return err
	//	}
	//
	//	return nil
	//}, &retry.RetryOptions{MaxRetry: 1, Delay: time.Second * 5})

	return err
}
