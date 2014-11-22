

window.addEventListener("load", function(){

	var w = window.innerWidth;
	var h = window.innerHeight;

	var map;

	var options = { 
		zoom: 11,
center: new google.maps.LatLng(43.25, -79.87),
mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	var mapDiv = document.getElementById("map");
	map = new google.maps.Map(mapDiv, options);

	// wrapper div to hold the control
	var searchDiv = document.createElement('div');

	searchDiv.style.padding = '5px';
	searchDiv.style.marginBottom = '5px';



	var dnow = new Date();
	var schedDay = dnow.getDay();

	if ( schedDay == 0 ){
		schedDay = "sunday";
	} else if ( schedDay == 6 ) {
		schedDay = "saturday";
	} else schedDay = "weekday";

	//schedDay = "sunday";

	var option, textnode;
	var dayofweekUI = document.createElement('select');
	dayofweekUI.style.padding = "5px";

	option = document.createElement('option');
	option.value = schedDay;
	textnode = document.createTextNode("weekday");
	option.appendChild(textnode);
	if ( schedDay == "weekday") {
		option.setAttribute("selected", "selected");
	}
	dayofweekUI.appendChild(option);

	option = document.createElement('option');
	option.value = "saturday";
	textnode = document.createTextNode("saturday");
	option.appendChild(textnode);
	if ( schedDay == "saturday") {
		option.setAttribute("selected", "selected");
	}
	dayofweekUI.appendChild(option);

	option = document.createElement('option');
	option.value = "sunday";
	textnode = document.createTextNode("sunday");
	option.appendChild(textnode);
	if ( schedDay == "sunday") {
		option.setAttribute("selected", "selected");
	}
	dayofweekUI.appendChild(option);

	searchDiv.appendChild(dayofweekUI);

	var hourofdayUI = document.createElement('select');
	hourofdayUI.style.padding = "5px";

	var hourOfDay = dnow.getHours();

	var ampm, time, i;
	for ( i = 4; i < 28; i++ ){
		ampm = "pm";
		if ( i < 12 || i > 23 ) { ampm = "am"; }

		time = i;
		if ( i > 12 ) { time = i - 12; }

		if ( i > 24 ){
			time = i - 24;
		}

		option = document.createElement('option');
		option.value = i;

		if ( i == hourOfDay || ( i > 24 && time == hourOfDay ) || (i == 24 && hourOfDay == 0) ) {
			option.setAttribute("selected", "selected");
		}

		textnode = document.createTextNode(time + " " + ampm);
		option.appendChild(textnode);
		hourofdayUI.appendChild(option);

	}

	searchDiv.appendChild(hourofdayUI);


	var searchUI = document.createElement('input');

	searchUI.style.padding = "5px";
	searchUI.setAttribute("type", "text");
	searchUI.setAttribute("id", "myinput");
	searchUI.setAttribute("placeholder", "search bus stop");

	searchDiv.appendChild(searchUI);


	google.maps.event.addDomListener(searchUI, 'keyup', function(){
		handleInput(this);
	});

	map.controls[google.maps.ControlPosition.BOTTOM_LEFT].push(searchDiv);

	// wrapper to hold results div
	var resultsDiv = document.createElement('div');
	resultsDiv.style.padding = '20px';
	resultsDiv.style.display = 'none';
	resultsDiv.index = -1; //ensures atop defaultUI
	//does seem to work

	var resultsInnerDiv = document.createElement('div');
	resultsInnerDiv.style.backgroundColor = 'rgba(' + 255 + ',' + 255 + ',' + 255 + ',' + 0.3 + ')';
	resultsInnerDiv.style.borderStyle = 'solid';
	resultsInnerDiv.style.borderWidth = '1px';
	resultsInnerDiv.style.borderColor = '#acacac';
	resultsInnerDiv.style.cursor = 'pointer';
	resultsInnerDiv.style.title = 'select stop';
	resultsInnerDiv.style.width = 0.8*w + "px";
	resultsInnerDiv.style.height = 0.8*h + "px";

	resultsInnerDiv.id = "scrollable";

	resultsDiv.appendChild(resultsInnerDiv);

	google.maps.event.addDomListener(resultsDiv, 'click');

	map.controls[google.maps.ControlPosition.TOP_RIGHT].push(resultsDiv);


	searchUI.addEventListener('focus', function(){
		clearpolysMarkers();
		resultsDiv.style.display = "block";
	}, false);


	function handleInput(inputDiv){
		var inputString= inputDiv.value;

		//var prevString = inputDiv.getAttribute("data-prevString");
		//inputDiv.setAttribute("data-prevString", inputString); 
		//  no need to prevent hammering

		if ( inputString.length > 2 ) {
			findMatches(inputString);
		}
	}


	var day, hr;
	var busstopstrings = [];
	var sslength = stopssearchable.length;
	var substrRegex = [], qpieces = [], qplength, i, j;
	function findMatches(q) {
		//console.log("input search query: ", q);
		qpieces = q.split(" ");
		qplength = qpieces.length;

		for ( i = 0; i < qplength; i++ ){
			if ( qpieces[i] == "and" || qpieces[i] == "AND" ){
				qpieces[i] = "&";
			}
			substrRegex[i] = new RegExp(qpieces[i], 'i');
		}	

		var matches = [], allmatched, str;
		for ( i = 0; i < sslength; i++ ){
			allmatched = true;
			str = stopssearchable[i];
			for ( j = 0; j < qplength; j++ ){
				if ( !(substrRegex[j].test(str)) ){
					allmatched = false;
				}
			}
			if ( allmatched ){
				matches.push(str);
			}
		}

		var html = '';
		for ( i = 0 ; i < matches.length; i++ ){
			//console.log(matches[i]);
			html += '<p class="busstopstring">' + matches[i] + '</p>';
		}
		resultsInnerDiv.innerHTML = html;

		busstopstrings = document.getElementsByClassName("busstopstring");
		var busstopstringslength = busstopstrings.length;
		for ( i = 0; i < busstopstringslength; i++ ){
			busstopstrings[i].addEventListener("click", function(){
				day = dayofweekUI.value;
				hr = hourofdayUI.value;
				console.log(this.innerHTML);
				placemarker(this.innerHTML);
				resultsDiv.style.display = "none";
			}, false);
		}
	};



	var markers = [], l = [], mark, jmsg; // 0-depart, 1-destination
	function placemarker(stop){
		var pieces = (stop).split(", ");
		var stopid = pieces[0]; //2267
		//console.log(stopscoordinates[stopid]);
		l[0] = stopscoordinates[stopid][0];
		l[1] = stopscoordinates[stopid][1];

		jmsg = {"Stopid":stopid,"Day":day,"Hour":hr};
		jmsg = JSON.stringify(jmsg);
		console.log(jmsg);

		// l[0] is lat, l[1] is lng
		var loc = new google.maps.LatLng(l[0],l[1]);
		var mark = new google.maps.Marker({
			position: loc, title: stopid,
		    icon: "img/flag.png", map: map
		});

		var loc = new google.maps.LatLng(l[0],l[1]);

		//center map to new marker
		map.setOptions({center:loc});

		xmlhttp = new XMLHttpRequest();
		xmlhttp.onreadystatechange = function(){
			if (xmlhttp.readyState==4 && xmlhttp.status==200){
				//console.log('done');
				pasteTrips(xmlhttp.responseText, l);
			}
		};
		xmlhttp.open("POST","searchStopsAjax",true);
		xmlhttp.send(jmsg);


		if (markers.length > 1){ //max 2 flags
			markers[0].setMap(null);
			markers.shift();
		}
		markers.push(mark);



		var triplinks, triplinkslength;
		var stophtml, stophr, stopmin, inf, loc;

		function pasteTrips(res, l){
			res = JSON.parse(res);
			console.log(res);

			triplinks = document.getElementsByClassName("triplink");
			triplinkslength = triplinks.length;
			// so as to avoid registering click handler twice
			for ( i = 0; i < triplinkslength; i++ ) { 
				//console.log(i, triplinks[0]);
				triplinks[0].className = "tripl"; //
				// it does shift() after changeing className
				// so always use index 0
			}

			stophtml = "<div class='info'><h4>" + stop + "</h4>";

			if ( res == null ) {
				stophtml += "<p>Sorry, no buses running.</p>";
			} else {

				var reslength = res.length;

				for ( i = 0; i < reslength; i++ ){
					//console.log(res[i]);
					stophr = res[i].DepartureTimeHr;
					if ( stophr < 10 ) { stophr = "0" + stophr; }
					stopmin = res[i].DepartureTimeMin;
					if ( stopmin < 10 ) { stopmin = "0" + stopmin; }

					stophtml += "<p id='" + res[i].TripId + "' class='triplink'>[ " + stophr + ":" + stopmin + " ] " + res[i].HeadSign + " (" + res[i].RouteName + ")</p>";

				}

			} 

			stophtml += "</div>";
			//console.log(stophtml);

			inf = new google.maps.InfoWindow({
				content: stophtml 
			});


			google.maps.event.addListener(inf, "domready", function(){

				triplinks = document.getElementsByClassName("triplink");
				triplinkslength = triplinks.length;

				for ( i = 0; i < triplinkslength; i++ ){
					triplinks[i].addEventListener("click", function(){
						console.log(this, this.id);
					}, false);

				}

			});


			loc = new google.maps.LatLng(l[0],l[1]);

			//center map to new marker
			map.setOptions({center:loc});

			// inf.open(map, mark); // auto open
			// seems to open and position poorly

			regOpenIHand(inf, mark); // register InfoWindow Open Handler

			function regOpenIHand(i, m){

				google.maps.event.addListener(m, 'click', function(){ 
					i.open(map, m); 
				});

				i.open(map, m);

			}


		}


	}







	var polys = [], ptspoly = [], j;

	function makePolygon(ptsgps){

		ptspoly = [];

		for ( j = 0; j < 4; j++ ){
			ptspoly.push( new google.maps.LatLng(ptsgps[j][0], ptsgps[j][1]) );
		}

		polys.push(new google.maps.Polygon({
			path: ptspoly, map: map, 
			strokeColor: '#000',
			strokeOpacity: 0.9,
			strokeWeight: 0.1,
			fillColor: '#000',
			fillOpacity: 0 }) );
	}

	var x1, y1, x2, y2;
	x1 = 43.143082;
	x2 = 43.415024;
	y1 = -80.087585;
	y2 = -79.554749;

	//   00.100000
	x1 = 43.153082; // + north
	y1 = -80.032585; // - west

	pts1gps = [ [x1,y1], [x2, y1], [x2, y2], [x1, y2] ];
	pts2gps = [ [x1,y1], [x2, y1], [x2, y2], [x1, y2] ];
	pts3gps = [ [x1,y1], [x2, y1], [x2, y2], [x1, y2] ];


	var ptss = [], ptsgps = [], gridsize = 0.015;
	//ptss = [ pts1gps ];

	for ( m = 1; m < 24; m++ ){
		for ( n = 0; n < 14; n++ ){
			x2 = x1 + gridsize;
			y2 = y1 + gridsize;
			ptss.push([ [x1, y1], [x2, y1], [x2, y2], [x1, y2] ]);
			x1 = x2;
		}
		x1 = 43.1530832;
		y1 = -80.032585 + m*gridsize;
		//y1 = y2;
	}

	ptssLength = ptss.length;

	for ( i = 0; i < ptssLength; i++ ){ 

		makePolygon(ptss[i]);

	}


	var mark, matchedpolysMarkers = [], k;
	var matchedpolysMarkersLength = matchedpolysMarkers.length;

	function clearpolysMarkers(){
		matchedpolysMarkersLength = matchedpolysMarkers.length;
		for ( k = 0; k < matchedpolysMarkersLength; k++ ){
			matchedpolysMarkers[k].setMap(null);
		}

	}

	function addStopMarker(stopid, stoppt){
		mark = new google.maps.Marker({ 
			position: stoppt, map: map, 
		     icon: "img/flag.png", title: stopid });
		matchedpolysMarkers.push(mark);
		//register marker click infoWindow
		registerMarkerClick(mark, stoppt);

	}

	function registerMarkerClick(mark, stoppt){

		google.maps.event.addListener(mark, 'click', function() {
			clearpolysMarkers();

			day = dayofweekUI.value;
			hr = hourofdayUI.value;
			stopid = this.title;
			jmsg = {"Stopid":stopid,"Day":day,"Hour":hr};
			jmsg = JSON.stringify(jmsg);
			//console.log(jmsg);

			xmlhttpGrid = new XMLHttpRequest();
			xmlhttpGrid.onreadystatechange = function(){
				if (xmlhttpGrid.readyState==4 && xmlhttpGrid.status==200){
					//console.log('done');
					pasteTripsGrid(xmlhttpGrid.responseText, mark, stopid, stoppt);
				}
			};
			xmlhttpGrid.open("POST","searchStopsAjax",true);
			xmlhttpGrid.send(jmsg);

		});

	}

	var stopTitle;
	function pasteTripsGrid(res, mark, stopid, stoppt){

		for ( i = 0; i < sslength; i++ ){
			qpieces = (stopssearchable[i]).split(",");
			//console.log(qpieces[0]);	
			if ( stopid == qpieces[0] ) {
				stopTitle = stopssearchable[i];
				//console.log(stopTitle);
				break;
			}
		}

		res = JSON.parse(res);
		//console.log(mark.position, mark.title, mark.icon);

		mark = new google.maps.Marker({
			position: mark.position, title: mark.title,
		     icon: "img/flag.png", map: map });
		// need replacement, so it won't get deleted 
		// when clicking second marker in another polygon

		if (markers.length > 1){ //max 2 flags
			markers[0].setMap(null);
			markers.shift();
		}
		markers.push(mark);
		// make a infoWindow and paste the res
		triplinks = document.getElementsByClassName("triplink");
		triplinkslength = triplinks.length;
		// so as to avoid registering click handler twice
		for ( i = 0; i < triplinkslength; i++ ) { 
			//console.log(i, triplinks[0]);
			triplinks[0].className = "tripl"; //
			// it does shift() after changeing className
			// so always use index 0
		}

		stophtml = "<div class='info'><h4>" + stopTitle + "</h4>";

		if ( res == null ) {
			stophtml += "<p>Sorry, no buses running.</p>";
		} else {

			var reslength = res.length;

			for ( i = 0; i < reslength; i++ ){
				stophr = res[i].DepartureTimeHr;
				if ( stophr < 10 ) { stophr = "0" + stophr; }
				stopmin = res[i].DepartureTimeMin;
				if ( stopmin < 10 ) { stopmin = "0" + stopmin; }

				stophtml += "<p id='" + res[i].TripId + "' class='triplink'>[ " + stophr + ":" + stopmin + " ] " + res[i].HeadSign + " (" + res[i].RouteName + ")</p>";

			}

		}

		stophtml += "</div>";
		//console.log(stophtml);

		inf = new google.maps.InfoWindow({
			content: stophtml 
		});


		google.maps.event.addListener(inf, "domready", function(){

			triplinks = document.getElementsByClassName("triplink");
			triplinkslength = triplinks.length;

			for ( i = 0; i < triplinkslength; i++ ){
				triplinks[i].addEventListener("click", function(){
					console.log(this, this.id);
				}, false);

			}

		});

		map.setOptions({center: stoppt});

		regOpenIGridHand(inf, mark);

		function regOpenIGridHand(i, m){

			google.maps.event.addListener(m, 'click', function(){ 
				i.open(map, m); 
			});

			i.open(map, m);

		}



	}

	var polysLength = polys.length, 
	    stoppt, matchedpolys = [];

	function regPolyClickHand(poly){
		//console.log(poly);
		google.maps.event.addListener(poly, 'click', function(e){
			map.fitBounds(poly.getBounds());

			matchedpolys = [];

			clearpolysMarkers();

			for ( stopid in stopscoordinates ) {
				stoppt = new google.maps.LatLng(stopscoordinates[stopid][0], stopscoordinates[stopid][1]);
				if ( google.maps.geometry.poly.containsLocation(stoppt, poly) ) {
					matchedpolys.push(stopid);
					addStopMarker(stopid, stoppt);
				}
			}
			//console.log("markers count: ", matchedpolys.length); // 100 is too many
			//console.log("all markers count: ", stopssearchable.length); // 100 is too many


		});


	}


	for ( i = 0; i < polysLength; i++ ){
		regPolyClickHand(polys[i]);
	}


	google.maps.Polygon.prototype.getBounds = function() {
		var bounds = new google.maps.LatLngBounds();
		var paths = this.getPaths();
		var path;        
		for (var i = 0; i < paths.getLength(); i++) {
			path = paths.getAt(i);
			for (var ii = 0; ii < path.getLength(); ii++) {
				bounds.extend(path.getAt(ii));
			}
		}
		return bounds;
	}











}, false);

