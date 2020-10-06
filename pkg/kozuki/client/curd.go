package client

import (
	"bytes"
	"encoding/json"
	"laugh-tale/pkg/kozuki/types"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type CURLClient struct {
	URL        string
	HTTPClient *http.Client
}

func (crud *CURLClient) Create(k types.Key) (types.Key, error) {
	return crud.doHTTPRequest(k, "POST")
}

func (crud *CURLClient) Retrieve(k types.Key) (types.Key, error) {
	if k.ID.String() == "" && k.ImageID == "" {
		return types.Key{}, errors.New("either of key id and image id must be set")
	}
	reqURL := path.Join(crud.URL, "key")
	urlObj, err := url.Parse(reqURL)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to parse url")
	}
	params := url.Values{}
	if k.ID.String() != "" {
		params.Add("key_id", k.ID.String())
	}
	if k.ImageID != "" {
		params.Add("image_id", k.ImageID)
	}
	urlObj.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", urlObj.String(), nil)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to create GET request")
	}
	resp, err := crud.HTTPClient.Do(req)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to perform GET request")
	}
	return readKeyFromResponse(resp)
}

func (crud *CURLClient) Update(k types.Key) (types.Key, error) {
	return crud.doHTTPRequest(k, "PATCH")
}

func (crud *CURLClient) Delete(k types.Key) error {
	if k.ID.String() == "" {
		return errors.New("key id must be set")
	}
	reqURL := path.Join(crud.URL, "key")
	urlObj, err := url.Parse(reqURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse url")
	}
	params := url.Values{}
	params.Add("key_id", k.ID.String())
	urlObj.RawQuery = params.Encode()
	req, err := http.NewRequest("DELETE", urlObj.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create DELETE request")
	}
	_, err = crud.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform DELETE request")
	}
	return nil
}

func (crud *CURLClient) doHTTPRequest(k types.Key, m string) (types.Key, error) {
	reqURL := path.Join(crud.URL, "key")
	keyB, err := json.Marshal(k)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to marshal request body")
	}
	req, err := http.NewRequest(m, reqURL, bytes.NewBuffer(keyB))
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := crud.HTTPClient.Do(req)
	if err != nil {
		return types.Key{}, errors.Wrap(err, "failed to perform http "+m+" request")
	}
	return readKeyFromResponse(resp)
}
