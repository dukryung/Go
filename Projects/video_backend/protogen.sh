#!/bin/bash

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

for dir in $proto_dirs; do
  echo "$dir"
done

for dir in $proto_dirs; do
  protoc -I="proto" \
   -I="third_party/proto" \
   --go_out=. \
   --go-grpc_out=. \
   --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done


#cp -r github.com/hessegg/klaatoo-faucet/* ./
#rm -rf github.com
