#!/bin/sh

BIN_DIR=/usr/local/bin

# install curl, git
apk add curl git

# download kubectl which fits the clusters version
KUBERNETES_VERSION=$(curl -s --cacert /run/secrets/kubernetes.io/serviceaccount/ca.crt --header "Authorization: Bearer $(< /run/secrets/kubernetes.io/serviceaccount/token)" https://${KUBERNETES_SERVICE_HOST}/version | grep -Eo 'v1\.[0-9]+\.[0-9]+')
curl --output ${BIN_DIR}/kubectl --location https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl
chmod +x ${BIN_DIR}/kubectl

# download helm
HELM_VERSION=v3.2.3
curl --output ${BIN_DIR}/helm.tar.gz --location https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz
cd ${BIN_DIR}
tar xvzf helm.tar.gz
mv linux-amd64/helm .
rm -rf linux-amd64 helm.tar.gz
cd -

# download yq
YQ_VERSION=$(curl -s https://github.com/mikefarah/yq/releases/latest | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
curl --output ${BIN_DIR}/yq --location https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64
chmod +x ${BIN_DIR}/yq
