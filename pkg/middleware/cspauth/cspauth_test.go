package cspauth

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockHandler struct{}

func (r *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Response == nil {
		w.WriteHeader(http.StatusOK)

		return
	}

	w.WriteHeader(req.Response.StatusCode)
	resp, err := io.ReadAll(req.Response.Body)

	if err == nil {
		_, _ = w.Write(resp)
	}
}

func TestCSPAuth(t *testing.T) {

	tcs := []struct {
		appKey      string
		clientID    string
		authContext string
		sharedKey   string
		expCode     int
		body        string
	}{
		{
			"ankling123jerkins4junked",
			"cd1",
			"RUtKQ3UzYWRWTHhYQ1FrdkY0alprQzVvZCtySWE5c3FDQnVJTHdnQXlZZmoxWHcrMk1mQVpYWXVBQ21tNkhuTlpXMzZHb1h3akxhSUVVbU1CdW83NVhOaWZqUHZCSlQxWldiSG1Cem1OQ05ZYUVBY2VFeGtNdisvczlBbHlBSzdDMVlXVWJSTzRselQ4VjBZcWI5ZE85L3dvNjJrWkk4Y1VqUDVtUGhVT1U1NzJDYzl1S0JNYkpaSyt1S2l0SmlVcnVURTlHSUZSTjlibGFvT2IrRVRiOTlDWGtnRU9iR2ZURmhNeWd0ZnhJZVB6N0xxeFFXd1YvSmZsb1RWSm1JVHg1SW5HRHBiN1RZR3B5VXU3VmdpT3hGSnlmMGJxTXVlQkFxTHZ5Unpic0dHVlltWXVCNnM2NW0zT1l2NFBTdkhLNHgxZDFjcmxMYjE4cHk4VjNpNm5ia1dGTkwzRTJSOFh4MS9HSVVjZ2x5NVgvQlJUMHh5KzAwT1FDbGdQcDk3ODEzZTU1",
			"CSP_APP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"",
			"CSP_APP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak",
			"cd1",
			"RUtKQ3UzYWRWTHhYQ1FrdkY0alprQzVvZCtySWE5c3FDQnVJTHdnQXlZZmoxWHcrMk1mQVpYWXVBQ21tNkhuTlpXMzZHb1h3akxhSUVVbU1CdW83NVhOaWZqUHZCSlQxWldiSG1Cem1OQ05ZYUVBY2VFeGtNdisvczlBbHlBSzdDMVlXVWJSTzRselQ4VjBZcWI5ZE85L3dvNjJrWkk4Y1VqUDVtUGhVT1U1NzJDYzl1S0JNYkpaSyt1S2l0SmlVcnVURTlHSUZSTjlibGFvT2IrRVRiOTlDWGtnRU9iR2ZURmhNeWd0ZnhJZVB6N0xxeFFXd1YvSmZsb1RWSm1JVHg1SW5HRHBiN1RZR3B5VXU3VmdpT3hGSnlmMGJxTXVlQkFxTHZ5Unpic0dHVlltWXVCNnM2NW0zT1l2NFBTdkhLNHgxZDFjcmxMYjE4cHk4VjNpNm5ia1dGTkwzRTJSOFh4MS9HSVVjZ2x5NVgvQlJUMHh5KzAwT1FDbGdQcDk3ODEzZTU1",
			"CSP_APP_SHARED_KEY",
			400,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"",
			"RUtKQ3UzYWRWTHhYQ1FrdkY0alprQzVvZCtySWE5c3FDQnVJTHdnQXlZZmoxWHcrMk1mQVpYWXVBQ21tNkhuTlpXMzZHb1h3akxhSUVVbU1CdW83NVhOaWZqUHZCSlQxWldiSG1Cem1OQ05ZYUVBY2VFeGtNdisvczlBbHlBSzdDMVlXVWJSTzRselQ4VjBZcWI5ZE85L3dvNjJrWkk4Y1VqUDVtUGhVT1U1NzJDYzl1S0JNYkpaSyt1S2l0SmlVcnVURTlHSUZSTjlibGFvT2IrRVRiOTlDWGtnRU9iR2ZURmhNeWd0ZnhJZVB6N0xxeFFXd1YvSmZsb1RWSm1JVHg1SW5HRHBiN1RZR3B5VXU3VmdpT3hGSnlmMGJxTXVlQkFxTHZ5Unpic0dHVlltWXVCNnM2NW0zT1l2NFBTdkhLNHgxZDFjcmxMYjE4cHk4VjNpNm5ia1dGTkwzRTJSOFh4MS9HSVVjZ2x5NVgvQlJUMHh5KzAwT1FDbGdQcDk3ODEzZTU1",
			"CSP_APP_SHARED_KEY",
			400,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"RUtKQ3UzYWRWTHhYQ1FrdkY0alprQzVvZCtySWE5c3FDQnVJTHdnQXlZZmoxWHcrMk1mQVpYWXVBQ21tNkhuTlpXMzZHb1h3akxhSUVVbU1CdW83NVhOaWZqUHZCSlQxWldiSG1Cem1OQ05ZYUVBY2VFeGtNdisvczlBbHlBSzdDMVlXVWJSTzRselQ4VjBZcWI5ZE85L3dvNjJrWkk4Y1VqUDVtUGhVT1U1NzJDYzl1S0JNYkpaSyt1S2l0SmlVcnVURTlHSUZSTjlibGFvT2IrRVRiOTlDWGtnRU9iR2ZURmhNeWd0ZnhJZVB6N0xxeFFXd1YvSmZsb1RWSm1JVHg1SW5HRHBiN1RZR3B5VXU3VmdpT3hGSnlmMGJxTXVlQkFxTHZ5Unpic0dHVlltWXVCNnM2NW0zT1l2NFBTdkhLNHgxZDFjcmxMYjE4cHk4VjNpNm5ia1dGTkwzRTJSOFh4MS9HSVVjZ2x5NVgvQlJUMHh5KzAwT1FDbGdQcDk3ODEzZTU1",
			"",
			400,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"c29tZSB",
			"CSP_APP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"RUtKQ3UzYWRWTHhYQ1FrdkY0alprQzVvZCtySWE5c3FDQnVJTHdnQXlZZmoxWHcrMk1mQVpYWXVBQ21tNkhuTlpXMzZHb1h3akxhSUVVbU1CdW83NVRWTnFmb2ZCdHVGZWE1NjhmOXB0VmpqZkFNeTJWdDFwMGcvZmpJaUw1MitBeFJuOVNISXQvNU1vWnRKSDFwbjYrbFV6QXRjeU5za0xLcFkzclJIQnQ5S2J5YlpwQlVaM3FtVVdDSE9jd1loMythUW1RK1RSc3dxOSsvYUtwYVdDZ1VIUGZVWWdYZEtjejVRdTRCeHA4MHJYM2UwV2FMeXczV3BtNndPMENhNFZhU2dNSXdvV3hzcDg1TUhVNWx3Ry94eXBPUlRmOTlGU0FHN2pUcU91ZzBoN1crYWxmeStXUlNuMy9vQ2ZrczVLNGxTSzF4OFpKNWZFMzFtMG9KeXo5Z2hSUytidWc3Q0xUR0lSMVRoYnA0ajhvYWFKVGFFd0taU0tMSDBwakdSOWE4NTFh",
			"CSP_APP_SHARED_KEY",
			200,
			"Dummy body",
		},
	}

	for i, tc := range tcs {
		os.Setenv("CSP_APP_SHARED_KEY", tc.sharedKey)
		body := bytes.NewReader([]byte(tc.body))
		if tc.body == "" {
			body = nil
		}
		req := httptest.NewRequest("GET", "/dummy", body)
		req.Header.Set("ak", tc.appKey)
		req.Header.Set("cd", tc.clientID)
		req.Header.Set("ac", tc.authContext)
		w := httptest.NewRecorder()

		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)

		handler := CSPAuth(logger)(&MockHandler{})
		handler.ServeHTTP(w, req)

		if w.Code != tc.expCode {
			t.Errorf("TESTCASE[%v]\nexpected code %v,\ngot %v", i, tc.expCode, w.Code)
		}
	}
}

func Test_getBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "/dummy", nil)
	req.Body = nil

	b := getBody(req)
	if len(b) != 0 {
		t.Errorf("expected empty slice, got %v", b)
	}
}

func Test_pkcs7Pad(t *testing.T) {
	tcs := []struct {
		blockSize int
		err       error
	}{
		{0, errInvalidBlockSize},
		{1, errInvalidPKCS7Data},
	}

	for i, tc := range tcs {
		_, err := pkcs7Pad(nil, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}

func Test_pkcs7Unpad(t *testing.T) {
	tcs := []struct {
		data      []byte
		blockSize int
		err       error
	}{
		{nil,0, errInvalidBlockSize},
		{nil,1, errInvalidPKCS7Data},
		{[]byte{1,2,3},2, errInvalidPKCS7Padding},
		{[]byte{1,2,3,5},2, errInvalidPKCS7Padding},
	}

	for i, tc := range tcs {
		_, err := pkcs7Unpad(tc.data, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}
