package client

import (
	"errors"
	"net"
	"net/url"
)

func validateURLSyntax(urlString string) error {

	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return err
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	if u.Host == "" {
		return errors.New("invalid url provided, no host found")
	}

	if u.Scheme == "" {
		return errors.New("invalid url provided, no scheme found")
	}

	return nil
}

func hostExists(host string) bool {

	//retrieves the IP address(es) associated with the specified host
	_, err := net.LookupHost(host)
	if err != nil {
		return false
	}

	return true
}

func validateURLExists(urlString string) error {

	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	if u.Host == "" {
		return errors.New("invalid url provided, no host found")
	}

	if !hostExists(u.Host) {
		return errors.New("invalid url provided, host does not exist")
	}

	return nil
}

func validateURLLength(url string) error {

	if len(url) == 0 {
		return errors.New("url is empty")
	}

	if len(url) < len("http://a.co") {
		return errors.New("url is too short")
	}

	if len(url) > 2000 {
		return errors.New("url is too long")
	}

	return nil
}
