module github.com/junaidk/image-backup-controller

go 1.16

require (
	cloud.google.com/go v0.81.0 // indirect
	github.com/containers/image/v5 v5.17.0
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/utils v0.0.0-20211116205334-6203023598ed
	sigs.k8s.io/controller-runtime v0.10.0
)
