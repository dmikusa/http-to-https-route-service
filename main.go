// Copyright [yyyy] [name of copyright owner]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	DEFAULT_PORT            = "8080"
	CF_FORWARDED_URL_HEADER = "X-Cf-Forwarded-Url"
	X_FORWARDED_PROTO       = "X-Forwarded-Proto"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	log.SetOutput(os.Stdout)

	proxy := newProxy()
	redirectWrapper := newRedirectWrapperHandler(proxy)

	log.Fatal(http.ListenAndServe(":"+port, redirectWrapper))
}

func newRedirectWrapperHandler(wrappedHandler http.Handler) httpsRedirectWrapperHandler {
	return httpsRedirectWrapperHandler{
		code:           302,
		wrappedHandler: wrappedHandler,
	}
}

type httpsRedirectWrapperHandler struct {
	code           int
	wrappedHandler http.Handler
}

func (rh httpsRedirectWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forwardedURL := r.Header.Get(CF_FORWARDED_URL_HEADER)
	httpProto := r.Header.Get(X_FORWARDED_PROTO)

	if httpProto == "http" {
		log.Printf("Insecure request [%s] being redirected", forwardedURL)
		url, err := url.Parse(forwardedURL)
		if err != nil {
			log.Fatalln(err.Error())
		}
		url.Scheme = "https"

		http.Redirect(w, r, url.String(), rh.code)
	} else {
		rh.wrappedHandler.ServeHTTP(w, r)
	}
}

func newProxy() http.Handler {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)

			targetURL, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			req.URL = targetURL
			req.Host = targetURL.Host
		},
	}
	return reverseProxy
}
