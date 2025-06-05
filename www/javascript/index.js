window.addEventListener("load", function load(event){

    var upload_el = document.getElementById("upload");
    var image_btn = document.getElementById("image-button");
    var start_video_btn = document.getElementById("start-video");
    var stop_video_btn = document.getElementById("stop-video");            

    const video = document.getElementById("video");
    const colours = document.getElementById("colours");    
    
    sfomuseum.golang.wasm.fetch("wasm/extrude.wasm").then((rsp) => {

	var show_colours = function(rsp){
	    const data = JSON.parse(rsp);
	    console.log("colours", data);

	    const count = data.length;

	    for (var i=0; i < count; i++){

		// const ext = data[i];

		const swatches = data.swatches;
		const count_swathces = swatches.length;
		
		for (var k=0; k < count_swatches; k++){

		    const sw = swatches[k];
		    const hex = sw.colour.hex.trim("#");
		    
		    console.log(hex);
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
		console.log("SAD", err);
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
	
	image_btn.onclick = function(){

	    try {
		process_upload();
	    } catch(err) {
		console.error(err);
	    }

	    return false;
	};

	start_video_btn.onclick = function(){

	    navigator.mediaDevices.getUserMedia({ video: { facingMode: "environment" } }).then(function(stream) {

		stop_video_btn.onclick = function(){

		    video.pause();

		    stream.getTracks().forEach((track) => {
			if (track.readyState == 'live') {
			    track.stop();
			}
		    });

		    stop_video_btn.setAttribute("disabled", "disabled");
		};

		start_video_btn.setAttribute("disabled", "disabled");		
		stop_video_btn.removeAttribute("disabled");
		
		process_video(stream);
		
	    }).catch((err) => {
		console.error(err);
	    });
	};
	
	upload_el.removeAttribute("disabled");
	image_btn.removeAttribute("disabled");
	start_video_btn.removeAttribute("disabled");	
	
    }).catch((err) => {
	alert("Failed to load age WebAssembly functions");
	console.error("Failed to load WASMbinary", err);
        return;
    });
	
});
