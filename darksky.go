/*
Package darksky provides a Go API for accessing the DarkSky HTTP API.

For Dark Sky API documentation refer to:

	https://darksky.net/dev/docs

Requires an API Key to use. To register go to:

	https://darksky.net/dev/register
*/
package darksky

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"fmt"
)

// Forecast is the top level representation of the weather forecast for a location.
type Forecast struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timezone  string    `json:"timezone"`
	Offset    int       `json:"offset"`
	Currently DataPoint `json:"currently,omitempty"`
	Minutely  DataBlock `json:"minutely,omitempty"`
	Hourly    DataBlock `json:"hourly,omitempty"`
	Daily     DataBlock `json:"daily,omitempty"`
	Alerts    []Alert   `json:"alerts,omitempty"`
	Flags     Flags     `json:"flags,omitempty"`
}

// DataPoint is the current weather data for a single point in time.
type DataPoint struct {
	Time                   int64   `json:"time"`
	Summary                string  `json:"summary"`
	Icon                   string  `json:"icon"`
	SunriseTime            int64   `json:"sunriseTime"`
	SunsetTime             int64   `json:"sunsetTime"`
	PrecipIntensity        float64 `json:"precipIntensity"`
	PrecipIntensityMax     float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime int64   `json:"precipIntensityMaxTime"`
	PrecipProbability      float64 `json:"precipProbability"`
	PrecipType             string  `json:"precipType"`
	PrecipAccumulation     float64 `json:"precipAccumulation"`
	Temperature            float64 `json:"temperature"`
	TemperatureMin         float64 `json:"temperatureMin"`
	TemperatureMinTime     int64   `json:"temperatureMinTime"`
	TemperatureMax         float64 `json:"temperatureMax"`
	TemperatureMaxTime     int64   `json:"temperatureMaxTime"`
	ApparentTemperature    float64 `json:"apparentTemperature"`
	DewPoint               float64 `json:"dewPoint"`
	WindSpeed              float64 `json:"windSpeed"`
	WindBearing            float64 `json:"windBearing"`
	CloudCover             float64 `json:"cloudCover"`
	Humidity               float64 `json:"humidity"`
	Pressure               float64 `json:"pressure"`
	Visibility             float64 `json:"visibility"`
	Ozone                  float64 `json:"ozone"`
	MoonPhase              float64 `json:"moonPhase"`
}

// WindDirection converts the numerical WindBearing value in degrees to directional text. (ex: 200 => "SW")
func (dp DataPoint) WindDirection() string {
	direction := ""

	if dp.WindBearing > 293 || dp.WindBearing < 67 {
		direction += "N"
	}

	if dp.WindBearing < 247 && dp.WindBearing > 113 {
		direction += "S"
	}

	if dp.WindBearing > 22 && dp.WindBearing < 157 {
		direction += "E"
	}

	if dp.WindBearing < 337 && dp.WindBearing > 203 {
		direction += "W"
	}

	return direction
}

// DataBlock is a collection of data points over a period of time.
type DataBlock struct {
	Summary string      `json:"summary"`
	Icon    string      `json:"icon"`
	Data    []DataPoint `json:"data"`
}

// Alert is a potentially serious weather condition.
type Alert struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Expires     int64  `json:"expires"`
	URI         string `json:"uri"`
}

// Flags contains meta data about the Forecast.
type Flags struct {
	DarkSkyUnavailable string   `json:"darksky-unavailable"`
	DarkSkyStations    []string `json:"darksky-stations"`
	DataPointStations  []string `json:"datapoint-stations"`
	ISDStations        []string `json:"isds-stations"`
	LAMPStations       []string `json:"lamp-stations"`
	METARStations      []string `json:"metars-stations"`
	METNOLicense       string   `json:"metnol-license"`
	Sources            []string `json:"sources"`
	Units              string   `json:"units"`
}

// Options are used to modify the data representation of the Forecast.
type Options struct {
	Time         int64
	Exclude      []string
	ExtendHourly bool
	Lang         Lang
	Units        Units
}

type ForecastRequest struct {
	Key string
	Lat float64
	Lng float64
	Time int64
	Lang Lang
	Units Units
	ExtendHourly bool
	Exclude []string
	baseURL string
}

/*
ForecastResponse is a wrapper struct for a response from the DarkSky API.

Errors are included to make it easier to pass single values via channel from a goroutine.
*/
type ForecastResponse struct {
	Forecast     Forecast
	APICallCount int
	Error        error
}

func MakeRequest(key string, latitude float64, longitude float64) *ForecastRequest {
	return &ForecastRequest{
		Key: key,
		Lat: latitude,
		Lng: longitude,
		Time: -1,
		Lang: English,
		Units: US,
		ExtendHourly: false,
		Exclude: []string{},
		baseURL: "https://api.darksky.net/forecast",
	}
}

func (f *ForecastRequest) Get() ForecastResponse {
	forecastResponse := ForecastResponse{}

	requestURL := fmt.Sprintf("%v/%v/%v,%v", f.baseURL, f.Key, f.Lat, f.Lng)

	res, err := http.Get(requestURL)

	if err != nil {
		forecastResponse.Error = err
		return forecastResponse
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		forecastResponse.Error = err
		return forecastResponse
	}

	if res.StatusCode >= 400 {
		forecastResponse.Error = errors.New(string(body))
		return forecastResponse
	}

	callCount, err := strconv.Atoi(res.Header.Get(APICallsHeader))

	if err == nil {
		forecastResponse.APICallCount = callCount
	}

	forecast, err := fromJSON(body)

	if err != nil {
		forecastResponse.Error = err
		return forecastResponse
	}

	forecastResponse.Forecast = *forecast

	return forecastResponse
}

func (f *ForecastRequest) WithBaseURL(baseURL string) *ForecastRequest {
	f.baseURL = baseURL
	return f
}

func (f *ForecastRequest) WithTime(t int64) *ForecastRequest {
	f.Time = t
	return f
}

func (f *ForecastRequest) WithLang(l Lang) *ForecastRequest {
	f.Lang = l
	return f
}

func (f *ForecastRequest) WithUnit(u Units) *ForecastRequest {
	f.Units = u
	return f
}

// Units defines the possible options for measurement units used in the response.
type Units string

const (
	US   Units = "us"
	SI   Units = "si"
	CA   Units = "ca"
	UK   Units = "uk"
	UK2  Units = "uk2"
	AUTO Units = "auto"
)

// Lang defines the possible options for the text summary language.
type Lang string

const (
	Arabic             Lang = "ar"
	Bosnian            Lang = "bs"
	German             Lang = "de"
	Greek              Lang = "el"
	English            Lang = "en"
	Spanish            Lang = "es"
	French             Lang = "fr"
	Croatian           Lang = "hr"
	Italian            Lang = "it"
	Dutch              Lang = "nl"
	Polish             Lang = "pl"
	Portuguese         Lang = "pt"
	Russian            Lang = "ru"
	Slovak             Lang = "sk"
	Swedish            Lang = "sv"
	Tetum              Lang = "tet"
	Turkish            Lang = "tr"
	Ukranian           Lang = "uk"
	PigLatin           Lang = "x-pig-latin"
	Chinese            Lang = "zh"
	TraditionalChinese Lang = "zh-tw"
)

// APICallsHeader is the HTTP Header that contains the number of API calls made by the given key for the current 24 period.
const APICallsHeader = "X-Forecast-API-Calls"

func fromJSON(jsonBlob []byte) (*Forecast, error) {
	var f Forecast

	err := json.Unmarshal(jsonBlob, &f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}