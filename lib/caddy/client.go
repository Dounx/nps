package caddy

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const reverseProxyJSONTemplate = `
{
    "@id": "",
    "handle": [
        {
            "handler": "subroute",
            "routes": [
                {
                    "handle": [
                        {
                            "handler": "reverse_proxy",
                            "upstreams": [
                                {
                                    "dial": "127.0.0.1:3000"
                                }
                            ]
                        }
                    ],
                    "match": [
                        {
                            "path": [
                                "/*"
                            ]
                        }
                    ]
                }
            ]
        }
    ],
    "match": [
        {
            "host": [
                "example.com"
            ]
        }
    ],
    "terminal": true
}
`

type ReverseConfig struct {
	ID           int64
	UpstreamPath string
	MatchHost    string
	MatchPath    string
}

func (config *ReverseConfig) ToJSON() (string, error) {
	value, err := sjson.Set(reverseProxyJSONTemplate, "@id", config.ID)
	if err != nil {
		return "", err
	}

	value, err = sjson.Set(
		value, "handle.0.routes.0.handle.0.upstreams.0.dial",
		config.UpstreamPath,
	)
	if err != nil {
		return "", err
	}

	value, err = sjson.Set(value, "handle.0.routes.0.match.0.path.0", config.MatchPath)
	if err != nil {
		return "", err
	}

	value, err = sjson.Set(value, "match.0.host.0", config.MatchHost)
	if err != nil {
		return "", err
	}

	return value, nil
}

type Client struct {
	Host string
	Port string
}

func (client *Client) LoadConfig(config string) error {
	url := fmt.Sprintf("http://%s:%s/load", client.Host, client.Port)

	return client.do("POST", url, strings.NewReader(config))
}

func (client *Client) GetConfig(path string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/config%s", client.Host, client.Port, path)

	return client.get(url)
}

func (client *Client) Stop() error {
	url := fmt.Sprintf("http://%s:%s/stop", client.Host, client.Port)

	return client.do("POST", url, nil)
}

func (client *Client) GetReverseProxyList() ([]*ReverseConfig, error) {
	config, err := client.GetConfig("/apps/http/servers/srv0/routes")
	if err != nil {
		return nil, err
	}

	var reverseConfigArray []*ReverseConfig
	configArray := gjson.Get(config, "@this").Array()

	for _, c := range configArray {
		routeConfig := gjson.Get(c.String(), "handle.0.routes.0").String()

		reverseConfigArray =
			append(reverseConfigArray,
				&ReverseConfig{
					ID:           gjson.Get(c.String(), "@id").Int(),
					UpstreamPath: gjson.Get(routeConfig, "handle.0.upstreams.0.dial").String(),
					MatchPath:    gjson.Get(routeConfig, "match.0.path.0").String(),
					MatchHost:    gjson.Get(c.String(), "match.0.host.0").String(),
				})
	}

	return reverseConfigArray, nil
}

func (client *Client) GetReverseProxyLastID() (int64, error) {
	list, err := client.GetReverseProxyList()
	if err != nil {
		return -1, err
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID > list[j].ID
	})

	return list[0].ID, nil
}

func (client *Client) GetReverseProxy(id int64) (*ReverseConfig, error) {
	url := fmt.Sprintf("http://%s:%s/id/%v", client.Host, client.Port, id)

	config, err := client.get(url)
	if err != nil {
		return nil, err
	}

	routeConfig := gjson.Get(config, "handle.0.routes.0").String()
	return &ReverseConfig{
		ID:           gjson.Get(config, "@id").Int(),
		UpstreamPath: gjson.Get(routeConfig, "handle.0.upstreams.0.dial").String(),
		MatchPath:    gjson.Get(routeConfig, "match.0.path.0").String(),
		MatchHost:    gjson.Get(config, "match.0.host.0").String(),
	}, nil
}

func (client *Client) AddReverseProxy(config *ReverseConfig) error {
	url := fmt.Sprintf("http://%s:%s/config/apps/http/servers/srv0/routes/", client.Host, client.Port)

	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	return client.do("POST", url, strings.NewReader(configJSON))
}

func (client *Client) UpdateReverseProxy(config *ReverseConfig) error {
	url := fmt.Sprintf("http://%s:%s/id/%v", client.Host, client.Port, config.ID)

	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	return client.do("PATCH", url, strings.NewReader(configJSON))
}

func (client *Client) DeleteReverseProxy(id int64) error {
	url := fmt.Sprintf("http://%s:%s/id/%v", client.Host, client.Port, id)
	return client.do("DELETE", url, nil)
}

func (client *Client) get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (client *Client) do(action string, url string, body io.Reader) error {
	req, err := http.NewRequest(action, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	return nil
}
