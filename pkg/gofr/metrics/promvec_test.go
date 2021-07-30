package metrics

import (
	"reflect"
	"strings"
	"testing"
)

func Test_IncCounter(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		err  error
	}{
		{"success-case", "test", nil},
		{"error-case", "test1", metricNotFound},
	}

	p := newPromVec()

	err := p.registerCounter("test", "testing method", "code", "method")
	if err != nil {
		t.Errorf("error while creating counter")
	}

	for i, tc := range tcs {
		err = p.IncCounter(tc.name, "200", "POST")
		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("TESTCASE[%v] expected error %T, got %T", i, tc.err, err)
		}
	}
}

func Test_AddCounter(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		err  error
	}{
		{"success-case", "test_counter", nil},
		{"error-case", "test1", metricNotFound},
	}

	p := newPromVec()

	err := p.registerCounter("test_counter", "testing method", "code", "method")
	if err != nil {
		t.Errorf("error while creating counter")
	}

	for i, tc := range tcs {
		err = p.AddCounter(tc.name, float64(i), "200", "POST")
		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("TESTCASE[%v] expected error %v, got %v", i, tc.err, err)
		}
	}
}

func Test_ObserveHistogram(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		err  error
	}{
		{"success-case", "test_histogram", nil},
		{"error-case", "test", metricNotFound},
	}

	p := newPromVec()

	err := p.registerHistogram("test_histogram", "testing method",
		[]float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30},
		"code", "method")
	if err != nil {
		t.Errorf("error while creating histogram")
	}

	for i, tc := range tcs {
		err = p.ObserveHistogram(tc.name, float64(i), "200", "POST")
		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("TESTCASE[%v] expected error %v, got %v", i, tc.err, err)
		}
	}
}

func Test_SetGauge(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		err  error
	}{
		{"success-case", "test_gauge", nil},
		{"error-case", "test", metricNotFound},
	}

	p := newPromVec()

	err := p.registerGauge("test_gauge", "set value of gauge", "no_of_go_routines")
	if err != nil {
		t.Errorf("error while creating gauge")
	}

	for i, tc := range tcs {
		err = p.SetGauge(tc.name, float64(i), "no_of_go_routines")
		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("TESTCASE[%v] expected error %v, got %v", i, tc.err, err)
		}
	}
}

func Test_ObserveSummary(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		err  error
	}{
		{"success-case", "test_summary", nil},
		{"error-case", "test", metricNotFound},
	}

	p := newPromVec()

	err := p.registerSummary("test_summary", "testing method", "code", "method")
	if err != nil {
		t.Errorf("error while creating summary")
	}

	for i, tc := range tcs {
		err = p.ObserveSummary(tc.name, float64(i), "200", "POST")
		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("TESTCASE[%v] expected error %v, got %v", i, tc.err, err)
		}
	}
}

func Test_InvalidLabelError(t *testing.T) {
	errString := "inconsistent label cardinality"

	p := newPromVec()

	err := p.registerCounter("label_counter", "testing method", "code", "method")
	if err != nil {
		t.Errorf("error while creating counter")
	}

	err = p.registerHistogram("label_histogram", "testing method",
		[]float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30},
		"code", "method")
	if err != nil {
		t.Errorf("error while creating histogram")
	}

	err = p.registerGauge("label_gauge", "set value of gauge", "no_of_go_routines")
	if err != nil {
		t.Errorf("error while creating gauge")
	}

	err = p.registerSummary("label_summary", "testing method", "code", "method")
	if err != nil {
		t.Errorf("error while creating summary")
	}

	err = p.IncCounter("label_counter")
	if err == nil || !strings.Contains(err.Error(), errString) {
		t.Errorf("expected err string %v, got %v", errString, err)
	}

	err = p.AddCounter("label_counter", 2)
	if err == nil || !strings.Contains(err.Error(), errString) {
		t.Errorf("expected err string %v, got %v", errString, err)
	}

	err = p.SetGauge("label_gauge", 2)
	if err == nil || !strings.Contains(err.Error(), errString) {
		t.Errorf("expected err string %v, got %v", errString, err)
	}

	err = p.ObserveSummary("label_summary", 2)
	if err == nil || !strings.Contains(err.Error(), errString) {
		t.Errorf("expected err string %v, got %v", errString, err)
	}

	err = p.ObserveHistogram("label_histogram", 2)
	if err == nil || !strings.Contains(err.Error(), errString) {
		t.Errorf("expected err string %v, got %v", errString, err)
	}
}
