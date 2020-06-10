package utils

import (
	"io"
	"net/http"
)

func HttpClient(url string) (io.ReadCloser, error) {
	client := &http.Client{}
	//request, _ := http.NewRequest("GET", "https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1628265209,77028583&fm=15&gp=0.jpg", nil)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	request.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	request.Header.Set("Cache-Control", "max-age=0")
	request.Header.Set("Connection", "keep-alive")

	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		return response.Body, nil
	} else {
		return nil, nil
	}
}
