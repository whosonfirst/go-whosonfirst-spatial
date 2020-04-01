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

    var layers = L.layerGroup();
    layers.addTo(map);

    map.on("moveend", function(e){

	var pos = map.getCenter();	

	var args = {
	    'latitude': pos['lat'],
	    'longitude': pos['lng'],
	    'format': 'geojson',
	};

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

	    var table = whosonfirst.spatial.pip.render_properties_table(features);
	    
	    var matches = document.getElementById("pip-matches");
	    matches.innerHTML = "";
	    matches.appendChild(table);
	    
	};

	var on_error = function(err){
	    console.log("SAD", err);
	}

	whosonfirst.spatial.api.point_in_polygon(args, on_success, on_error);
    });
    
    map.setView([37.604, -122.405], 13);

    slippymap.crosshairs.init(map);    
});
