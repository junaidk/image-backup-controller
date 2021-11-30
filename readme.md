
operator-sdk init --domain junaidk --repo github.com/junaidk/image-backup-operator

operator-sdk create api  --version v1alpha1 --kind ImageBackup --resource --controlle

make generate
make manifests


operator-sdk olm install

make bundle bundle-build bundle-push junaidk/image-backup-operator:1.0



docker login -u junaidk
aae4f2b6-2f5d-4e9a-9611-f8d178c144dc

kubectl create secret docker-registry regcred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username=junaidk \
    --docker-password=aae4f2b6-2f5d-4e9a-9611-f8d178c144dc 

index.docker.io/<user>/reponame:tag

https in url inside secret

----

ENV variables
export BACKUP_REGISTRY_URL="index.docker.io"
export BACKUP_REGISTRY_USERNAME="junaidk"
export BACKUP_REGISTRY_PASSWORD="aae4f2b6-2f5d-4e9a-9611-f8d178c144dc"
export IGNORE_NAMESPACES="kube-system,kube-public,kube-node-lease,olm,operators"
