package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// Parse cmdline arguments using flag package
	oidcUrl := flag.String("oidc-url", "", "OIDC IdP's URL to get thumbprint for")
	port := flag.Uint("port", 443, "Port that has TLS")
	flag.Parse()
	if *oidcUrl == "" {
		flag.Usage()
		os.Exit(1)
	}

	config := getOpenIdConfiguration(*oidcUrl)
	server := config.getServerFromOpenIdConfiguration()

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", *server, *port), &tls.Config{})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	// Get the ConnectionState struct as that's the one which gives us x509.Certificate struct
	certificates := conn.ConnectionState().PeerCertificates
	cert := certificates[len(certificates)-1]
	fingerprint := sha1.Sum(cert.Raw)
	var buf bytes.Buffer
	for _, f := range fingerprint {
		fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Println(buf.String())

	conn.Close()
}

type OpenIdConfiguration struct {
	JwksUri string `json:"jwks_uri"`
}

func getOpenIdConfiguration(uri string) OpenIdConfiguration {
	uri += ".well-known/openid-configuration"
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var openidConfig OpenIdConfiguration
	err = json.Unmarshal(body, &openidConfig)
	if err != nil {
		panic(err)
	}

	return openidConfig
}

func (config *OpenIdConfiguration) getServerFromOpenIdConfiguration() *string {
	u, err := url.Parse(config.JwksUri)
	if err != nil {
		panic(err)
	}

	var hostname string
	hostname = u.Hostname()

	return &hostname
}
