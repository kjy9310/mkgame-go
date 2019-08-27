
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var type = document.getElementById("type");
    var input = document.getElementById("input");
    var ws;
    var pingFunc;
    var pingSent = (new Date()).getTime();
    var diffTime = 0;
    var print = function(message) {
        var d = document.createElement("div");
            d.innerHTML = message;
            output.appendChild(d);
        };
   
     document.getElementById("open").onclick = function(evt) {
	var hostAddress = document.getElementById("address").value;
	if (ws) {
	    return false;
	}
	ws = new WebSocket("ws://"+hostAddress+"/ws");
	ws.onopen = function(evt) {
	    print("OPEN");
	}
	ws.onclose = function(evt) {
	    print("CLOSE");
	    ws = null;
	    clearInterval(pingFunc);
	    pingFunc = null;
	}
	ws.onmessage = function(evt) {
	    if (!evt.data) {
		return
	    }
	    var jsonData = JSON.parse(evt.data);
	    var serverTime = jsonData.Time;
	    switch (jsonData.Action){
	        case 'pong':
		    diffTime = pingSent - serverTime;
		    break;
	    }
	    print("RESPONSE: " + evt.data);
	}
	ws.onerror = function(evt) {
	    print("ERROR: " + evt.data);
	}
	pingFunc = setInterval(sendPingRequest, 5000);
	return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        var sendData = {
		actionType:type.value,
		value:input.value,
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

    function sendPingRequest() {
	if (!ws) {
	    return;
	}
	pingSent = (new Date()).getTime();
        var pingData = {
		actionType:'ping',
		time:(pingSent-diffTime)
	};
	ws.send(JSON.stringify(pingData));
    }
});
