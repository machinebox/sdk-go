// Package videobox provides a client for accessing Videobox services.
package videobox

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/pkg/errors"
)

// Video represents a video.
type Video struct {
	ID                          string      `json:"id"`
	Status                      VideoStatus `json:"status"`
	Error                       string      `json:"error"`
	DownloadTotal               int64       `json:"downloadTotal,omitempty"`
	DownloadComplete            int64       `json:"downloadComplete,omitempty"`
	DownloadEstimatedCompletion *time.Time  `json:"downloadCompleteEstimate,omitempty"`
	FramesCount                 int         `json:"framesCount,omitempty"`
	FramesComplete              int         `json:"framesComplete"`
	LastFrameBase64             string      `json:"lastFrameBase64,omitempty"`
	MillisecondsComplete        int         `json:"millisecondsComplete"`
	Expires                     *time.Time  `json:"expires,omitempty"`
}

// Client is an HTTP client that can make requests to the box.
type Client struct {
	addr string

	// HTTPClient is the http.Client that will be used to
	// make requests.
	HTTPClient *http.Client
}

// New makes a new Client.
func New(addr string) *Client {
	return &Client{
		addr: addr,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Info gets the details about the box.
func (c *Client) Info() (*boxutil.Info, error) {
	var info boxutil.Info
	u, err := url.Parse(c.addr + "/info")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ErrVideobox represents an error from Videobox.
type ErrVideobox string

func (e ErrVideobox) Error() string {
	return "videobox: " + string(e)
}
