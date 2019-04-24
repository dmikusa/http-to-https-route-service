# HTTPS Redirect Route Service

A small route service which can be used to redirect HTTP traffic to HTTPS.

## Getting Started

- Download this repository or `git clone` it.
- Change directories into the repository
- Edit the manifest.yml file, minimally you'll need to change the route.
- Run `cf push http-to-https-rs`.
- Create a user-provided route service and bind that to the route or routes of your choosing. [See docs](http://docs.cloudfoundry.org/services/route-services.html#user-provided).
- Run `cf logs http-to-https-rs`. You should see log messages like `Insecure request [<url>] being redirected` when HTTP requests are redirected to HTTPS.

## How It Works

Requests are delivered to this route service. The route service examines the `X-Forwarded-Proto` header for each request and if that is set to `http` then it will return a 302 redirect in response to that request. The redirect will have the `Location` header set to the URL indicated in `X-Cf-Forwarded-Url` but it will make sure that the scheme is set to `https`. This should cause well behaved clients to resend their request but using HTTPS.

Requests that have the `X-Forwarded-Proto` set to something other than `http` will be proxied through to `X-Cf-Forwarded-Url` as these are already using HTTPS.
