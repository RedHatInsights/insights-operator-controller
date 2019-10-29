package logging

import (
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk"
)

type Client struct {
	ClientImpl *splunk.Client
}

func NewClient(address string, token string, source string, source_type string, index string) Client {
	url := address + "/services/collector/raw"
	splunk := splunk.NewClient(nil, url, token, source, source_type, index)
	return Client{ClientImpl: splunk}
}

func (client Client) Log(key string, value string) error {
	err := client.ClientImpl.Log(
		map[string]string{key: value})
	return err
}

func (client Client) LogWithTime(time int64, key string, value string) error {
	err := client.ClientImpl.LogWithTime(
		time,
		map[string]string{key: value})
	return err
}
