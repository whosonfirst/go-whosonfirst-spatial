window.addEventListener("load", function load(event){

    var api_key = document.body.getAttribute("data-nextzen-api-key");
    var style_url = document.body.getAttribute("data-nextzen-style-url");
    var tile_url = document.body.getAttribute("data-nextzen-tile-url");    

    if (! api_key){
	console.log("Missing API key");
	return;
    }
    
    if (! style_url){
	console.log("Missing style URL");
	return;
    }
    
    if (! tile_url){
	console.log("Missing tile URL");
	return;
    }

    var pip_wrapper = document.getElementById("point-in-polygon");

    if (! pip_wrapper){
	console.log("Missing 'point-in-polygon' element.");
	return;
    }
    
    var init_lat = pip_wrapper.getAttribute("data-initial-latitude");

    if (! init_lat){
	console.log("Missing initial latitude");
	return;
    }
    
    var init_lon = pip_wrapper.getAttribute("data-initial-longitude");

    if (! init_lon){
	console.log("Missing initial longitude");
	return;
    }
    
    var init_zoom = pip_wrapper.getAttribute("data-initial-zoom");    

    if (! init_zoom){
	console.log("Missing initial zoom");
	return;
    }
    
    var map_el = document.getElementById("map");

    if (! map_el){
	console.log("Missing map element");	
	return;
    }

    var map_args = {
	"api_key": api_key,
	"style_url": style_url,
	"tile_url": tile_url,
    };

    // we need to do this _before_ Tangram starts trying to draw things
    // map_el.style.display = "block";
    
    var map = whosonfirst.spatial.maps.getMap(map_el, map_args);

    if (! map){
	console.log("Unable to instantiate map");
	return;
    }

    var hash = new L.Hash(map);
    
    var layers = L.layerGroup();
    layers.addTo(map);

    var filter_ids = [
	"is_current",
	"is_deprecated",
	"is_ceased",
	"is_superseded",
	"is_superseding",
    ];

    var count_filters = filter_ids.length;
    	
    var update_map = function(e){

	var pos = map.getCenter();	

	// PLEASE MAKE ME DYNAMIC (or at least not hard-coded here)
	
	var properties = [
	    "sfomuseum:placetype",
	    "edtf:inception",
	    "edtf:cessation",	    
	];
	
	var args = {
	    'latitude': pos['lat'],
	    'longitude': pos['lng'],
	    'format': 'geojson',
	    'properties': properties.join(","),
	};

	for (var i=0; i < count_filters; i++){
	    
	    var f = filter_ids[i];
	    var el = document.getElementById(f);

	    if ((el) && (el.checked)){
		args[f] = 1;
	    }
	}

	var on_success = function(rsp){

	    layers.clearLayers();
	    
	    var l = L.geoJSON(rsp, {
		style: function(feature){
		    return whosonfirst.spatial.pip.named_style("match");
		},
	    });
	    
	    layers.addLayer(l);
	    l.bringToFront();

	    //
	    
	    var features = rsp["features"];

	    var table_props = whosonfirst.spatial.pip.default_properties();

	    var count_properties = properties.length;

	    for (var i=0; i < count_properties; i++){
		table_props[properties[i]] = "";
	    }
	    
	    var table = whosonfirst.spatial.pip.render_properties_table(features, table_props);
	    
	    var matches = document.getElementById("pip-matches");
	    matches.innerHTML = "";
	    matches.appendChild(table);
	    
	};

	var on_error = function(err){
	    console.log("SAD", err);
	}

	whosonfirst.spatial.api.point_in_polygon(args, on_success, on_error);
    };
    
    map.on("moveend", update_map);

    for (var i=0; i < count_filters; i++){
	    
	var f = filter_ids[i];
	var el = document.getElementById(f);
	
	if (! el){
	    continue
	}

	el.onchange = update_map;
    }

    
    var hash_str = location.hash;

    if (hash_str){

	var parsed = whosonfirst.spatial.maps.parseHash(hash_str);

	if (parsed){
	    init_lat = parsed['latitude'];
	    init_lon = parsed['longitude'];
	    init_zoom = parsed['zoom'];
	}
    }
    
    map.setView([init_lat, init_lon], init_zoom);    

    slippymap.crosshairs.init(map);    
});
