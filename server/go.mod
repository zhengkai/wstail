module wstail

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/gorilla/websocket v1.4.0
	github.com/zhengkai/rome v1.0.0
	golang.org/x/sys v0.0.0-20190219203350-90b0e4468f99 // indirect
	pb v0.0.0
)

replace pb => ./pb

replace github.com/zhengkai/rome => /go/src/github.com/zhengkai/rome
