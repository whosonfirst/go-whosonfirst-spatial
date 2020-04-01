var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.pip = (function(){

    var styles = {
	"match": {
	    "color": "#000",
	    "weight": 1,
	    "opacity": 1,
	    "fillColor": "#00308F",
	    "fillOpacity": 0.05
	}
    };
    
    var self = {

	'named_style': function(name){
	    return styles[name];
	},
	
	'default_properties': function(){

	    var props_table = {
		"wof:id":"",
		"wof:name":"",
		"wof:placetype":"",
		"wof:repo":"",
	    };

	    return props_table;
	},
	
	'render_properties_table': function(features, props_table){

	    if (! props_table){
		props_table = self.default_properties();
	    }
	    
	    var count = features.length;
	    
	    var table = document.createElement("table");
	    table.setAttribute("class", "table table-striped");
	    
	    var tr = document.createElement("tr");
	    
	    for (var k in props_table){
		
		var v = k;	// props_table[k]
		var th = document.createElement("th");
		th.appendChild(document.createTextNode(v));
		
		tr.appendChild(th);
	    }

	    var thead = document.createElement("thead");
	    thead.setAttribute("class", "thead-dark");
	    thead.appendChild(tr);

	    table.appendChild(thead);
	    
	    for (var i=0; i < count; i++){

		var f = features[i];
		var props = f["properties"];

		var wof_id = props["wof:id"];
		
		var tr = document.createElement("tr");
		tr.setAttribute("id", "tr-" + wof_id);
		
		for (var k in props_table){

		    var v = props[k];
		    var td = document.createElement("td");
		    
		    td.appendChild(document.createTextNode(v));
		    tr.appendChild(td);
		    table.appendChild(tr);
		}
		
	    }

	    var wrapper = document.createElement("div");
	    wrapper.setAttribute("class", "table-responsive");

	    wrapper.appendChild(table);
	    return wrapper;
	}
	
    };

    return self;
    
})();
