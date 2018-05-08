#!/usr/bin/env bash
#start api connect logic 
baseDir=`pwd`
export GOPATH=$baseDir/api
cd $baseDir/api/src
go get 
go build apimain.go
nohup ./apimain & >/dev/null 2>&1 
echo "api start success"
export GOPATH=$baseDir/connect
cd $baseDir/connect/src
go get
go build connectmain.go
nohup ./connectmain & >/dev/null 2>&1 
echo "connect start success"
export GOPATH=$baseDir/logic
cd $baseDir/logic/src
go get
go build logicmain.go
nohup ./logicmain & >/dev/null 2>&1 
echo "logic start success"
