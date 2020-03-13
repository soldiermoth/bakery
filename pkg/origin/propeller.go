package origin

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/propeller-client-go/pkg/client"
)

//Propeller struct holds basic config of a Propeller Channel
type Propeller struct {
	URL       string
	OrgID     string
	ChannelID string
}

//GetPlaybackURL will retrieve url
func (p *Propeller) GetPlaybackURL() string {
	return p.URL
}

//FetchManifest will grab manifest contents of configured origin
func (p *Propeller) FetchManifest(c config.Config) (string, error) {
	return fetch(c, p.URL)
}

//NewPropeller returns a propeller struct
func NewPropeller(c config.Config, orgID string, channelID string, renditionURL string) (*Propeller, error) {
	var propellerURL string
	var err error

	switch renditionURL {
	case "": //fetch master
		propellerURL, err = getPropellerChannelURL(c.PropellerHost, orgID, channelID)
		if err != nil {
			return &Propeller{}, fmt.Errorf("fetching propeller channel: %w", err)
		}
	default:
		propellerURL, err = decodeRenditionURL(renditionURL)
		if err != nil {
			return &Propeller{}, fmt.Errorf("fetching propeller rendition: %w", err)
		}
	}

	return &Propeller{
		URL:       propellerURL,
		OrgID:     orgID,
		ChannelID: channelID,
	}, nil
}

func decodeRenditionURL(rendition string) (string, error) {
	rendition = strings.TrimSuffix(rendition, ".m3u8")
	url, err := base64.StdEncoding.DecodeString(rendition)
	if err != nil {
		return "", fmt.Errorf("decoding rendition: %w", err)
	}

	return string(url), nil
}

func getPropellerChannelURL(host string, orgID string, channelID string) (string, error) {
	pURL, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("parsing propeller host url: %w", err)
	}
	p := client.NewClient(pURL)

	channel, err := p.GetChannel(orgID, channelID)
	if err != nil {
		return "", fmt.Errorf("fetching channel from propeller: %w", err)
	}

	manifestURL, err := channel.URL()
	if err != nil {
		return "", fmt.Errorf("reading url from propeller channel: %w", err)
	}

	return manifestURL.String(), nil
}
