var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.pip = (function(){

    var self = {

	'default_properties': function(){

	    var props_table = {
		"wof:id":"",
		"wof:name":"",
		"wof:placetype":"",
	    };

	    return props_table;
	},
	
	'render_properties_table': function(features, props_table){

	    if (! props_table){
		props_table = self.default_properties();
	    }
	    
	    var count = features.length;
	    
	    var table = document.createElement("table");
	    table.setAttribute("class", "table");

	    var tr = document.createElement("tr");

	    for (var k in props_table){
		
		var v = k;	// props_table[k]
		var th = document.createElement("th");
		th.appendChild(document.createTextNode(v));
		
		tr.appendChild(th);
		table.appendChild(tr);
	    }
	    
	    for (var i=0; i < count; i++){

		var f = features[i];
		var props = f["properties"];

		var tr = document.createElement("tr");
		
		for (var k in props_table){

		    var v = props[k];
		    var td = document.createElement("td");
		    
		    td.appendChild(document.createTextNode(v));
		    tr.appendChild(td);
		    table.appendChild(tr);
		}
		
	    }

	    return table;
	}
	
    };

    return self;
    
})();
