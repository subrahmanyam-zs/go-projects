package cspauth

import "testing"

func Test_Validate(t *testing.T) {
	tcs := []struct {
		opts *Options
		err  error
	}{
		{&Options{"",  "192.168.0.1",   "",  "CSP_SHARED_KEY",  "cd1",},ErrEmptyAppKey},
		{&Options{"Ubuntu",  "",   "ak11127983471298348912734",  "",  "cd1",},ErrEmptySharedKey},
		{&Options{"Ubuntu",  "",   "ak11127983471298348912734",  "CSP_SHARED_KEY",  "",},ErrEmptyAppID},
		{&Options{"",  "192.168.0.1",   "ak11127983471298348912734",  "CSP_SHARED_KEY",  "cd1",},nil},
	}

	for i,tc := range tcs {
		err := tc.opts.validate()
		if tc.err != err {
			t.Errorf("TESTCASE[%v]\nExpected:\n%v\nGot:\n%v",i,tc.err,err)
		}
	}
}
