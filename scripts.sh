#!/bin/bash

install_path() {
  	INSTALL_COMMAND="export PATH=\${PATH}:$(pwd)/bin"
  	COMMAND_EXIST=$(cat ~/.zshrc | grep ${INSTALL_COMMAND})
  	if [[ -z $(cat ~/.zshrc | grep "${INSTALL_COMMAND}") ]]; then echo "${INSTALL_COMMAND}" >> ~/.zshrc; fi
}

build() {
  make build
}

clean() {
  make clean
}


install () {
  clean
  build
  install_path
  printf "\\033[1;32m%s\\033[0m\\n" "Installation Complete"
}



install