package logfmt_test

import (
	"testing"

	"github.com/brotich/go-log-forwarder/internal/logfmt"
	qt "github.com/frankban/quicktest"
)

func TestToLogfmt(t *testing.T) {
	c := qt.New(t)
	c.Parallel()

	tests := []struct {
		about string

		input map[string]any

		expectedError error
		expectOutput  string
	}{{
		about: "success: convert json to logfmt",
		input: map[string]any{
			"time":                "2025-05-07T10:00:24.41618973Z",
			"level":               "INFO",
			"msg":                 "Request",
			"path":                "/up",
			"status":              200,
			"dur":                 8,
			"method":              "GET",
			"req_content_length":  0,
			"req_content_type":    "",
			"resp_content_length": 73,
			"resp_content_type":   "text/html; charset=utf-8", "remote_addr": "300.180.67.252, 10.10.10.4",
			"user_agent": "Blackbox Exporter/0.26.0",
			"cache":      "miss", "query": "",
		},
		expectOutput: `cache=miss dur=8 level=INFO method=GET msg=Request path=/up query= remote_addr="300.180.67.252, 10.10.10.4" req_content_length=0 req_content_type= resp_content_length=73 resp_content_type="text/html; charset=utf-8" status=200 time=2025-05-07T10:00:24.41618973Z user_agent="Blackbox Exporter/0.26.0"`,
	}}

	for _, test := range tests {
		c.Run(test.about, func(c *qt.C) {
			c.Parallel()

			t := logfmt.JSONFmt{}
			out, err := t.ToLogfmt(test.input)
			if test.expectedError != nil {
				c.Assert(err, qt.ErrorMatches, test.expectedError)
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(string(out), qt.Equals, test.expectOutput)
			}
		})
	}
}
