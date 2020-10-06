package client

import (
	"encoding/json"
	"io/ioutil"
	"laugh-tale/pkg/kozuki/types"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	URL        string
	HTTPClient *http.Client
}

// called by roger
// generate a container encryption key in kozuki database
func (c *Client) CreateKey() (types.Key, error) {
	reqURL := path.Join(c.URL, "key/create")
	resp, err := c.HTTPClient.Post(reqURL, "application/json", nil)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to create key with key server")
	}
	defer resp.Body.Close()
	return readKeyFromResponse(resp)
}

// called by roger
// generate a container encryption key in kozuki database
func (c *Client) BindKey(keyID, imageID string) (types.Key, error) {
	reqURL := path.Join(c.URL, "key/bind")
	req, err := http.NewRequest("POST", reqURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Form.Add("key_id", keyID)
	req.Form.Add("image_id", imageID)
	_, err = c.HTTPClient.Do(req)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to bind key with image id on key server")
	}
	return types.Key{}, nil
}

// called by ohara
// get runtime decryption key with certificate signed by image creator
func (c *Client) GetKey(imageID string) (types.Key, error) {
	reqURL := path.Join(c.URL, "key")
	reqURL = path.Join(reqURL, imageID)
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to get key from key server")
	}
	defer resp.Body.Close()
	return readKeyFromResponse(resp)
}

func readKeyFromResponse(resp *http.Response) (types.Key, error) {
	retJsonB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to read response body")
	}
	retKey := types.Key{}
	err = json.Unmarshal(retJsonB, &retKey)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to unmarshal response body")
	}
	return retKey, nil
}
