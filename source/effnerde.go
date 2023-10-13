package source

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type EffnerDESource struct {
	Password string
}

func (effnerSrc *EffnerDESource) Load() (string, error) {
	// TODO Move this somewhere else.
	form := url.Values{}
	form.Add("post_password", effnerSrc.Password)

	// load the html from effner.de
	req, err := http.NewRequest("POST", "https://effner.de/wp-login.php?action=postpass", strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Referer", "https://effner.de/service/vertretungsplan/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	jar, err := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: jar,
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	link := res.Header.Get("Link")
	linkParts := strings.Split(link, ",")

	final := strings.Split(linkParts[1], ";")[0]
	final = final[2 : len(final)-1]

	req, err = http.NewRequest("GET", final, nil)
	res, err = httpClient.Do(req)

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
