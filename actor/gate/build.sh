#!/bin/bash
#export PROTOPATH=../../.godeps/src
#export PROTOPATH2=../

cd `dirname $0`
protoc3 -I=. -I=$GOPATH/src/github.com/gogo/protobuf/protobuf -I=$GOPATH/src --csharp_out=. --gogoslick_out=plugins=grpc:. protos.proto  && cp Protos.cs ../../cscode/Client/Assets/pb3net/GateProtos.cs
cd -
