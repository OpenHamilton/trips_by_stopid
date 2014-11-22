

package main

import (
	"fmt"
	"github.com/nadeemelahi/gtfs/lib" 
)

func main(){

	fmt.Println("main() gtfs ");

	// download the zip file
	// gtfs.Downloadzipfile();

	// extract the data and load into mongo db
	// gtfs.Loadmgo();

	// gtfs.InspectTripidsForEachStopid()
	// so for each stop_id in stops.txt
	// there are multiple entries in stop_times.txt
	// for various trips_ids & departure times


	//gtfs.WriteStopsSearchJson()

	gtfs.Www()


}
