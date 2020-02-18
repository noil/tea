package tea

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"
)

func timestamp() string {
	return fmt.Sprintf("%v", time.Now().UTC().Unix())
}

func nonce() string {
	nonce := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(nonce)
}

func (tea *TwitterEngagementAPI) oauthHeader(signingURL string) string {
	oauthParams := map[string]string{
		"oauth_consumer_key":     tea.token.ConsumerKey,
		"oauth_nonce":            nonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        timestamp(),
		"oauth_token":            tea.token.Access,
		"oauth_version":          "1.0",
	}
	signatureParts := []string{
		"POST",
		url.QueryEscape(signingURL),
		url.QueryEscape(sortedQueryString(oauthParams)),
	}
	signatureBase := strings.Join(signatureParts, "&")
	signingKey := tea.token.ConsumerKeySecret + "&" + tea.token.AccessSecret
	signer := hmac.New(sha1.New, []byte(signingKey))
	signer.Write([]byte(signatureBase))
	oauthParams["oauth_signature"] = base64.StdEncoding.EncodeToString(signer.Sum(nil))

	oauthParamKeys := make([]string, 0)
	for k := range oauthParams {
		oauthParamKeys = append(oauthParamKeys, k)
	}
	sort.Strings(oauthParamKeys)
	buf := new(bytes.Buffer)
	buf.WriteString("OAuth")
	for _, v := range oauthParamKeys {
		buf.WriteByte(' ')
		buf.WriteString(rfc3986Escape(v))
		buf.WriteString("=\"")
		buf.WriteString(rfc3986Escape(oauthParams[v]))
		buf.WriteString("\",")
	}
	oauth := buf.String()

	return oauth[:len(oauth)-1]
}

func sortedQueryString(values map[string]string) string {
	if len(values) == 0 {
		return ""
	}
	pairs := make(sortedPairs, 0)
	for k, v := range values {
		pairs = append(pairs, pair{rfc3986Escape(k), rfc3986Escape(v)})
	}
	sort.Sort(pairs)
	buf := new(bytes.Buffer)
	buf.WriteString(pairs[0].k)
	buf.WriteByte('=')
	buf.WriteString(pairs[0].v)

	for _, p := range pairs[1:] {
		buf.WriteByte('&')
		buf.WriteString(p.k)
		buf.WriteByte('=')
		buf.WriteString(p.v)
	}
	return buf.String()
}

type pair struct{ k, v string }
type sortedPairs []pair

func (sp sortedPairs) Len() int {
	return len(sp)
}

func (sp sortedPairs) Swap(i, j int) {
	sp[i], sp[j] = sp[j], sp[i]
}

func (sp sortedPairs) Less(i, j int) bool {
	if sp[i].k == sp[j].k {
		return sp[i].v < sp[j].v
	}
	return sp[i].k < sp[j].k
}

// Escapes a string more in line with Rfc3986 than http.URLEscape.
// URLEscape was converting spaces to "+" instead of "%20", which was messing up
// the signing of requests.
func rfc3986Escape(input string) string {
	firstEsc := -1
	b := []byte(input)
	for i, c := range b {
		if !isSafeChar(c) {
			firstEsc = i
			break
		}
	}

	// If nothing needed to be escaped, then the input is clean and
	// we're done.
	if firstEsc == -1 {
		return input
	}

	// If something did need to be escaped, write the prefix that was
	// fine to the buffer and iterate through the rest of the bytes.
	output := new(bytes.Buffer)
	output.Write(b[:firstEsc])

	for _, c := range b[firstEsc:] {
		if isSafeChar(c) {
			output.WriteByte(c)
		} else {
			fmt.Fprintf(output, "%%%02X", c)
		}
	}
	return output.String()
}

func isSafeChar(c byte) bool {
	return ('0' <= c && c <= '9') ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		c == '-' || c == '.' || c == '_' || c == '~'
}
