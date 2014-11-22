
package gtfs

import(
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
)

func WriteStopsSearchJson(){
	fmt.Println("WriteStopsSearchJson() ");

	// dial for a mongo db session
	sess, err := mgo.Dial("localhost")
	if err != nil { panic(err) }
	defer sess.Close()
	sess.SetMode(mgo.Monotonic, true)

	c := sess.DB("gtfs").C("stops")

	var stops []Stops

	query := c.Find(nil)
	iter := query.Iter()
	err = iter.All(&stops)

	fileout := "/home/n/Go/src/github.com/nadeemelahi/gtfs/public/js/stopssearchable.js"

	fileOUT, err := os.Create(fileout)
	if err != nil { panic(err) }
	defer fileOUT.Close()

	fmt.Fprintf(fileOUT, `var stopssearchable = [`)


	for i, s := range stops {
		s.StopDesc = strings.Replace(s.StopDesc, "'", "\\'", -1)
		s.StopName = strings.Replace(s.StopName, "'", "\\'", -1)
		if i == 0 {
			fmt.Fprintf(fileOUT, "'%d, %s, %s'", s.StopId, s.StopDesc, s.StopName)
		}
		fmt.Fprintf(fileOUT, ",'%d, %s, %s'", s.StopId, s.StopDesc, s.StopName)
		// MUST DELETE THE LAST COMMA IN THE JS ARRAY AFTER RUNNING THIS SCRIPT
	}

	fmt.Fprintf(fileOUT, "];\n")


	fmt.Fprintf(fileOUT, `var stopscoordinates = {`)
	
	for i, s := range stops {
		if i == 0 {
			fmt.Fprintf(fileOUT, "%d:[%f,%f]", s.StopId, s.StopLat, s.StopLon)
		}
		fmt.Fprintf(fileOUT, ",%d:[%f,%f]", s.StopId, s.StopLat, s.StopLon)
	}
	fmt.Fprintf(fileOUT, "};")
}
