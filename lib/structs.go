package gtfs

var Zippath string = "/home/n/Go/src/github.com/nadeemelahi/gtfs/zipfile/hsrgtfs.zip";

type Trip struct {
	TripId int
	ServiceId int
	ShapeId int
	Headsign string
	Routeid int
}

type Calendar struct { // only 7 total, and 6 saved since one is out of service
	ServidId int // ex 1_, 2_, 3_merged_942008
	// we only hold that last number 942008
        ServiceDay int // set 1=weekday, 2=sat, 3=sun
}

type Stops struct {
	StopId int
	StopDesc string
	StopName string

	StopLat float32 // probably these not needed
	StopLon float32
}

type Stoptimes struct {
	TripId int
	StopId int
	DepartureTimeHr int
	DepartureTimeMin int
}

type StopSearchTrips struct { // decode $.post 
	Stopid string 
	Day string
	Hour string
}

type StopSearchTripsResponse struct {
	TripId int
	RouteId int
	RouteName string
	HeadSign string
	DepartureTimeHr int
	DepartureTimeMin int
}


