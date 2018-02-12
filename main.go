package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"shortLink/data"
	"shortLink/server"
)

func handleStub(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}

const ShortPrefix = "/t"
const AddPrefix = "/add"

type shortLinkHandler struct {
}

func (c shortLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[1:]
	rawLink, err := data.GetRawLink(token)
	if err != nil {
		println("Failed to serve token:", err.Error())
		if data.NoTokenErr(err) {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, http.StatusText(http.StatusNotFound))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	println("Serve successly:", token, rawLink)
	http.Redirect(w, r, rawLink, http.StatusFound)
}

type addHandler struct {
}

func (c addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			resp := server.NewResp(w)
			if err, ok := err.(server.BadRequestError); ok {
				println("Bad add request")
				resp.Code = http.StatusBadRequest
				resp.DefContByCode()

				println("     " + err.Error())
				resp.WriteMsg(err.Error())
			}
			if err, ok := err.(server.ServerError); ok {
				println("Add error: ")
				resp.Code = http.StatusInternalServerError
				resp.DefContByCode()

				println("     " + err.Error())
				resp.WriteMsg(err.Error())
			}
			if err, ok := err.(error); ok {
				println(err.Error())
			}
			resp.Exec()
		}
	}()

	if r.Method != http.MethodPost {
		panic(server.BadRequestError{"Method Should Be Post"})
	}

	rawLink := r.PostFormValue("rawLink")
	if rawLink == "" {
		panic(server.BadRequestError{"Invalid Request Body"})
	}

	// TODO: 实现客户端自定义token
	// ...
	var token string
	md5 := md5.Sum([]byte(rawLink))
	token = hex.EncodeToString(md5[:])

	var id int64
	var err error
	id, err = data.InsertLink(token, rawLink)
	if err != nil {
		panic(server.ServerError{err})
	}

	println("Add link:", token, rawLink)
	resp := server.NewResp(w)
	cont := resp.Content
	cont.Success = true
	cont.Data.Id = id
	cont.Data.Link = ShortPrefix + "/" + token
	cont.Data.RawLink = rawLink
	resp.Code = http.StatusOK
	resp.WriteMsg("操作成功")
	resp.Exec()
}

func main() {
	http.Handle(ShortPrefix+"/", http.StripPrefix(ShortPrefix, shortLinkHandler{}))
	http.Handle(AddPrefix+"/", http.StripPrefix(AddPrefix, addHandler{}))
	http.Handle("/", http.FileServer(http.Dir("html/")))
	http.ListenAndServe(":80", nil)
}
