
package gtfs

import (
	"fmt"
	"archive/zip"
	"encoding/csv"
	"strconv"
	"gopkg.in/mgo.v2"
	//"strings"
)

func loadstops(f *zip.File){
	// f is our stop_times.txt
	fmt.Println("loadstops()")

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

	c := sess.DB("gtfs").C("stops")
	c.DropCollection()

	for _, g := range records[1:]{

		stopId, err := strconv.Atoi(g[4])
		if err != nil { 
			// some stopId's may have been merged
			// ie 324324_merged_32432
			// fmt.Println("stopId conv: ", g[4]); 
			// we'll just skip those for now cause we don't know what they mean
			continue // break current iteration but continue the looping
		}

		stopLat, err := strconv.ParseFloat(g[0], 32)
		if err != nil { 
			fmt.Println("stopLat conv"); 
			panic(err) 
		}
		stopLat32 := float32(stopLat)

		stopLon, err := strconv.ParseFloat(g[3], 32)
		if err != nil { 
			fmt.Println("stopLon conv"); 
			panic(err) 
		}
		stopLon32 := float32(stopLon)

		stop := Stops{ 
			stopId,
			g[7], // stopId
			g[8], // stopname
			stopLat32, 
			stopLon32, 
		}

		//fmt.Println(stop)

		err = c.Insert(stop)

		if err != nil { panic(err) }

		//fmt.Println("i: ", i)
	}

	fmt.Println("loadstops() done")

}

