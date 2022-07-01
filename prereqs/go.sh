#!/bin/bash
cd $HOME
wget -c https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
sudo rm -rf go
tar xvf go1.18.3.linux-amd64.tar.gz
cat >> $HOME/.bashrc <<- EOM
export GOPATH=\$HOME
export GOROOT=\$GOPATH/go
export GOBIN=\$GOPATH/bin
export PATH=\$GOBIN:\$GOROOT/bin:\$PATH
EOM
source $HOME/.bashrc
