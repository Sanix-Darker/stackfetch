package cloudmeta

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Info struct {
	Provider     string `json:"provider,omitempty"`
	InstanceID   string `json:"instanceId,omitempty"`
	InstanceType string `json:"instanceType,omitempty"`
	Region       string `json:"region,omitempty"`
}

var hc = &http.Client{
	Timeout: 600 * time.Millisecond,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{Timeout: 300 * time.Millisecond}).DialContext,
	},
}

func Collect() Info {
	if v, ok := aws(); ok {
		return v
	}
	if v, ok := gcp(); ok {
		return v
	}
	if v, ok := azure(); ok {
		return v
	}
	return Info{}
}

func aws() (Info, bool) {
	id, err := fetch("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return Info{}, false
	}
	itype, _ := fetch("http://169.254.169.254/latest/meta-data/instance-type")
	reg, _ := fetch("http://169.254.169.254/latest/meta-data/placement/region")
	return Info{"AWS", id, itype, reg}, true
}

func gcp() (Info, bool) {
	h := map[string]string{"Metadata-Flavor": "Google"}
	id, err := fetchH("http://metadata.google.internal/computeMetadata/v1/instance/id", h)
	if err != nil {
		return Info{}, false
	}
	mt, _ := fetchH("http://metadata.google.internal/computeMetadata/v1/instance/machine-type", h)
	z, _ := fetchH("http://metadata.google.internal/computeMetadata/v1/instance/zone", h)
	mt = last(mt)
	z = last(z)
	return Info{"GCP", id, mt, z}, true
}

func azure() (Info, bool) {
	req, _ := http.NewRequest("GET",
		"http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
	req.Header.Set("Metadata", "true")
	resp, err := hc.Do(req)
	if err != nil {
		return Info{}, false
	}
	defer resp.Body.Close()
	var d struct {
		Compute struct {
			VmId     string `json:"vmId"`
			VmSize   string `json:"vmSize"`
			Location string `json:"location"`
		} `json:"compute"`
	}
	if json.NewDecoder(resp.Body).Decode(&d) != nil {
		return Info{}, false
	}
	return Info{"Azure", d.Compute.VmId, d.Compute.VmSize, d.Compute.Location}, true
}

func fetchH(url string, hdr map[string]string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(b)), nil
}

func fetch(url string) (string, error) { return fetchH(url, nil) }

func last(s string) string {
	if i := strings.LastIndexByte(s, '/'); i >= 0 {
		return s[i+1:]
	}
	return s
}
