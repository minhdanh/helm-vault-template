#!/bin/bash

set -ueo pipefail

VERSION=v0.1.0

function isAlreadyInstalled() {
  hash helm-vault-template 2>/dev/null && [[ $(helm-vault-template -v | cut -d " " -f 3) == ${VERSION} ]]
}

if isAlreadyInstalled; then
  echo "helm-vault-template is already installed"
else
  echo "Downloading helm-vault-template version ${VERSION}"
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  URL=https://github.com/minhdanh/helm-vault-template/releases/download/${VERSION}/helm-vault-template_${OS}_amd64

  temp_file=$(mktemp)
  trap "rm ${temp_file}" EXIT

  statuscode=$(curl -w "%{http_code}" -sL ${URL} -o ${temp_file})

  if [[ ! "${statuscode}" == "200" ]]; then
    echo "Failed to download binary"
    exit 1
  fi

  cp ${temp_file} /usr/local/bin/helm-vault-template
  chmod +x /usr/local/bin/helm-vault-template
fi
