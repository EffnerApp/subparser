package dsb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DSB struct {
	User      string
	Pass      string
	AuthToken string
	Documents []DSBDocument
}

type DSBDocument struct {
	Id       string
	Date     string
	Title    string
	Detail   string
	Children []DSBDocument `json:"Childs"` // SIC!
}

const (
	BaseUrl    = "https://mobileapi.dsbcontrol.de"
	BundleId   = "de.heinekingmedia.dsbmobile"
	AppVersion = "36"
	OsVersion  = "36"
)

func NewDSB(user string, pass string) DSB {
	return DSB{
		User: user,
		Pass: pass,
	}
}

var ErrUnknownLogin = errors.New("unknown login error")
var ErrLoginFailed = errors.New("login failed")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrDataRequestFailed = errors.New("timetables could not be loaded")

// Login DSB presents an error as an empty auth code???
func (dsb *DSB) Login() error {
	// build the url to login with DSB
	url := fmt.Sprintf("%s/authid?bundleid=%s&appversion=%s&osversion=%s&pushid&user=%s&password=%s", BaseUrl, BundleId, AppVersion, OsVersion, dsb.User, dsb.Pass)

	// make the request
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return ErrUnknownLogin
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	bodyStr := string(body)
	if len(bodyStr) >= 2 {
		bodyStr = bodyStr[1 : len(bodyStr)-1]
	} else {
		bodyStr = ""
	}

	if bodyStr == "" {
		// Login failed
		return ErrLoginFailed
	}

	dsb.AuthToken = bodyStr

	return nil
}

func (dsb *DSB) LoadTimetables() error {
	// check if we obtained an auth token yet
	if dsb.AuthToken == "" {
		return ErrNotLoggedIn
	}

	// build url to get timetables
	url := fmt.Sprintf("%s/dsbtimetables?authid=%s", BaseUrl, dsb.AuthToken)

	// make the request
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return ErrDataRequestFailed
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &dsb.Documents)

	if err != nil {
		return err
	}
	return nil
}

func (document *DSBDocument) Download() ([]byte, error) {
	// make the request
	req, err := http.NewRequest(http.MethodGet, document.Detail, nil)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrDataRequestFailed
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}
