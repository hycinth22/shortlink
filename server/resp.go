package server

import (
	"encoding/json"
	"net/http"
)

type RespContent struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data struct {
		Id      int64  `json:"id,omitempty"`
		Link    string `json:"link,omitempty"`
		RawLink string `json:"rawlink,omitempty"`
	} `json:"data,omitempty"`
}

type Resp struct {
	// Code is constants http.Status*
	Code    int
	w       http.ResponseWriter
	Content *RespContent
}

func NewResp(w http.ResponseWriter) Resp {
	c := Resp{
		w:       w,
		Code:    http.StatusOK,
		Content: nil,
	}
	c.Content = new(RespContent)
	c.Content.Success = false
	c.Content.Msg = ""
	c.Content.Data.Id = 0
	c.Content.Data.Link = ""
	c.Content.Data.RawLink = ""
	return c
}

func (c Resp) DefContByCode() {
	c.Content.Success = c.Code == http.StatusOK
	c.Content.Msg = http.StatusText(c.Code)
}

func (c Resp) WriteMsg(msg string) {
	if c.Content.Msg != "" {
		c.Content.Msg += "\n"
	}
	c.Content.Msg += msg
}

func (c Resp) Exec() error {
	c.w.WriteHeader(c.Code)
	bytes, err := json.Marshal(c.Content)
	if err != nil {
		return err
	}
	c.w.Write(bytes)
	return nil
}
