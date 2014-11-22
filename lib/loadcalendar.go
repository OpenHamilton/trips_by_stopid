
package gtfs

import (
	"fmt"
	"encoding/csv"
	"archive/zip"
	"gopkg.in/mgo.v2"
	"strings"
	"strconv"
)

func loadcalendar(f * zip.File) {
	fmt.Println("loadcalendar()")

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

	c := sess.DB("gtfs").C("calendars")
	c.DropCollection()

	var pieces []string;
	var wksatsun int;
	var serviceid int;
	var calendar Calendar;
	for _, g := range records[1:]{

		if g[0] == "4" {
			//pieces[2] = "4"
			continue // 4 is no service
		} else {
			pieces = strings.Split(g[0], "_")
		}
		
		serviceid, _ = strconv.Atoi(pieces[2]);
		if err != nil { panic(err) }

		fmt.Println(g[3], g[8], g[9], g)

		if g[3] == "1" {
			wksatsun = 1 // weekday
			fmt.Println("weekday", wksatsun)
		} else if g[8] == "1" {
			wksatsun = 2 // sat
			fmt.Println("sat", g[8])
		} else if g[9] == "1" {
			wksatsun = 3 // sun
			fmt.Println("sun", g[9])
		}

		calendar = Calendar{ serviceid, wksatsun }

		err = c.Insert(calendar)
		if err != nil { panic(err) }

	}


}
