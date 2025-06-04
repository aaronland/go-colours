window.addEventListener("load", function load(event){
    
    sfomuseum.golang.wasm.fetch("wasm/extrude.wasm").then((rsp) => {

    	console.log("OKAY");
	
    }).catch((err) => {
	alert("Failed to load age WebAssembly functions");
	console.error("Failed to load WASMbinary", err);
        return;
    });
});
