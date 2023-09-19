package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

var convExtCandidates = []string{".jpg", ".png", ".jpeg"}

func doWebp(req *http.Request) (*http.Response, error) {
	var orgRes *http.Response
	for _, cExt := range convExtCandidates {
		orgURL := req.URL
		newPath := orgURL.Path[:len(orgURL.Path)-len(".webp")] + cExt
		newOrgURL, err := url.Parse(fmt.Sprintf("%s://%s%s?%s", orgURL.Scheme, orgURL.Host, newPath, orgURL.RawQuery))
		if err != nil {
			log.Println(err)
			continue
		}
		newReq, err := http.NewRequest("GET", newOrgURL.String(), nil)
		newReq.Header = req.Header
		if err != nil {
			log.Println(err)
			continue
		}
		orgRes, err = client.Do(newReq)
		if err != nil {
			log.Println(err)
			continue
		}
		if orgRes.StatusCode != 200 {
			log.Println(orgRes.Status)
			continue
		} else {
			break
		}
	}
	if orgRes == nil {
		return nil, fmt.Errorf("get origin failed")
	}
	return orgRes, nil
}

func convWebp(src io.Reader, params []string) (*bytes.Buffer, error) {
	f, err := os.CreateTemp("/tmp", "")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	_, err = io.Copy(f, src)
	if err != nil {
		return nil, err
	}
	params = append(params, "-quiet", "-mt", "-jpeg_like", f.Name(), "-o", "-")
	out, err := exec.Command("cwebp", params...).Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bytes.NewBuffer(out), nil
}
