package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

var convExtCandidates = []string{".jpg", ".png", ".jpeg"}

func Do(req *http.Request) (*http.Response, error) {
	orgURL := req.URL
	var orgRes *http.Response
	pathExt := filepath.Ext(req.URL.Path)
	if pathExt == ".webp" {
		for _, cExt := range convExtCandidates {

			newPath := orgURL.Path[:len(orgURL.Path)-len(pathExt)] + cExt
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
	} else {
		orgRes, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("get origin failed")
		}
		return orgRes, nil
	}
}
