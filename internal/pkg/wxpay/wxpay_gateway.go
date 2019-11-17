package wxpay

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	httpTransport *http.Transport
)

func init() {
	//要求发起Http请求时候，忽略证书签名。否则，可能触发（x509: certificate signed by unknown authority）
	httpTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func Post(xml string, url string) (string, error) {

	//httpClient := &http.Client{Transport: httpTransport}
	//resp, err := httpClient.Post(url, "text/xml", strings.NewReader(xml))
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
