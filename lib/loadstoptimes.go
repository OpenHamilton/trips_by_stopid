
package gtfs

import (
	"fmt"
	"archive/zip"
	"encoding/csv"
	"strconv"
	"gopkg.in/mgo.v2"
	"strings"
)

func loadstoptimes(f *zip.File){
	// f is our stop_times.txt
	fmt.Println("loadstoptimes()")

	// open and read routes.txt
	rc, err := f.Open()
	if err != nil { panic(err) }
	defer rc.Close()

	reader := csv.NewReader(rc)
	// read it line by line -records is an array
	records, err := reader.ReadAll()
	if err != nil { panic(err) }

	fmt.Println("lines: ", len(records))

	// dial for a mongo db session
	sess, err := mgo.Dial("localhost")
	if err != nil { panic(err) }
	defer sess.Close()
	sess.SetMode(mgo.Monotonic, true)

	c := sess.DB("gtfs").C("stoptimes")
	c.DropCollection()
	
	var pieces []string;
	var departureHr int;
	var departureMin int;

	for _, g := range records[1:]{


		tripId, err := strconv.Atoi(g[0])
		if err != nil { 
			fmt.Println("tripId conv"); 
			panic(err);
		}

		stopId, err := strconv.Atoi(g[3])
		if err != nil { 
			//fmt.Println("stopId conv"); 
			// empty
			stopId = 0
		}

		// g[2] is the time
		pieces = strings.Split(g[2], ":");

		departureHr, err = strconv.Atoi(pieces[0]);
		if err != nil { panic(err); }
		departureMin, err = strconv.Atoi(pieces[1]);
		if err != nil { panic(err); }

		fmt.Println(departureHr, departureMin);

		err = c.Insert(&Stoptimes{
			tripId,
			stopId,
			departureHr,
			departureMin,
		})

		if err != nil { panic(err) }

	}

	fmt.Println("loadstops() done")

}

