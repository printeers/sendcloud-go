package method

import (
	"net/url"
	"strconv"

	"github.com/afosto/sendcloud-go"
)

type Client struct {
	apiKey    string
	apiSecret string
}

func New(apiKey string, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

type MethodOption func(*methodOptions)

type methodOptions struct {
	senderAddressID *int64
}

func WithSenderAddress(id int64) MethodOption {
	return func(o *methodOptions) {
		o.senderAddressID = &id
	}
}

// Get all shipment methods with optional filters
func (c *Client) GetMethods(opts ...MethodOption) ([]*sendcloud.Method, error) {
	options := &methodOptions{}
	for _, opt := range opts {
		opt(options)
	}

	smr := sendcloud.MethodListResponseContainer{}

	values := url.Values{}
	if options.senderAddressID != nil {
		values.Set("sender_address", strconv.FormatInt(*options.senderAddressID, 10))
	}

	reqURL := "/api/v2/shipping_methods"
	if len(values) > 0 {
		reqURL += "?" + values.Encode()
	}

	err := sendcloud.Request("GET", reqURL, nil, c.apiKey, c.apiSecret, &smr)
	if err != nil {
		return nil, err
	}
	return smr.GetResponse().([]*sendcloud.Method), nil
}

func (c *Client) GetReturnMethods() ([]*sendcloud.Method, error) {
	smr := sendcloud.MethodListResponseContainer{}
	err := sendcloud.Request("GET", "/api/v2/shipping_methods?is_return=true", nil, c.apiKey, c.apiSecret, &smr)
	if err != nil {
		return nil, err
	}
	return smr.GetResponse().([]*sendcloud.Method), nil
}

// Get a single method with optional sender_address filter
func (c *Client) GetMethod(id int64, opts ...MethodOption) (*sendcloud.Method, error) {
	options := &methodOptions{}
	for _, opt := range opts {
		opt(options)
	}

	values := url.Values{}
	if options.senderAddressID != nil {
		values.Set("sender_address", strconv.FormatInt(*options.senderAddressID, 10))
	}

	reqURL := "/api/v2/shipping_methods/" + strconv.Itoa(int(id))
	if len(values) > 0 {
		reqURL += "?" + values.Encode()
	}

	mr := sendcloud.MethodResponseContainer{}
	err := sendcloud.Request("GET", reqURL, nil, c.apiKey, c.apiSecret, &mr)
	if err != nil {
		return nil, err
	}
	return mr.GetResponse().(*sendcloud.Method), nil
}
