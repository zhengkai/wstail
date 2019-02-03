SHELL:=/bin/bash

default:
	#server
	mkdir -p ../server/pb
	protoc --go_out=../server/pb *.proto
	#client
	mkdir -p ../client/src/pb
	../client/node_modules/protobufjs/bin/pbjs -t static-module -w es6 --es6 --force-number -o ../client/src/pb/pb.js *.proto
	../client/node_modules/protobufjs/bin/pbts -o ../client/src/pb/pb.d.ts ../client/src/pb/pb.js
