package common

import (
	"crypto/tls"
	"image"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func HttpBody(url, method, param string, header map[string]string) (int, []byte) {
	re_status := 100
	par := strings.NewReader(param)
	req, err := http.NewRequest(method, url, par)

	if err != nil {
		panic(err)
		return re_status, nil
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	if len(header) < 1 {
		header["Content-Type"] = "text/plain; charset=UTF8"
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36"
		header["Accept"] = "application/json, text/plain, */*"
	}
	for h_k, h_v := range header {
		req.Header.Set(h_k, h_v)
	}

	tr := &http.Transport{DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}

	//处理返回结果
	resp, err := client.Do(req)
	if err != nil {
		return re_status, nil
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	re_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return re_status, nil
	}

	re_status = resp.StatusCode

	return re_status, re_body
}

func HttpBodyResponse(url, method, param string, header map[string]string) (*http.Response, error) {
	par := strings.NewReader(param)
	req, err := http.NewRequest(method, url, par)
	if err != nil {
		panic(err)
		return nil, err
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	if len(header) < 1 {
		header["Content-Type"] = "text/plain; charset=UTF8"
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36"
		header["Accept"] = "application/json, text/plain, */*"
	}
	for h_k, h_v := range header {
		req.Header.Set(h_k, h_v)
	}

	//处理返回结果
	tr := &http.Transport{DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	return resp, err
}

func HttpsBodyResponse(url, method, param string, header map[string]string) (*http.Response, error) {
	par := strings.NewReader(param)

	req, err := http.NewRequest(method, url, par)
	if err != nil {
		return nil, err
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	if len(header) > 0 {
		for h_key, h_val := range header {
			req.Header.Set(h_key, h_val)
		}
	} else {
		req.Header.Set("Content-Type", "charset=UTF-8")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Connection", "	keep-alive")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	}
	_tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: _tlsConfig,
		DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	return resp, err
}

func HttpBodyByImg(url, method, param string, header map[string]string) (int, image.Image) {
	re_status := 100
	par := strings.NewReader(param)
	req, err := http.NewRequest(method, url, par)

	if err != nil {
		panic(err)
		return re_status, nil
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	if len(header) < 1 {
		header["Content-Type"] = "text/plain; charset=UTF8"
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36"
		header["Accept"] = "application/json, text/plain, */*"
	}
	for h_k, h_v := range header {
		req.Header.Set(h_k, h_v)
	}

	tr := &http.Transport{DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}

	//处理返回结果
	resp, err := client.Do(req)
	if err != nil {
		return re_status, nil
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return re_status, nil
	}

	re_status = resp.StatusCode

	img, _, err := image.Decode(resp.Body)

	return re_status, img
}

func PushHttpBody(url, method, param string, header map[string]string) (int, []byte) {
	re_status := 100
	par := strings.NewReader(param)
	req, err := http.NewRequest(method, url, par)

	if err != nil {
		panic(err)
		return re_status, nil
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	if len(header) < 1 {
		header["Content-Type"] = "text/plain; charset=UTF8"
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36"
		header["Accept"] = "application/json, text/plain, */*"
	}
	for h_k, h_v := range header {
		req.Header.Set(h_k, h_v)
	}

	client := &http.Client{}

	//处理返回结果
	resp, err := client.Do(req)
	if err != nil {
		return re_status, nil
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	re_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return re_status, nil
	}

	re_status = resp.StatusCode

	return re_status, re_body
}
