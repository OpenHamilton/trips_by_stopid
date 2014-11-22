
package gtfs

import(
	"fmt"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func stopSearchTripsAjaxH(res http.ResponseWriter, req *http.Request){

	fmt.Println("stopsSearchAjaxHandler() ")

	decoder := json.NewDecoder(req.Body)

	var stopSearchTrips StopSearchTrips

	err := decoder.Decode(&stopSearchTrips)
	if err != nil { 
		fmt.Println("failed post to stopSearch")
		fmt.Println(err)
		return
	}

	fmt.Println(stopSearchTrips.Stopid, stopSearchTrips.Day, stopSearchTrips.Hour)


	// dial for a mongo db session
	sess, err := mgo.Dial("localhost")
	if err != nil { panic(err) }
	defer sess.Close()
	sess.SetMode(mgo.Monotonic, true)

	cstoptimes := sess.DB("gtfs").C("stoptimes")

	// query by stopid, day or week, and time
	var stoptimes []Stoptimes;

	stopid, err := strconv.Atoi(stopSearchTrips.Stopid)
	if err != nil { panic(err) }

	stophr, err := strconv.Atoi(stopSearchTrips.Hour)
	if err != nil { panic(err) }

	query := cstoptimes.Find( bson.M{"stopid":stopid, "departuretimehr": stophr} )
	iter := query.Iter()
	err = iter.All(&stoptimes)
	if err != nil { panic(err); }

	ctrips := sess.DB("gtfs").C("trips")
	var trip Trip;

	m := make(map[int]string)
	m = map[int]string{942008: "sunday", 942004: "weekday", 942005: "sunday", 942006: "saturday", 942009: "saturday", 942007: "weekday", }

	r := make(map[int]string)
	r = map[int]string{
		2903:"08 YORK",
		2902:"07 LOCKE",
		2862:"51 UNIVERSITY",
		2863:"52 DUNDAS LOCAL",
		2864:"55 STONEY CREEK CENTRAL",
		2865:"56 CENTENNIAL",
		2848:"20 A-LINE EXPRESS",
		2849:"21 UPPER KENILWORTH",
		2846:"16 ANCASTER",
		2847:"18 WATERDOWN",
		2844:"11 PARKDALE",
		2845:"12 WENTWORTH",
		2842:"09 ROCK GARDENS",
		2843:"10 B-LINE EXPRESS",
		2905:"02 BARTON",
		2926:"33 SANATORIUM",
		2930:"43 STONE CHURCH",
		2860:"43 STONE CHURCH",
		2913:"10 B-LINE EXPRESS",
		2922:"24 UPPER SHERMAN",
		2907:"04 BAYFRONT",
		2928:"35 COLLEGE",
		2906:"03 CANNON",
		2866:"58 STONEY CREEK LOCAL",
		2910:"07 LOCKE",
		2923:"25 UPPER WENTWORTH",
		2867:"99 WATERFRONT SHUTTLE",
		2904:"01 KING",
		2931:"44 RYMAL",
		2861:"44 RYMAL",
		2921:"23 UPPER GAGE",
		2935:"56 CENTENNIAL",
		2920:"22 UPPER OTTAWA",
		2918:"20 A-LINE EXPRESS",
		2909:"06 ABERDEEN",
		2934:"55 STONEY CREEK CENTRAL",
		2859:"41 MOHAWK",
		2858:"35 COLLEGE",
		2912:"09 ROCK GARDENS",
		2908:"05 DELAWARE",
		2924:"26 UPPER WELLINGTON",
		2915:"12 WENTWORTH",
		2916:"16 ANCASTER",
		2917:"18 WATERDOWN",
		2851:"23 UPPER GAGE",
		2850:"22 UPPER OTTAWA",
		2853:"25 UPPER WENTWORTH",
		2852:"24 UPPER SHERMAN",
		2855:"27 UPPER JAMES",
		2854:"26 UPPER WELLINGTON",
		2857:"34 UPPER PARADISE",
		2856:"33 SANATORIUM",
		2932:"51 UNIVERSITY",
		2927:"34 UPPER PARADISE",
		2839:"06 ABERDEEN",
		2838:"05 DELAWARE",
		2837:"04 BAYFRONT",
		2836:"03 CANNON",
		2835:"02 BARTON",
		2834:"01 KING",
		2936:"58 STONEY CREEK LOCAL",
		2925:"27 UPPER JAMES",
		2919:"21 UPPER KENILWORTH",
		2933:"52 DUNDAS LOCAL",
		2937:"99 WATERFRONT SHUTTLE",
		2911:"08 YORK",
		2929:"41 MOHAWK",
		2914:"11 PARKDALE",
	};


	var searchresponse StopSearchTripsResponse
	var jresp []StopSearchTripsResponse

	for _, stimes := range stoptimes {

		fmt.Println("-------------")
		fmt.Println(stopid, stimes.TripId)

		

		err = ctrips.Find( bson.M{ "tripid": stimes.TripId } ).One(&trip);
		if err != nil { 
			// its not guaranteed that all tripid's in stoptimes are found in trips, so if not found just skip.
			continue
			//fmt.Println("ctrips Find error for tripid: ", stimes.TripId); 
		} else {
			fmt.Println("length of trip", trip)
		}

		

		fmt.Println(stopSearchTrips.Day, m[trip.ServiceId], trip.TripId, trip.ServiceId, trip.ShapeId, trip.Headsign, trip.Routeid)

		


		
		if stopSearchTrips.Day == m[trip.ServiceId] {
			searchresponse = StopSearchTripsResponse{
				stimes.TripId,
				trip.Routeid,
				r[trip.Routeid],
				trip.Headsign,
				stimes.DepartureTimeHr,
				stimes.DepartureTimeMin,
			}
			jresp = append(jresp, searchresponse);
		}

		// for each tripid in stimes
		// cross check with calendar
		// for stopSearchTrips.Day
		
	}


	
	for _, j := range jresp { 
		fmt.Println( j.RouteId, j.HeadSign, j.DepartureTimeHr, j.DepartureTimeMin )
	}

	fmt.Println("results count: ", len(jresp))
	

	encoder := json.NewEncoder(res)


	err = encoder.Encode(jresp)
	if err != nil { panic(err); }


}
