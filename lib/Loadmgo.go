

package gtfs 

import (
	"fmt"
	"archive/zip"
)

func Loadmgo(){
	fmt.Println("Loadmgo()");

	// Open a zip archive for reading.
	r, err := zip.OpenReader(Zippath); 
	if err != nil { panic(err) }
	defer r.Close();

	// Iterate files in the archive,
	for _, f := range r.File {

		switch f.Name {

			//case "routes.txt" : loadmgoroutes(f);
			//case "trips.txt" : loadmgotrips(f);
			//case "shapes.txt" : loadmgoshapes(f);

			// case "stops.txt" : loadstops(f); // loaded!
			// case "stop_times.txt" : loadstoptimes(f); // loaded!
			// case "calendar.txt" : loadcalendar(f); // loaded!
			// case "trips.txt" : loadtrips(f); // loaded!

		}
	}

}

