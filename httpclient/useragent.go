package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform/version"
)

const userAgentFormat = "Terraform/%s"
const uaEnvVar = "TF_APPEND_USER_AGENT"

func UserAgentString() string {
	ua := fmt.Sprintf(userAgentFormat, version.Version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", rt.userAgent)
	}
	return rt.inner.RoundTrip(req)
}

func TerraformUserAgent(version string) *userAgent {
	return &userAgent{[]*UAProduct{
		{"HashiCorp", "1.0", ""},
		{"Terraform", version, ""},
	}}
}

type UAProduct struct {
	Product string
	Version string
	Comment string
}

func (uap *UAProduct) String() string {
	var b strings.Builder
	b.WriteString(uap.Product)
	if uap.Version != "" {
		b.WriteString(fmt.Sprintf("/%s", uap.Version))
	}
	if uap.Comment != "" {
		b.WriteString(fmt.Sprintf(" (%s)", uap.Comment))
	}
	return b.String()
}

func (uap *UAProduct) Equal(p *UAProduct) bool {
	if uap.Product == p.Product && uap.Version == p.Version && uap.Comment == p.Comment {
		return true
	}
	return false
}

type userAgent struct {
	products []*UAProduct
}

func (ua *userAgent) Products() []*UAProduct {
	return ua.products
}

func (ua *userAgent) String() string {
	var b strings.Builder
	for i, p := range ua.products {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(p.String())
	}

	return b.String()
}

func (ua *userAgent) Append(uap []*UAProduct) *userAgent {
	ua.products = append(ua.products, uap...)
	return ua
}

func (ua *userAgent) Equal(userAgent *userAgent) bool {
	if len(ua.products) != len(userAgent.products) {
		return false
	}

	for i, p := range ua.products {
		if p != userAgent.products[i] {
			return false
		}
	}

	return true
}

func UserAgent(products []*UAProduct) *userAgent {
	ua := &userAgent{products}

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua.Append(ParseUserAgentString(add))
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

func ParseUserAgentString(uaString string) []*UAProduct {
	// TODO
	return []*UAProduct{}
}
