package requestx

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ParsePathMap
/*
	从请求体中获取path请求参数
	return: map[string]any
*/
func ParsePathMap(r *gin.Context) (map[string]any, error) {
	vars := r.Params
	paramMap := make(map[string]any)
	for _, p := range vars {
		paramMap[p.Key] = p.Value
	}

	return paramMap, nil
}

// ParamFormMap
/*
	从请求体中获取form请求参数
	return: map[string]any
*/
func ParamFormMap(r *gin.Context) (map[string]any, error) {
	if err := r.Request.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.Request.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]any, len(r.Request.Form))
	for name := range r.Request.Form {
		formValue := r.Request.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	if len(params) == 0 {
		return nil, nil
	}
	return params, nil
}

// ParamJsonMap
/*
	从请求体中获取json请求参数
	return: map[string]any
*/
func ParamJsonMap(r *gin.Context) (map[string]any, error) {
	var buf bytes.Buffer
	var bodyReader *bytes.Reader
	if _, err := io.Copy(&buf, r.Request.Body); err != nil {
		return nil, err
	} else {
		bodyReader = bytes.NewReader(buf.Bytes())
	}
	defer func() {
		r.Request.Body = io.NopCloser(&buf)
	}()

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		if err.Error() != "unexpected end of JSON input" {
			return nil, err
		}
	}
	return result, nil
}

// DeepCopyRequest
/*
	深拷贝一个*http.Request
*/
func DeepCopyRequest(r *http.Request) (*http.Request, error) {
	var buf bytes.Buffer
	var bodyReader *bytes.Reader
	if _, err := io.Copy(&buf, r.Body); err != nil {
		return nil, err
	} else {
		bodyReader = bytes.NewReader(buf.Bytes())
	}
	defer func() {
		r.Body = io.NopCloser(&buf)
	}()

	// 创建新的Body
	newBody := io.NopCloser(bodyReader)

	// 创建http.Request的副本
	var clonedRequest http.Request

	// 拷贝 Method
	clonedRequest.Method = r.Method

	// 拷贝 URL
	clonedRequest.URL = r.URL

	// 拷贝 Proto
	clonedRequest.Proto = r.Proto

	// 拷贝 ProtoMajor
	clonedRequest.ProtoMajor = r.ProtoMajor

	// 拷贝 ProtoMinor
	clonedRequest.ProtoMinor = r.ProtoMinor

	// 拷贝 Header
	clonedRequest.Header = make(http.Header)
	for key, values := range r.Header {
		for _, value := range values {
			clonedRequest.Header.Add(key, value)
		}
	}

	// 拷贝 Body
	clonedRequest.Body = newBody

	// 拷贝 ContentLength
	clonedRequest.ContentLength = r.ContentLength

	// 拷贝 TransferEncoding
	clonedRequest.TransferEncoding = append([]string(nil), r.TransferEncoding...)

	// 拷贝 Close
	clonedRequest.Close = r.Close

	// 拷贝 Host
	clonedRequest.Host = r.Host

	// 拷贝 Form
	clonedRequest.Form = r.Form

	// 拷贝 PostForm
	clonedRequest.PostForm = r.PostForm

	// 拷贝 MultipartForm
	clonedRequest.MultipartForm = r.MultipartForm

	// 拷贝 Trailer
	clonedRequest.Trailer = r.Trailer

	// 拷贝 RemoteAddr
	clonedRequest.RemoteAddr = r.RemoteAddr

	// 拷贝 RequestURI
	clonedRequest.RequestURI = r.RequestURI

	// 拷贝 TLS
	clonedRequest.TLS = r.TLS

	// 拷贝 Context
	clonedRequest = *clonedRequest.WithContext(r.Context())

	return &clonedRequest, nil
}
