# DarkSky API for Go

## Requirements

* Go 1.5+
* Valid API key from https://darksky.net/dev.

## Usage

    import "go.larrymyers.com/darksky"
    
    r := darksky.MakeRequest("my_key", "41.8781", "-87.6297").Get()
    
    feelsLike := r.Forecast.Currently.ApparentTemperature

    // All time fields are represented as seconds since epoch, so
    // conversion to a time.Time representation is straight forward.
    currentTime := time.Unix(r.Forecast.Currently.Time, 0)

    fmt.Printf("It feels like %v degrees outside at %v.", feelsLike, currentTime)

## Notes

All time based fields are stored as int64 values, which contain the seconds since epoch.

Conversion can be done using time.Unix.

## Run Tests With Coverage

    go test -coverprofile=cover.out && go tool cover -html=cover.out
