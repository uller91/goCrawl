package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorExpected bool		
	}{
		{
			name:     "remove scheme",
			inputURL: "https://crawler-test.com/path",
			expected: "crawler-test.com/path",
			errorExpected: false,
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://crawler-test.com/path/",
			expected: "crawler-test.com/path",
			errorExpected: false,
		},
		{
			name:     "lowercase capital letters",
			inputURL: "https://CRAWLER-TEST.com/PATH",
			expected: "crawler-test.com/path",
			errorExpected: false,
		},
		{
			name:     "remove scheme and capitals and trailing slash",
			inputURL: "http://CRAWLER-TEST.com/path/",
			expected: "crawler-test.com/path",
			errorExpected: false,
		},
		{
			name:          "handle invalid URL",
			inputURL:      `:\\invalidURL`,
			expected:      "",
			errorExpected: true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				if !tc.errorExpected {
					t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
					return
				}				
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}