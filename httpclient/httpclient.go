package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const tag = "httpclient:"

var defaultRequestSettings = RequestSetting{
	Retries:   1,
	DumpBody:  true,
	UserAgent: "httpclient",
	Debug:     true,
}

var mux sync.Mutex
var defaultClient = http.DefaultClient

func Client(client *http.Client) {
	if client == nil {
		panic("client is nil")
	}
	mux.Lock()
	defer mux.Unlock()
	defaultClient = client
}
func RequestSettings(settings RequestSetting) {
	mux.Lock()
	defer mux.Unlock()
	defaultRequestSettings = settings
}

type RequestSetting struct {
	// if set to -1 means will retry forever
	Retries   int
	UserAgent string
	Debug     bool
	DumpBody  bool
}

type HttpRequest struct {
	url      string
	params   map[string][]string
	files    map[string]string
	req      *http.Request
	resp     *http.Response
	body     []byte
	dump     []byte
	settings RequestSetting
	finished bool
}

func Get(url string) *HttpRequest {
	return NewHttpRequest(url, "GET")
}

func Post(url string) *HttpRequest {
	return NewHttpRequest(url, "POST")
}

func Put(url string) *HttpRequest {
	return NewHttpRequest(url, "PUT")
}

func Delete(url string) *HttpRequest {
	return NewHttpRequest(url, "DELETE")
}

func Head(url string) *HttpRequest {
	return NewHttpRequest(url, "HEAD")
}

func NewHttpRequest(rawurl, method string) *HttpRequest {
	if defaultClient == nil {
		panic("Should call DefaultClient() or Client() first")
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Println(tag, err)
		return nil
	}
	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &HttpRequest{
		url:      rawurl,
		params:   map[string][]string{},
		files:    map[string]string{},
		req:      &req,
		settings: defaultRequestSettings,
	}
}

func (r *HttpRequest) Param(key, value string) *HttpRequest {
	if param, ok := r.params[key]; ok {
		r.params[key] = append(param, value)
	} else {
		r.params[key] = []string{value}
	}
	return r
}

func (r *HttpRequest) File(formField, filepath string) *HttpRequest {
	r.files[formField] = filepath
	return r
}

func (r *HttpRequest) Body(data interface{}) *HttpRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		r.req.Body = ioutil.NopCloser(bf)
		r.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		r.req.Body = ioutil.NopCloser(bf)
		r.req.ContentLength = int64(len(t))
	}
	return r
}

func (r *HttpRequest) JSONBody(obj interface{}) (*HttpRequest, error) {
	if r.req.Body != nil {
		panic("the request already has a request body")
	}
	if obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return r, err
		}
		r.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		r.req.ContentLength = int64(len(byts))
		r.req.Header.Set("Content-Type", "application/json")
	}
	return r, nil
}

func (r *HttpRequest) AddHeader(key, value string) *HttpRequest {
	r.req.Header.Add(key, value)
	return r
}

func (r *HttpRequest) SetHeader(key, value string) *HttpRequest {
	r.req.Header.Set(key, value)
	return r
}

func (r *HttpRequest) AddCookie(cookie *http.Cookie) *HttpRequest {
	r.req.AddCookie(cookie)
	return r
}

func (r *HttpRequest) Retries(retries int) *HttpRequest {
	r.settings.Retries = retries
	return r
}

func (r *HttpRequest) DumpBody(dumpbody bool) *HttpRequest {
	r.settings.DumpBody = dumpbody
	return r
}

func (r *HttpRequest) UserAgent(useragent string) *HttpRequest {
	r.settings.UserAgent = useragent
	return r
}

func (r *HttpRequest) Debug(debug bool) *HttpRequest {
	r.settings.Debug = debug
	return r
}

func (r *HttpRequest) DoRequest() (resp *http.Response, err error) {
	var paramBody string
	length := len(r.params)
	if length > 0 {
		var buf bytes.Buffer
		keys := make([]string, 0, length)
		for k := range r.params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vs := r.params[k]
			prefix := url.QueryEscape(k) + "="
			for _, v := range vs {
				if buf.Len() > 0 {
					buf.WriteByte('&')
				}
				buf.WriteString(prefix)
				buf.WriteString(url.QueryEscape(v))
			}
		}
		paramBody = buf.String()
	}

	r.build(paramBody)

	httpurl, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	r.req.URL = httpurl

	if r.settings.UserAgent != "" && r.req.Header.Get("User-Agent") == "" {
		r.req.Header.Set("User-Agent", r.settings.UserAgent)
	}

	if r.settings.Debug {
		dump, err := httputil.DumpRequest(r.req, r.settings.DumpBody)
		if err != nil {
			log.Println(tag, err.Error())
		}
		r.dump = dump
	}

	for i := 0; r.settings.Retries == -1 || i <= r.settings.Retries; i++ {
		resp, err = defaultClient.Do(r.req)
		if err == nil {
			break
		}
	}
	return resp, err
}

func (r *HttpRequest) build(paramBody string) {
	if r.req.Method == "GET" {
		if len(paramBody) > 0 {
			if strings.Contains(r.url, "?") {
				r.url += "&" + paramBody
			} else {
				r.url += "?" + paramBody
			}
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if r.req.Method == "POST" || r.req.Method == "PUT" || r.req.Method == "PATCH" || r.req.Method == "DELETE" {
		// with files
		if len(r.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formField, file := range r.files {
					fileWriter, err := bodyWriter.CreateFormFile(formField, filepath.Base(file))
					if err != nil {
						log.Println(tag, err)
						continue
					}
					fh, err := os.Open(file)
					if err != nil {
						log.Println(tag, err)
						continue
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Println(tag, err)
					}
				}
				for k, v := range r.params {
					for _, vv := range v {
						err := bodyWriter.WriteField(k, vv)
						if err != nil {
							log.Println(tag, err)
						}
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			r.SetHeader("Content-Type", bodyWriter.FormDataContentType())
			r.req.Body = ioutil.NopCloser(pr)
			return
		}
	}

	if len(paramBody) > 0 {
		r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		r.Body(paramBody)
	}
}

func (r *HttpRequest) Response() (*http.Response, error) {
	return r.getResponse()
}

func (r *HttpRequest) Bytes() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	resp, err := r.getResponse()
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		r.body, err = ioutil.ReadAll(reader)
		return r.body, err
	}
	r.body, err = ioutil.ReadAll(resp.Body)
	return r.body, err
}

func (r *HttpRequest) ToString() (string, error) {
	data, err := r.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *HttpRequest) ToJSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (r *HttpRequest) ToFile(filename string, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := r.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func (r *HttpRequest) Dump() []byte {
	return r.dump
}

func (r *HttpRequest) getResponse() (*http.Response, error) {
	if r.finished {
		return r.resp, nil
	}
	resp, err := r.DoRequest()
	r.finished = true
	if err != nil {
		return nil, err
	}
	r.resp = resp
	return resp, nil
}
