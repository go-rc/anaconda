package anaconda

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type OEmbed struct {
	Type          string
	Width         int
	Cache_age     string
	Height        int
	Author_url    string
	Html          string
	Version       string
	Provider_name string
	Provider_url  string
	Url           string
	Author_name   string
}

// No authorization on this endpoint. Its the only one.
func (a TwitterApi) GetOEmbed(v url.Values) (o OEmbed, err error) {
	var f io.Writer

	filename := filepath.Join(append([]string{"json"}, strings.Split("/statuses/oembed.json", "/")...)...)
	fl, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return o, err
	}
	defer fl.Close()
	f = fl

	resp, err := http.Get(a.baseUrlV1() + "/statuses/oembed.json?" + v.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
	err = decodeResponse(resp, &o)
	return
}

// Calls GetOEmbed with the corresponding id. Convenience wrapper for GetOEmbed()
func (a TwitterApi) GetOEmbedId(id int64, v url.Values) (o OEmbed, err error) {
	if v == nil {
		v = url.Values{}
	}
	v.Set("id", strconv.FormatInt(id, 10))
	resp, err := http.Get(a.baseUrlV1() + "/statuses/oembed.json?" + v.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeResponse(resp, &o)
	return
}

func (a TwitterApi) baseUrlV1() string {
	if a.baseUrl == BaseUrl {
		return BaseUrlV1
	}

	if a.baseUrl == "" {
		return BaseUrlV1
	}

	return a.baseUrl
}
