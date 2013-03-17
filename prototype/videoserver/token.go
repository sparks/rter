// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// rtER access tokens
//
// An access token is a simple mechanism to let a HTTP server independently verify
// authorization of calls to its resources. We assume an application server hands out
// tokens to authenticated clients who use them for contacting resource servers
// such as video ingest points or file servers.
//
// An access token allows a resource server to independently (without querying the
// issuing application server):
//   - verify the token was issued for the requested resource
//   - verify the token was issued to the requesting consumer (same IP address)
//   - confirm token freshness (lifetime)
//   - confirm token validity
//

// rtER access token parameters
//
// rter_resource: string (required)
//   The resource URI including scheme (http, https), authority (hostname), port
//   and path to the resource. Scheme and authority MUST be lowercase. Port number
//   MUST NOT be included if standard port 80 (http) or 443 (HTTPS) are used.
//
//   Example:  https2://rter.cim.mcgill.ca:6660/v1/ingest/1
//
// valid_until: string (required)
//   The token validity timestamp is an absolute time expressed as number of
//   seconds since Jan 1, 1970 00:00:00 GMT. The field MUST contain a stringified
//   representation of timestamp
//
//   Example: 137131200
//
// rter_signature: string (required)
//   A signature over all token fields (except the signature itself) is generated
//   using the following procedure:
//
//   1. Create `name=value` strings for each token parameter (consumer, resource,
//      and validity period) . The parameter name MUST be separated from the
//      corresponding value by an '=' character (ASCII code 61), even if the
//      value is empty.
//
//   2. Sort the parameter strings using lexicographical byte value ordering.
//      If two parameter names are equal, the sort order will be determined by
//      their values.
//
//   3. Concatenate the sorted parameter strings by an '&' character (ASCII code 38)
//      and URL encode the resulting string (using the RFC3986 percent-encoding
//      mechanism.
//
//   4. Create a message digest with HMAC and SHA265 where the concatenated
//      parameters are used as value and the server secret is used as key.
//
//   5. Generate a Base64 encoding of the digest and URL encode the result.

// HTTP Auth header
//
// Access tokens MUST be sent by clients in the HTTP Authorization header of each
// request as follows:
//
//   1. Create `name="value"` strings for each token parameter. Parameter name
//      MUST be separated from the corresponding value by an '=' character
//      (ASCII code 61), even if the value is empty. The parameter value MUST be
//      enclosed by an '"' character (ASCII code 34) and MAY be empty.
//
//   2. Sort the parameter strings using lexicographical byte value ordering.
//      If two parameter names are equal, the sort order will be determined by
//      their values.
//
//   3. Concatenate the sorted parameter strings by a comma (ASCII code 44).
//      The concatenated string MAY contain linear whitespace around the
//      separating comma characters.
//
//   4. Prefix the parameter string with "rtER " to indicate the access token
//      type.
//
// Example HTTP Request
//
// GET /resource/1 HTTP/1.1
//     Host: example.com
//     Authorization: rtER rter_resource="http://example.com/resource/1",
//						   rter_signature="cnRlcl9yZXNvdXJjZSUzRGh0dHAlM0ElMkYlMkZleGFtcGxlLmNvbSUyRnJlc291cmNlJTJGMSUyNnJ0ZXJfdmFsaWRfdW50aWwlM0QxMzYzNDY2NTg1b1aKbWSYzXXD6E%2B3QwN1VvhGsjfDxSUe%2FK%2FlH%2FZtlCg%3D"
//						   rter_valid_until="1363466585"

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AuthToken struct {
	Resource    string `json:"rter_resource"`
	Valid_until string `json:"rter_valid_until"`
	Signature   string `json:"rter_signature"`

	// internal variables, will not be exported
	lifetime int64  `json:"-"`
	consumer string `json:"-"`
}

func NewAuthToken() *AuthToken {
	var t AuthToken
	return &t
}

func (t *AuthToken) Sign(key string) error {
	sig_base := "rter_consumer=" + t.consumer + "&rter_resource=" + t.Resource + "&rter_valid_until=" + t.Valid_until
	hmac := hmac.New(sha256.New, bytes.NewBufferString(key).Bytes())
	t.Signature = url.QueryEscape(base64.StdEncoding.EncodeToString(hmac.Sum(bytes.NewBufferString(url.QueryEscape(sig_base)).Bytes())))
	return nil
}

func (t *AuthToken) VerifySignature(key string) error {
	sig_base := "rter_consumer=" + t.consumer + "&rter_resource=" + t.Resource + "&rter_valid_until=" + t.Valid_until
	hmac := hmac.New(sha256.New, bytes.NewBufferString(key).Bytes())
	checksig := url.QueryEscape(base64.StdEncoding.EncodeToString(hmac.Sum(bytes.NewBufferString(url.QueryEscape(sig_base)).Bytes())))

	if t.Signature == checksig {
		return nil
	}
	return errors.New("signature verification failed")
}

func (t *AuthToken) VerifyLifetime() error {

	now := time.Now().UTC()
	eol_sec, err := strconv.ParseInt(t.Valid_until, 10, 64)
	if err != nil {
		return err
	}

	if eol := time.Unix(eol_sec, 0); eol.After(now) {
		return nil
	}

	return errors.New("Auth: token expired")
}

func NewFromHttpRequest(r *http.Request) (*AuthToken, error) {

	// find our rtER auth header, assuming there may be more than one auth header
	var s string
	for _, h := range r.Header["Authorization"] {
		if strings.HasPrefix(h, "rtER") {
			s = strings.Replace(h, "rtER", "", 1)
			break
		}
	}

	// auth string must be present
	if s == "" {
		return nil, errors.New("Auth: Authorization header not found")
	}

	// split the key-value pairs into a map
	m := make(map[string]string)
	pairs := strings.Split(strings.TrimSpace(s), ",")
	for _, p := range pairs {
		// split parameters into key and value before insterting them into map
		kv := strings.Split(strings.TrimSpace(p), "=")
		// remove '"' around values
		m[kv[0]], _ = strconv.Unquote(kv[1])
	}

	t := NewAuthToken()
	t.Resource = m["rter_resource"]
	t.Valid_until = m["rter_valid_until"]
	t.Signature = m["rter_signature"]
	t.lifetime, _ = strconv.ParseInt(t.Valid_until, 10, 64)

	// check if mandatory values are present and sane

	// URI must not be empty and have a valid structure
	if t.Resource == "" {
		return t, errors.New("Auth: resource field missing")
	}
	url, err := url.Parse(t.Resource)
	if err != nil {
		return t, errors.New("Auth: resource parse error")
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return t, errors.New("Auth: resource scheme invalid")
	}

	// lifetime must not be empty and > 0
	if t.Valid_until == "" {
		return t, errors.New("Auth: timestamp field missing")
	}
	if t.lifetime <= 0 {
		return t, errors.New("Auth: timestamp must be > 0")
	}

	// signature must not be empty
	if t.Signature == "" {
		return t, errors.New("Auth: signature field missing")
	}

	// add resource consumer (sender of the request)
	t.consumer = strings.Split(r.RemoteAddr, ":")[0]

	return t, nil
}

func (t *AuthToken) String() string {

	return fmt.Sprintf("rtER rter_resource=\"%s\", rter_signature=\"%s\", rter_valid_until=\"%s\"",
		t.Resource, t.Signature, t.Valid_until)
}

func (t *AuthToken) Json() string {

	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func AuthenticateRequest(r *http.Request, key string) *ServerError {

	// parse token from HTTP request Authorization header
	t, err := NewFromHttpRequest(r)
	if err != nil {
		if t == nil {
			return ServerErrorAuthTokenRequired
		} else {
			return ServerErrorAuthTokenInvalid
		}
	}

	// tokens issued for a resource are valid for all sub-resources
	if !strings.HasPrefix(r.URL.String(), t.Resource) {
		return ServerErrorAuthUrlMismatch
	}

	// tokens are only valid up to a specified liftiem
	if t.VerifyLifetime() != nil {
		return ServerErrorAuthTokenExpired
	}

	// tokens are only valid when their signature is correct
	if t.VerifySignature(key) != nil {
		return ServerErrorAuthBadSignature
	}

	return nil
}

func GenerateAuthToken(uri, consumer string, valid time.Duration, key string) (*AuthToken, error) {

	// check preconditions
	if uri == "" {
		return nil, errors.New("Auth: empty URI not allowed")
	}

	url, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("Auth: resource parse error")
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, errors.New("Auth: resource scheme invalid")
	}

	if consumer == "" {
		return nil, errors.New("Auth: empty consumer not allowed")
	}
	if key == "" {
		return nil, errors.New("Auth: empty key not allowed")
	}
	if valid <= 0 {
		return nil, errors.New("Auth: invalid lifetime")
	}

	t := NewAuthToken()

	// use the passed URI as resource id
	t.Resource = uri
	t.consumer = consumer

	// generate valid_until timestamp
	t.lifetime = time.Now().UTC().Add(valid).Unix()
	t.Valid_until = strconv.FormatInt(t.lifetime, 10)

	// sign the token with our server key
	if err := t.Sign(key); err != nil {
		return nil, err
	}

	return t, nil
}
