package angulardetector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"

	"github.com/grafana/grafana/pkg/plugins/log"
)

const (
	// gcomAngularPatternsPath is the relative path to the GCOM API handler that returns angular detection patterns.
	gcomAngularPatternsPath = "/api/plugins/angular_patterns"
)

var _ DetectorsProvider = &GCOMDetectorsProvider{}

// GCOMDetectorsProvider is a DetectorsProvider which fetches patterns from GCOM.
type GCOMDetectorsProvider struct {
	log log.Logger

	httpClient *http.Client

	baseURL string
}

// NewGCOMDetectorsProvider returns a new GCOMDetectorsProvider.
// baseURL is the GCOM base url, without /api and without a trailing slash (e.g.: https://grafana.com)
func NewGCOMDetectorsProvider(baseURL string) (DetectorsProvider, error) {
	cl, err := httpclient.New()
	if err != nil {
		return nil, fmt.Errorf("httpclient new: %w", err)
	}
	return &GCOMDetectorsProvider{
		log:        log.New("plugins.angulardetector.gcom"),
		baseURL:    baseURL,
		httpClient: cl,
	}, nil
}

// ProvideDetectors gets the dynamic detectors from the remote source.
// If an error occurs, the function fails silently by logging an error and it returns nil.
func (p *GCOMDetectorsProvider) ProvideDetectors(ctx context.Context) []Detector {
	patterns, err := p.fetch(ctx)
	if err != nil {
		p.log.Warn("Could not fetch remote angular patterns", "error", err)
		return nil
	}
	detectors, err := patterns.detectors()
	if err != nil {
		p.log.Warn("Could not convert angular patterns to detectors", "error", err)
		return nil
	}
	return detectors
}

// fetch fetches the angular patterns from GCOM and returns them as gcomPatterns.
// Call detectors() on the returned value to get the corresponding detectors.
func (p *GCOMDetectorsProvider) fetch(ctx context.Context) (gcomPatterns, error) {
	st := time.Now()

	reqURL, err := url.JoinPath(p.baseURL, gcomAngularPatternsPath)
	if err != nil {
		return nil, fmt.Errorf("url joinpath: %w", err)
	}

	p.log.Debug("Fetching dynamic angular detection patterns", "url", reqURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			p.log.Error("response body close error", "error", err)
		}
	}()
	var out gcomPatterns
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	p.log.Debug("Fetched dynamic angular detection patterns", "patterns", len(out), "duration", time.Since(st))
	return out, nil
}

// gcomPatternType is a pattern type returned by the GCOM API.
type gcomPatternType string

const (
	gcomPatternTypeContains gcomPatternType = "contains"
	gcomPatternTypeRegex    gcomPatternType = "regex"
)

// gcomPattern is an Angular detection pattern returned by the GCOM API.
type gcomPattern struct {
	Name    string
	Pattern string
	Type    gcomPatternType
}

// errUnknownPatternType is returned when a pattern type is not known.
var errUnknownPatternType = errors.New("unknown pattern type")

// Detector converts a gcomPattern into a Detector, based on its Type.
// If a pattern type is unknown, it returns an error wrapping errUnknownPatternType.
func (p *gcomPattern) detector() (Detector, error) {
	switch p.Type {
	case gcomPatternTypeContains:
		return &ContainsBytesDetector{Pattern: []byte(p.Pattern)}, nil
	case gcomPatternTypeRegex:
		re, err := regexp.Compile(p.Pattern)
		if err != nil {
			return nil, fmt.Errorf("%q regexp compile: %w", p.Pattern, err)
		}
		return &RegexDetector{Regex: re}, nil
	}
	return nil, fmt.Errorf("%q: %w", p.Type, errUnknownPatternType)
}

// gcomPatterns is a slice of gcomPattern s.
type gcomPatterns []gcomPattern

// detectors converts the slice of gcomPattern s into a slice of detectors, by calling Detector() on each gcomPattern.
func (p gcomPatterns) detectors() ([]Detector, error) {
	var finalErr error
	detectors := make([]Detector, 0, len(p))
	for _, pattern := range p {
		d, err := pattern.detector()
		if err != nil {
			// Fail silently in case of an errUnknownPatternType.
			// This allows us to introduce new pattern types without breaking old Grafana versions
			if !errors.Is(err, errUnknownPatternType) {
				// Other error, do not ignore it
				finalErr = errors.Join(finalErr, err)
			}
			continue
		}
		detectors = append(detectors, d)
	}
	if finalErr != nil {
		return nil, finalErr
	}
	return detectors, nil
}
