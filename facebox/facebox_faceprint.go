package facebox

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/machinebox/sdk-go/internal/mbhttp"
	"github.com/pkg/errors"
)

// CompareFaceprints returns the confidence of the comparsion between the target faceprint
// and each of faceprint of the slice of candidates and returns an array of confidence in the same order
// of the candidates
func (c *Client) CompareFaceprints(target string, faceprintCandidates []string) ([]float64, error) {
	if target == "" {
		return nil, errors.New("target can not be empty")
	}
	u, err := url.Parse(c.addr + "/facebox/faceprint/compare")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}

	request := compareFaceprintRequest{
		Target:     target,
		Faceprints: faceprintCandidates,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var compareFaceprintsResponse struct {
		Confidences []float64
	}
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &compareFaceprintsResponse)
	if err != nil {
		return nil, err
	}
	return compareFaceprintsResponse.Confidences, nil
}

type compareFaceprintRequest struct {
	Faceprints []string `json:"faceprints"`
	Target     string   `json:"target"`
}

type checkFaceprintRequest struct {
	Faceprints []string `json:"faceprints"`
}

// CheckFaceprints checks the list of faceprints to see if they
// match any known faces.
func (c *Client) CheckFaceprints(faceprints []string) ([]Face, error) {
	if len(faceprints) == 0 {
		return nil, errors.New("faceprints can not be empty")
	}
	u, err := url.Parse(c.addr + "/facebox/faceprint/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	request := checkFaceprintRequest{
		Faceprints: faceprints,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var checkResponse struct {
		Faceprints []Face
	}
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return nil, err
	}
	return checkResponse.Faceprints, nil
}
