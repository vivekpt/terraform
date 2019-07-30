package httpclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/version"
)

func TestUserAgentString_env(t *testing.T) {
	expectedBase := fmt.Sprintf(userAgentFormat, version.Version)
	if oldenv, isSet := os.LookupEnv(uaEnvVar); isSet {
		defer os.Setenv(uaEnvVar, oldenv)
	} else {
		defer os.Unsetenv(uaEnvVar)
	}

	for i, c := range []struct {
		expected   string
		additional string
	}{
		{expectedBase, ""},
		{expectedBase, " "},
		{expectedBase, " \n"},

		{fmt.Sprintf("%s test/1", expectedBase), "test/1"},
		{fmt.Sprintf("%s test/2", expectedBase), "test/2 "},
		{fmt.Sprintf("%s test/3", expectedBase), " test/3 "},
		{fmt.Sprintf("%s test/4", expectedBase), "test/4 \n"},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if c.additional == "" {
				os.Unsetenv(uaEnvVar)
			} else {
				os.Setenv(uaEnvVar, c.additional)
			}

			actual := UserAgentString()

			if c.expected != actual {
				t.Fatalf("Expected User-Agent '%s' does not match '%s'", c.expected, actual)
			}
		})
	}
}

func TestParseUserAgentString(t *testing.T) {
	testCases := []struct {
		uaString   string
		uaProducts []*UAProduct
	}{
		{
			"terraform-github-actions/1.0",
			[]*UAProduct{{"terraform-github-actions", "1.0", ""}},
		},
		{
			"TFE/a718e58f",
			[]*UAProduct{{"TFE", "a718e58f", ""}},
		},
		{
			"OneProduct/0.1.0 AnotherOne/1.2",
			[]*UAProduct{{"OneProduct", "0.1.0", ""}, {"AnotherOne", "1.2", ""}},
		},
		{
			"ProductWithComment/1.0.0 (a comment; goes; here)",
			[]*UAProduct{{"ProductWithComment", "1.0.0", "a comment; goes; here"}},
		},
		{
			"ProductWithComment/1.0.0 (a comment; goes; here) AnotherProductWithComment/5.5.0 (blah)",
			[]*UAProduct{
				{"ProductWithComment", "1.0.0", "a comment; goes; here"},
				{"AnotherProductWithComment", "5.5.0", "blah"},
			},
		},
		{
			"NoComment/1.0.0 AnotherProductWithComment/5.5.0 (blah)",
			[]*UAProduct{
				{"NoComment", "1.0.0", ""},
				{"AnotherProductWithComment", "5.5.0", "blah"},
			},
		},
		{
			"First/1.0.0 Second/5.5.0 Third/5.5.0",
			[]*UAProduct{
				{"First", "1.0.0", ""},
				{"Second", "5.5.0", ""},
				{"Third", "5.5.0", ""},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			givenUA := UserAgent(ParseUserAgentString(tc.uaString))
			expectedUA := UserAgent(tc.uaProducts)

			if !givenUA.Equal(expectedUA) {
				t.Fatalf("Unexpected User-Agent.\nExpected: %q\nGiven: %q\n", expectedUA, givenUA)
			}
		})
	}

}

func TestUserAgentAppendViaEnvVar(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestUserAgentString(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestUserAgentAppend(t *testing.T) {
	t.Fatal("Not implemented")
}
