package videobox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Check starts processing the video in the Reader.
// Videobox is asynchronous, you must use Status to check when a
// video processing operation has completed before using Results to
// get the results.
func (c *Client) Check(video io.Reader, options *CheckOptions) (*Video, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, video)
	if err != nil {
		return nil, err
	}
	if err := options.apply(w.WriteField); err != nil {
		return nil, errors.Wrap(err, "setting options")
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	u, err := url.Parse(c.addr + "/videobox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckURL starts processing the video at the specified URL.
// See Check for more information.
func (c *Client) CheckURL(videoURL *url.URL, options *CheckOptions) (*Video, error) {
	u, err := url.Parse(c.addr + "/videobox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	if !videoURL.IsAbs() {
		return nil, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", videoURL.String())
	formset := func(key, value string) error {
		form.Set(key, value)
		return nil
	}
	if err := options.apply(formset); err != nil {
		return nil, errors.Wrap(err, "setting options")
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckBase64 starts processing the video from the base64 encoded data string.
// See Check for more information.
func (c *Client) CheckBase64(data string, options *CheckOptions) (*Video, error) {
	u, err := url.Parse(c.addr + "/videobox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	formset := func(key, value string) error {
		form.Set(key, value)
		return nil
	}
	if err := options.apply(formset); err != nil {
		return nil, errors.Wrap(err, "setting options")
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckOptions are additional options that control
// the behaviour of Videobox when processing videos.
type CheckOptions struct {
	fields map[string]string
}

// NewCheckOptions makes a new CheckOptions object.
func NewCheckOptions() *CheckOptions {
	return &CheckOptions{
		fields: make(map[string]string),
	}
}

// ResultsDuration sets the duration results should be kept in Videobox
// before being garbage collected.
func (o *CheckOptions) ResultsDuration(duration time.Duration) {
	o.fields["resultsDuration"] = duration.String()
}

// SkipFrames sets the number of frames to skip between extractions.
func (o *CheckOptions) SkipFrames(frames int) {
	o.fields["skipframes"] = strconv.Itoa(frames)
}

// SkipSeconds sets the number of seconds to skip between frame extractions.
func (o *CheckOptions) SkipSeconds(seconds int) {
	o.fields["skipseconds"] = strconv.Itoa(seconds)
}

// FrameWidth sets the width of the frame to extract.
func (o *CheckOptions) FrameWidth(width int) {
	o.fields["frameWidth"] = strconv.Itoa(width)
}

// FrameHeight sets the height of the frame to extract.
func (o *CheckOptions) FrameHeight(height int) {
	o.fields["frameHeight"] = strconv.Itoa(height)
}

// FaceboxThreshold sets the minimum confidence threshold of Facebox
// matches required for it to be included in the results.
func (o *CheckOptions) FaceboxThreshold(v float64) {
	o.fields["faceboxThreshold"] = strconv.FormatFloat(v, 'f', -1, 64)
}

// TagboxIncludeAll includes all tags in the results.
func (o *CheckOptions) TagboxIncludeAll() {
	o.fields["tagboxInclude"] = "all"
}

// TagboxIncludeCustom includes only custom tags in the results.
func (o *CheckOptions) TagboxIncludeCustom() {
	o.fields["tagboxInclude"] = "custom"
}

// TagboxThreshold sets the minimum confidence threshold of Tagbox
// matches required for it to be included in the results.
func (o *CheckOptions) TagboxThreshold(v float64) {
	o.fields["tagboxThreshold"] = strconv.FormatFloat(v, 'f', -1, 64)
}

// NudeboxThreshold sets the minimum confidence threshold of Nudebox
// matches required for it to be included in the results.
func (o *CheckOptions) NudeboxThreshold(v float64) {
	o.fields["nudeboxThreshold"] = strconv.FormatFloat(v, 'f', -1, 64)
}

// apply calls writeField for each field.
// If o is nil, apply is noop.
func (o *CheckOptions) apply(writeField func(key, value string) error) error {
	if o == nil {
		return nil
	}
	for k, v := range o.fields {
		if err := writeField(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) parseCheckResponse(r io.Reader) (*Video, error) {
	var checkResponse struct {
		Success bool
		Error   string
		ID      string
	}
	if err := json.NewDecoder(r).Decode(&checkResponse); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !checkResponse.Success {
		return nil, ErrVideobox(checkResponse.Error)
	}
	v := &Video{
		ID: checkResponse.ID,
	}
	return v, nil
}
