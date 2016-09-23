package darksky

import (
	"fmt"
	"time"
)

func ExampleGet() {
	r := MakeRequest("my_key", 41.8781, -87.6297).Get()

	feelsLike := r.Forecast.Currently.ApparentTemperature

	// All time fields are represented as seconds since epoch, so
	// conversion to a time.Time representation is straight forward.
	currentTime := time.Unix(r.Forecast.Currently.Time, 0)

	fmt.Printf("It feels like %v degrees outside at %v.", feelsLike, currentTime)
}
