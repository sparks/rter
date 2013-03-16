// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// Test cases for token authorization

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

const (
	TEST_TOKEN_URI          string        = "http://example.com/resource/1"
	TEST_TOKEN_SECRET       string        = "1122AABBCCDDEEFF"
	TEST_TOKEN_LIFETIME     time.Duration = time.Duration(3600) * time.Second
	TEST_TOKEN_LIFETIME_STR string        = "1363466585"
	TEST_TOKEN_MAX_TIMEDIFF int64         = 1
)

// t.Fail() 		- fails test but continues execution
// t.FailNow()  	- fails test case immedieately and continues with next testcase
// t.Error("",...) 	- = Log() + FailNow()
// t.Fatal("",...) 	- = Log() + FailNow()
// t.Log("",...)    - log formatting like Println()

// ----------------------------------------------------------------------------
// Helper Functions just used during tests
//
// check whether token fields contain the expected values
func CheckTokenFields(t *testing.T, tok *AuthToken, uri, lifetime, sig string) {

	if tok == nil {
		t.Fatal("token = nil")
	}
	if tok.Resource != uri {
		t.Error("resource emty")
	}
	if tok.Valid_until != lifetime {
		t.Error("valid_until empty")
	}
	if sig != "" {
		if sig != tok.Signature {
			t.Error("signature empty")
		}
	}
}

func CheckTokenFieldsWithTime(t *testing.T, tok *AuthToken, uri string, lifetime, issuetime int64, sig string) {

	if tok == nil {
		t.Fatal("token = nil")
	}

	// test token resource field
	if tok.Resource != uri {
		t.Error("resource empty")
	}

	// test correct value of time field
	check_time, err := strconv.ParseInt(tok.Valid_until, 10, 64)
	if err != nil || check_time < 0 {
		t.Error("lifetime malformed")
	}
	if check_time-lifetime-issuetime > TEST_TOKEN_MAX_TIMEDIFF {
		t.Error("lifetime invalid")
	}

	// empty signature always fails
	if tok.Signature == "" {
		t.Error("signature empty")
	}

	// check signature if sig is not ""
	if sig != "" {
		if sig != tok.Signature {
			t.Error("signature mismatch")
		}
	}
}

func FakeAuthToken(uri, valid, sig string) *AuthToken {
	tok := AuthToken{
		Resource:    uri,
		Valid_until: valid,
		Signature:   sig,
	}
	return &tok
}

// ----------------------------------------------------------------------------
// Unit Tests

// test signing function independent of others
func TestAuthTokenSigning(t *testing.T) {

	tok := FakeAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "")
	if err := tok.Sign(TEST_TOKEN_SECRET); err != nil {
		t.Fatal("Sign() failed")
	}

	// verify signature manually
	var sig_base string
	sig_base = "rter_resource=" + TEST_TOKEN_URI + "&rter_valid_until=" + TEST_TOKEN_LIFETIME_STR
	hmac := hmac.New(sha256.New, bytes.NewBufferString(TEST_TOKEN_SECRET).Bytes())
	sig := url.QueryEscape(base64.StdEncoding.EncodeToString(hmac.Sum(bytes.NewBufferString(url.QueryEscape(sig_base)).Bytes())))

	if sig != tok.Signature {
		t.Error("signature mismatch")
	}
}

// test if token generation works
func TestAuthTokenGeneration(t *testing.T) {

	// save current time in UNIX seconds for later reference
	tm := time.Now().UTC().Unix()

	// create signed token
	tok, err := GenerateAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME, TEST_TOKEN_SECRET)

	if err != nil {
		t.Fatal("token generation failed")
	}

	// check fields
	CheckTokenFieldsWithTime(t, tok, TEST_TOKEN_URI,
		int64(TEST_TOKEN_LIFETIME.Seconds()), tm, "")

	// check token signature
	if tok.VerifySignature(TEST_TOKEN_SECRET) != nil {
		t.Error("signature verification failed")
	}

}

func TestAuthTokenGenerationFailing(t *testing.T) {

	// test token generation with wrong input parameters
	// - empty resource
	_, err := GenerateAuthToken("", TEST_TOKEN_LIFETIME, TEST_TOKEN_SECRET)
	if err != nil {
		t.Error("token created when url was empty")
	}

	// - malformed resource (no url)
	_, err = GenerateAuthToken("http:/nourl", TEST_TOKEN_LIFETIME, TEST_TOKEN_SECRET)
	if err != nil {
		t.Error("token created when url was malformed")
	}

	// - empty signing key
	_, err = GenerateAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME, "")
	if err != nil {
		t.Error("token created when signing key was empty")
	}

	// - negative or zero validity time
	_, err = GenerateAuthToken(TEST_TOKEN_URI, time.Duration(-1), TEST_TOKEN_SECRET)
	if err != nil {
		t.Error("token created when lifetime was negative")
	}

	// - negative or zero validity time
	_, err = GenerateAuthToken(TEST_TOKEN_URI, time.Duration(0), TEST_TOKEN_SECRET)
	if err != nil {
		t.Error("token created when lifetime was zero")
	}
}

func TestAuthTokenOutput(t *testing.T) {

	tok := FakeAuthToken("aaa", "bbb", "ccc")

	// string
	if tok.String() != `rtER rter_resource="aaa", rter_signature="ccc", rter_valid_until="bbb"` {
		t.Error("bad string output: '", tok.String(), "'")
	}

	// JSON
	var tok2 AuthToken
	err := json.Unmarshal(bytes.NewBufferString(tok.Json()).Bytes(), &tok2)

	if err != nil ||
		tok2.Resource != tok.Resource ||
		tok2.Valid_until != tok.Valid_until ||
		tok2.Signature != tok.Signature {
		t.Error("bad JSON output '", tok.Json(), "'")
	}
}

func TestAuthTokenInput(t *testing.T) {

	// from http.Request.Header with a single valid rtER Auth header
	r, err := http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "sig").String())
	tok, _ := NewFromHttpRequest(r)
	CheckTokenFields(t, tok, TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "sig")

	// from http.Request.Header with multiple Auth headers and a valid rtER header
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.SetBasicAuth("user", "pass")
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "sig").String())
	tok, _ = NewFromHttpRequest(r)
	CheckTokenFields(t, tok, TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "sig")

	// from http.Request.Header with an invalid rtER Auth header (empty sig)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME_STR, "").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with empty signature")
	}

	// from http.Request.Header with an invalid rtER Auth header (empty uri)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken("", TEST_TOKEN_LIFETIME_STR, "sig").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with empty uri")
	}

	// from http.Request.Header with an invalid rtER Auth header (empty lifetime)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, "", "sig").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with empty lifetime")
	}

	// from http.Request.Header with a malformed rtER Auth header (invalid URI syntax)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken("http;//example.com", TEST_TOKEN_LIFETIME_STR, "sig").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with invalid uri")
	}

	// from http.Request.Header with a malformed rtER Auth header (invalid lifetime)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, "10a", "sig").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with invalid lifetime")
	}

	// from http.Request.Header with a malformed rtER Auth header (negative lifetime)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", FakeAuthToken(TEST_TOKEN_URI, "-10", "sig").String())
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with negative lifetime")
	}

	// from http.Request.Header with a malformed rtER Auth header (invalid escaping: extra '"')
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", `rtER rter_resource="http://example.com",rter_valid_until="1234",rter_signature="si"g"`)
	tok, _ = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with invalid escaping, (extra ')")
	}

	// from http.Request.Header with a malformed rtER Auth header (invalid whitespace)
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", `rtER rter_resource ="http://example.com",rter_valid_until="1234",rter_signature="sig"`)
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with invalid whitespace")
	}

	// from http.Request.Header with a malformed rtER Auth header (missing '"')
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", `rtER rter_resource=http://example.com,rter_valid_until=1234,rter_signature=sig`)
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with missing \"")
	}

	// from http.Request.Header with a malformed rtER Auth header (invalid escaping: extra '=')
	r, _ = http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", `rtER rter_resource="http://example.com",rter_valid_until="1234",rter_signature="s=ig"`)
	tok, err = NewFromHttpRequest(r)
	if err == nil {
		t.Error("accepted auth header with with invalid escaping (extra '='")
	}

}

func TestAuthTokenSignatureVerification(t *testing.T) {

	// generate a valid token signature, ut do not use AuthToken.Sign()
	var sig_base string
	sig_base = "rter_resource=" + TEST_TOKEN_URI + "&rter_valid_until=3600"
	hmac := hmac.New(sha256.New, bytes.NewBufferString(TEST_TOKEN_SECRET).Bytes())
	sig := url.QueryEscape(base64.StdEncoding.EncodeToString(hmac.Sum(bytes.NewBufferString(url.QueryEscape(sig_base)).Bytes())))

	// fake and test a valid token for testing
	tok := FakeAuthToken(TEST_TOKEN_URI, "3600", sig)
	if tok.VerifySignature(TEST_TOKEN_SECRET) != nil {
		t.Error("valid signature rejected")
	}

	// fake and test a token with invalid URI
	tok = FakeAuthToken("http://invalid.uri", "3600", sig)
	if tok.VerifySignature(TEST_TOKEN_SECRET) == nil {
		t.Error("invalid uri validated")
	}

	// fake and test a token with invalid livetime
	tok = FakeAuthToken(TEST_TOKEN_URI, "1", sig)
	if tok.VerifySignature(TEST_TOKEN_SECRET) == nil {
		t.Error("invalid lifetime validated")
	}

	// fake and test a token with invalid signature
	tok = FakeAuthToken(TEST_TOKEN_URI, "3600", sig+"x")
	if tok.VerifySignature(TEST_TOKEN_SECRET) == nil {
		t.Error("invalid signature validated")
	}

	// fake and test a token with empty signature
	tok = FakeAuthToken(TEST_TOKEN_URI, "3600", "")
	if tok.VerifySignature(TEST_TOKEN_SECRET) == nil {
		t.Error("empty signature validated")
	}
}

func TestAuthTokenLifetimeVerification(t *testing.T) {

	// create fake token with valid lifetime
	future := time.Now().UTC().Unix() + 10
	tok_valid := FakeAuthToken("", strconv.FormatInt(future, 10), "")
	if tok_valid.VerifyLifetime() != nil {
		t.Error("valid lifetime rejected")
	}

	// create fake token with invalid lifetime
	past := future - 20
	tok_invalid := FakeAuthToken("", strconv.FormatInt(past, 10), "")
	if tok_invalid.VerifyLifetime() == nil {
		t.Error("invalid lifetime accepted")
	}
}

// test end-to-end producer - consumer interaction
func TestAuthTokenVerification(t *testing.T) {

	// correct
	tok, err := GenerateAuthToken(TEST_TOKEN_URI, TEST_TOKEN_LIFETIME, TEST_TOKEN_SECRET)
	if err != nil {
		t.Fatal("token generation failed")
	}
	r, _ := http.NewRequest("POST", TEST_TOKEN_URI, nil)
	r.Header.Add("Authorization", tok.String())
	if err := AuthenticateRequest(r, TEST_TOKEN_SECRET); err != nil {
		t.Error(fmt.Sprintf("auth failed: %s (%s)", err.Error(), tok.String()))
	}

	// wrong request url
	r, _ = http.NewRequest("POST", "http://example.org", nil)
	r.Header.Add("Authorization", tok.String())
	if err := AuthenticateRequest(r, TEST_TOKEN_SECRET); err == nil {
		t.Error(fmt.Sprintf("invalid request url for auth token: %s", err.Error()))
	}
}
