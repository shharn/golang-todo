package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SanitizePageNoTestCase struct {
	input []string
	expected int
}

func TestSanitizePageNo(t *testing.T) {
	tcs := []SanitizePageNoTestCase{
		0: SanitizePageNoTestCase{
			input: []string{"1"},
			expected: 1,
		},
		1: SanitizePageNoTestCase{
			input: []string{"-1"},
			expected: 1,
		},
		2: SanitizePageNoTestCase{
			input: []string{"0"},
			expected: 1,
		},
		3: SanitizePageNoTestCase{
			input: []string{ "5", "1" },
			expected: 5,
		},
		4: SanitizePageNoTestCase{
			input: []string{},
			expected: 1,
		},
		5: SanitizePageNoTestCase{
			input: []string{"asdf"},
			expected: 1,
		},
		6: SanitizePageNoTestCase{
			input: []string{"10", "asdf"},
			expected: 10,
		},
	}
	for _, tc := range tcs {
		actual := sanitizePageNo(tc.input)
		assert.Equal(t, tc.expected, actual)
	}
}

type HandlerTestCase struct {
	method string
	path string
	rq *http.Request
	recorder *httptest.ResponseRecorder
	resp *http.Response
	handler http.HandlerFunc
	expectedStatusCode int
}

func (tc *HandlerTestCase) Execute() {
	tc.rq = httptest.NewRequest(tc.method, fmt.Sprintf("http://localhost/%s", tc.path), nil)
	tc.recorder = httptest.NewRecorder()
	tc.handler(tc.recorder, tc.rq)
	tc.resp = tc.recorder.Result()
}

func (tc *HandlerTestCase) GetStatusCode() int {
	return tc.resp.StatusCode
}

type MockFinder struct {}

func (mf MockFinder) GetAsByteArray() ([]byte, error) {
	return []byte{}, nil
}

func TestIndexHandler(t *testing.T) {
	mockFinder := MockFinder{}
	tcs := []HandlerTestCase{
		0: HandlerTestCase{
			method: "GET",
			path: "",
			handler: WithFinder(IndexHandler, mockFinder),
			expectedStatusCode: http.StatusOK,
		},
		1: HandlerTestCase{
			method: "GET",
			path: "notfound",
			handler: WithFinder(IndexHandler, mockFinder),
			expectedStatusCode: http.StatusNotFound,
		},
	}
	
	for _, tc := range tcs {
		tc.Execute()
		assert.Equal(t, tc.expectedStatusCode, tc.GetStatusCode())
	}
}