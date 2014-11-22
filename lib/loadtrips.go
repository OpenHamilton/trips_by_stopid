
package gtfs

import (
	"fmt"
	"encoding/csv"
	"archive/zip"
	"gopkg.in/mgo.v2"
	"strings"
	"strconv"
)

func loadtrips(f * zip.File) {
	fmt.Println("loadtrips()")

	// open and read 
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

	c := sess.DB("gtfs").C("trips")
	c.DropCollection()

	var tripid, serviceid, shapeid, routeid int;
	var pieces []string;
	var trip Trip;

	for _, g := range records[1:]{

		if g[5] == "4" {
			// serviceid = 4
			continue // 4 means no service
		} else {
			pieces = strings.Split(g[5], "_")
			serviceid, err = strconv.Atoi(pieces[2])
			if err != nil { panic(err); }
		}
		tripid, err = strconv.Atoi(g[6])
		if err != nil { panic(err); }
		// g[5] serviceid
		shapeid, err = strconv.Atoi(g[4])
		if err != nil { panic(err); }
		// g[3] headsign
		routeid, err = strconv.Atoi(g[1])
		if err != nil { panic(err); }

		trip = Trip{ tripid, serviceid, shapeid, g[3], routeid }

		err = c.Insert(trip)
		if err != nil { panic(err) }

	}

}
