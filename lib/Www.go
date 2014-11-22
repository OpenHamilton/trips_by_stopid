
package gtfs

import(
	"fmt"
	"net/http"
)

func Www(){
	fmt.Println("Wwww() ")
	http.HandleFunc("/searchStopsAjax", stopSearchTripsAjaxH)
	http.Handle("/", http.FileServer(http.Dir("/home/n/Go/src/github.com/nadeemelahi/gtfs/public/")))
	http.ListenAndServe(":8082", nil)
}
