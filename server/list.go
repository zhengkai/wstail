package main

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"pb"

	"github.com/gogo/protobuf/proto"
)

var (
	fileNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\.\-]*\.(txt|log)$`)
)

func doFile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(`Access-Control-Allow-Origin`, `*`)

	list := listFile()
	if list == nil {
		return
	}

	b, _ := proto.Marshal(list)

	w.Write(b)
}

func listFile() (r *pb.FileReturn) {

	r = &pb.FileReturn{
		Base: makeOpBaseReturn(),
	}

	files, err := ioutil.ReadDir(dirBase)
	if err != nil {
		r.Base.Error = `read dir fail`
		return
	}

	list := []string{}
	for _, file := range files {
		name := file.Name()
		if !checkFileName(name) {
			continue
		}
		list = append(list, name)
	}

	r.File = list

	return
}

func checkFileName(s string) (ok bool) {
	return fileNameRegexp.MatchString(s)
}
