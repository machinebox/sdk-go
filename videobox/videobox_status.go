package videobox

import (
	"net/http"
	"net/url"

	"github.com/machinebox/sdk-go/internal/mbhttp"
	"github.com/pkg/errors"
)

// VideoStatus holds the status of a video processing operation.
type VideoStatus string

const (
	// StatusPending indicates that a video operation is pending.
	StatusPending VideoStatus = "pending"
	// StatusDownloading indicates that a video file is being downloaded.
	StatusDownloading VideoStatus = "downloading"
	// StatusProcessing indicates that a video is being processed.
	StatusProcessing VideoStatus = "processing"
	// StatusComplete indicates that a video operation has finished.
	StatusComplete VideoStatus = "complete"
	// StatusFailed indicates that a video operation has failed.
	StatusFailed VideoStatus = "failed"
)

// Status gets the status of a video operation.
func (c *Client) Status(id string) (*Video, error) {
	u, err := url.Parse(c.addr + "/videobox/status/" + id)
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
	var video Video
	_, err = mbhttp.New("videobox", c.HTTPClient).DoUnmarshal(req, &video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}
