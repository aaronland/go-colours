window.addEventListener("load", function load(event){

    var upload_el = document.getElementById("upload");
    var feedback_el = document.getElementById("feedback");
    
    var image_btn = document.getElementById("image-button");
    var start_video_btn = document.getElementById("start-video");
    var stop_video_btn = document.getElementById("stop-video");            

    const video = document.getElementById("video");
    const colours = document.getElementById("colours");    
    
    sfomuseum.golang.wasm.fetch("wasm/extrude.wasm").then((rsp) => {

	var show_colours = function(rsp){
	    const data = JSON.parse(rsp);

	    const count = data.length;

	    colours.innerHTML = "";
	    
	    for (var i=0; i < count; i++){

		const swatches = data[i].swatches;
		const count_swatches = swatches.length;
		
		for (var k=0; k < count_swatches; k++){

		    const sw = data[i].swatches[k];
		    const hex = sw.colour.hex.trim("#");
		    
		    const d = document.createElement("div");
		    d.setAttribute("class", "swatch");
		    d.setAttribute("style", "background-color:" + hex);
		    d.appendChild(document.createTextNode(" "));

		    colours.appendChild(d);
		}
	    }
	};
	
	var derive_colours = function(im_b64){

	    const opts = {
		"grid": "euclidian://",
		"palettes": [ "crayola://" ],
		"extruders": [ "marekm4://" ],
	    };

	    const str_opts = JSON.stringify(opts);

	    colours_extrude(str_opts, im_b64).then((rsp) => {
		show_colours(rsp);
	    }).catch((err) => {
		feedback_el.innerText = "Failed to show colours, " + err;
	    });
	};

	var process_video_tick = function(){

	    if (video.readyState === video.HAVE_ENOUGH_DATA) {

		const canvas = document.createElement("canvas");
		const context = canvas.getContext('2d');
		
		canvas.width = video.videoWidth;
		canvas.height = video.videoHeight;
		
		context.drawImage(video, 0, 0, canvas.width, canvas.height);
		const im_b64 = canvas.toDataURL('image/jpeg');
		
		derive_colours(im_b64.replace("data:image/jpeg;base64,", ""));		
	    }
	    
	    requestAnimationFrame(process_video_tick);
	};
	
	var process_video = function(stream){

	    video.style.display = "block";
	    
	    video.srcObject = stream;
	    video.setAttribute("playsinline", true); // required to tell iOS safari we don't want fullscreen
	    video.play();
	    
	    requestAnimationFrame(process_video_tick);
	}
	
	var process_upload = function(){

	    if (! upload_el.files.length){
		feedback_el.innerText = "There are no files to process";
		return;
	    }
	    
	    const file = upload_el.files[0];
	    
	    if (! file.type.startsWith('image/')){
		return false;
	    }

	    switch (file.type) {
		case "image/jpeg":
		case "image/png":
		case "image/gif":
		case "image/webp":
		    // pass
		    break;
		default:
		    feedback_el.innerText = "Unsupported file type: " + file.type;
		    return;
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
		    const prefix = "data:" + file.type + ";base64,";
		    derive_colours(im_b64.replace(prefix, ""));
		};

		reader.readAsDataURL(file);
		
	    }, 10)
	    
	};

	upload_el.onchange = function(){
	    colours.innerHTML = "";
	};
	
	image_btn.onclick = function(){

	    feedback_el.innerHTML = "";
	    
	    try {
		process_upload();
	    } catch(err) {
		console.error(err);
	    }

	    return false;
	};

	start_video_btn.onclick = function(){

	    feedback_el.innerHTML = "";
	    
	    navigator.mediaDevices.getUserMedia({ video: { facingMode: "environment" } }).then(function(stream) {

		stop_video_btn.onclick = function(){

		    video.pause();
		    video.srcObject = null;
		    
		    stream.getTracks().forEach((track) => {
			if (track.readyState == 'live') {
			    track.stop();
			}
		    });

		    stop_video_btn.setAttribute("disabled", "disabled");
		    start_video_btn.removeAttribute("disabled");		    
		};

		start_video_btn.setAttribute("disabled", "disabled");		
		stop_video_btn.removeAttribute("disabled");
		
		process_video(stream);
		
	    }).catch((err) => {
		feedback_el.innerText = "Failed to start video feed, " + err;
	    });
	};
	
	upload_el.removeAttribute("disabled");
	image_btn.removeAttribute("disabled");
	start_video_btn.removeAttribute("disabled");	
	
    }).catch((err) => {
	feedback_el.innerText = "Failed to load age WebAssembly functions, " + err;
        return false;
    });
	
});
