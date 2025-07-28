package webrequest

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type userAgentTransport struct {
	Transport http.RoundTripper
	UserAgent string
}

type LoginResponse struct {
	IntendedRoute string `json:"intended_route"`
}

func AuthenticateWithCookies(cookies string) (*http.Client, error) {
	client, err := newClient(cookies)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get("https://academy.hackthebox.com/dashboard")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Authentication Failed, refresh your cookies and try again!")
	}

	return client, nil
}

func (ua *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", ua.UserAgent)
	}
	return ua.Transport.RoundTrip(req)
}

func newClient(cookies string) (*http.Client, error) {
	// For proxy debugging
	//proxy, _ := url.Parse("http://localhost:8080")
	//transport := &userAgentTransport{
	//	Transport: &http.Transport{
	//		Proxy:           http.ProxyURL(proxy),
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//	UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
	//}
	transport := &userAgentTransport{
		Transport: http.DefaultTransport,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	if cookies != "" {
		addCookiesToJar(jar, cookies)
	}

	return client, nil
}

func addCookiesToJar(jar *cookiejar.Jar, cookies string) {
	cookiePairs := strings.Split(cookies, ";")
	cookieList := []*http.Cookie{}

	for _, pair := range cookiePairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookieList = append(cookieList, &http.Cookie{
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			})
		}
	}

	u, _ := url.Parse("https://academy.hackthebox.com")
	jar.SetCookies(u, cookieList)
}

func getXSRFToken(client *http.Client, urlStr string) (string, error) {
	u, _ := url.Parse(urlStr)
	cookies := client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			rawToken, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return "", err
			}
			return rawToken, nil
		}
	}

	return "", nil
}
