package api

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"regexp"
)

// ForgejoGateway acts as a proxy to Forgejo
type ForgejoGateway struct {
	proxy *httputil.ReverseProxy
}

var forgejoPathRegexp = regexp.MustCompile("^/forgejo/?")
var forgejoAllowedRegexp = regexp.MustCompile("^/forgejo/((git|contents|pulls|branches|merges|statuses|compare|commits)/?|(issues/(\\d+)/labels))")

func NewForgejoGateway() *ForgejoGateway {
	return &ForgejoGateway{
		proxy: &httputil.ReverseProxy{
			Director:     forgejoDirector,
			Transport:    &ForgejoTransport{},
			ErrorHandler: proxyErrorHandler,
		},
	}
}

func forgejoDirector(r *http.Request) {
	ctx := r.Context()
	target := getProxyTarget(ctx)
	accessToken := getAccessToken(ctx)

	targetQuery := target.RawQuery
	r.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.URL.Host = target.Host
	r.URL.Path = singleJoiningSlash(target.Path, forgejoPathRegexp.ReplaceAllString(r.URL.Path, "/"))
	if targetQuery == "" || r.URL.RawQuery == "" {
		r.URL.RawQuery = targetQuery + r.URL.RawQuery
	} else {
		r.URL.RawQuery = targetQuery + "&" + r.URL.RawQuery
	}
	if _, ok := r.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		r.Header.Set("User-Agent", "")
	}
	if r.Method != http.MethodOptions {
		r.Header.Set("Authorization", "Bearer "+accessToken)
	}

	log := getLogEntry(r)
	log.Infof("Proxying to Forgejo: %v", r.URL.String())
}

func (fj *ForgejoGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := getConfig(ctx)
	if config == nil || config.Forgejo.AccessToken == "" {
		handleError(notFoundError("No Forgejo Settings Configured"), w, r)
		return
	}

	if err := fj.authenticate(w, r); err != nil {
		handleError(unauthorizedError("%s", err.Error()), w, r)
		return
	}

	endpoint := config.Forgejo.Endpoint
	apiURL := singleJoiningSlash(endpoint, "/api/v1/repos/"+config.Forgejo.Repo)
	target, err := url.Parse(apiURL)
	if err != nil {
		handleError(internalServerError("Unable to process Forgejo endpoint"), w, r)
		return
	}
	ctx = withProxyTarget(ctx, target)
	ctx = withAccessToken(ctx, config.Forgejo.AccessToken)
	fj.proxy.ServeHTTP(w, r.WithContext(ctx))
}

func (fj *ForgejoGateway) authenticate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims := getClaims(ctx)
	config := getConfig(ctx)

	if claims == nil {
		return errors.New("Access to endpoint not allowed: no claims found in Bearer token")
	}

	if !forgejoAllowedRegexp.MatchString(r.URL.Path) {
		return errors.New("Access to endpoint not allowed: this part of Forgejo's API has been restricted")
	}

	if len(config.AcceptContentPaths) > 0 {
		prefix := "/forgejo/contents/"
		if strings.HasPrefix(r.URL.Path, prefix) {
			contentPath := strings.TrimPrefix(r.URL.Path, prefix)
			if !isPathAllowed(config.AcceptContentPaths, contentPath) {
				return errors.New("Access to endpoint not allowed: content path is restricted")
			}
		}
	}

	if len(config.Roles) == 0 {
		return nil
	}

	roles, ok := claims.AppMetaData["roles"]
	if ok {
		roleStrings, _ := roles.([]interface{})
		for _, data := range roleStrings {
			role, _ := data.(string)
			for _, adminRole := range config.Roles {
				if role == adminRole {
					return nil
				}
			}
		}
	}

	return errors.New("Access to endpoint not allowed: your role doesn't allow access")
}

type ForgejoTransport struct{}

func (t *ForgejoTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err == nil {
		// remove CORS headers from Forgejo and use our own
		resp.Header.Del("Access-Control-Allow-Origin")
	}
	return resp, err
}
