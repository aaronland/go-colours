window.addEventListener("load", function load(event){

    var upload_el = document.getElementById("upload");
    var submit_el = document.getElementById("submit");    
    
    sfomuseum.golang.wasm.fetch("wasm/extrude.wasm").then((rsp) => {

	var derive_colours = function(im_b64){

	    const opts = {
		"grid": "euclidian://",
		"palettes": [ "crayola://" ],
		"extruders": [ "marekm4://" ],
	    };

	    const str_opts = JSON.stringify(opts);

	    colours_extrude(str_opts, im_b64).then((rsp) => {
		console.log("OK", rsp);
	    }).catch((err) => {
		console.log("SAD", err);
	    });
	};
	
	var process_upload = function(){

	    const file = upload_el.files[0];
	    
	    if (! file.type.startsWith('image/')){
		return false;
	    }
	    
            const reader = new FileReader();

            reader.onload = function(e) {
		
		const img = document.createElement("img");
		img.setAttribute("style", "max-height: 400px; max-width:400px;");
		img.src = e.target.result;
		
		const wrapper = document.getElementById("image-wrapper");
		wrapper.innerHTML = "";
		wrapper.appendChild(img);
            };
	    
            reader.readAsDataURL(file);

	    setTimeout(function(){

		reader.onload = function(e) {
		    const im_b64 = e.target.result;
		    derive_colours(im_b64.replace("data:image/jpeg;base64,", ""));
		};

		reader.readAsDataURL(file);
		
	    }, 10)
	    
	};
	
	submit_el.onclick = function(){

	    try {
		process_upload();
	    } catch(err) {
		console.error(err);
	    }

	    return false;
	};
	    
	upload_el.removeAttribute("disabled");
	submit_el.removeAttribute("disabled");
	
    }).catch((err) => {
	alert("Failed to load age WebAssembly functions");
	console.error("Failed to load WASMbinary", err);
        return;
    });
	
});
