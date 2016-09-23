package darksky

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForecastRequest_Get(t *testing.T) {
	usingTestServer(validForecastHandler, func(testURL string) {
		resp := MakeRequest(key, 41.8781, -87.6297).WithBaseURL(testURL).Get()

		if resp.Error != nil {
			t.Error(resp.Error)
			return
		}

		if resp.APICallCount != 1 {
			t.Errorf("Expected APICallCount to be %v but was %v.", 1, resp.APICallCount)
		}

		forecast := resp.Forecast

		if len(forecast.Alerts) != 3 {
			t.Error("Expected 3 Alerts.")
		}
	})
}

func TestErrorResponse(t *testing.T) {
	usingTestServer(errorForecastHandler, func(testURL string) {
		resp := MakeRequest(key, 41.8781, -87.6297).WithBaseURL(testURL).Get()

		if resp.Error == nil {
			t.Error("Expected an HTTP Error Response to result in an error.")
		}

		if resp.Error.Error() != "A Server Error Occurred." {
			t.Error("Error() was not the expected value.")
		}
	})
}

func TestWindDirection(t *testing.T) {
	dp := DataPoint{WindBearing: 147}

	if dp.WindDirection() != "SE" {
		t.Errorf("Expected WindBearing of %v to be SE, was %v.", dp.WindBearing, dp.WindDirection())
	}

	dp.WindBearing = 0

	if dp.WindDirection() != "N" {
		t.Errorf("Expected WindBearing of %v to be N, was %v.", dp.WindBearing, dp.WindDirection())
	}

	dp.WindBearing = 336

	if dp.WindDirection() != "NW" {
		t.Errorf("Expected WindBearing of %v to be NW, was %v.", dp.WindBearing, dp.WindDirection())
	}
}

var key string = "test_key"

var validForecastHandler http.HandlerFunc = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
	jsonBytes, _ := ioutil.ReadFile("testdata/chicago_forecast.json")
	resp.Header().Add(APICallsHeader, "1")
	resp.WriteHeader(200)
	resp.Write(jsonBytes)
})

var errorForecastHandler http.HandlerFunc = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(500)
	resp.Write([]byte("A Server Error Occurred."))
})

func usingTestServer(handler http.HandlerFunc, runTest func(testURL string)) {
	ts := httptest.NewServer(handler)

	defer func() {
		ts.Close()
	}()

	runTest(ts.URL)
}
