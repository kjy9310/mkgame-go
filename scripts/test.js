
window.addEventListener("load", function(evt) {
    var canvas = document.getElementById("canvas");
	canvas.height=999;
	canvas.width=999;
    var canvasControl = function() {
	var Objects = [];
	function setObjects(objs) {
		Objects = objs;
	}
	function drawCanvas(canvas) {
		var ctx = canvas.getContext('2d');
		ctx.clearRect(0, 0, canvas.width, canvas.height);
		Objects.forEach(function(object){
			var x = parseInt(object.Position/10000,10);
			var y = object.Position-x*10000;
			if (object.Uuid == myUuid){
				ctx.fillStyle="blue";
			}
			ctx.fillRect(x,y,10,10)
			ctx.fillStyle="black";
		});
	}
	
	return {setObjects:setObjects, drawCanvas:drawCanvas}
    };
    var simulation = canvasControl();
    var simulatingInterval
    var output = document.getElementById("output");
    var type = document.getElementById("type");
    var input = document.getElementById("input");
    var ws;
    var myUuid;
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
	simulatingInterval = setInterval(function(){simulation.drawCanvas(canvas)},500);

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
		case 'status':
		    simulation.setObjects(Object.values(jsonData.Data.Objects.start));
		    break;
		case 'login':
		    myUuid = jsonData.Data.Uuid;
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
	clearInterval(simulatingInterval)
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

    document.getElementById("left").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: left");
        var sendData = {
		actionType:'move',
		value:{
			Direction:4.7,
			Speed:1
		},
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
        return false;
    };
    document.getElementById("right").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: right");
        var sendData = {
		actionType:'move',
		value:{
			Direction:1.5,
			Speed:1
		},
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
        return false;
    };
    document.getElementById("up").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: up");
        var sendData = {
		actionType:'move',
		value:{
			Direction:0,
			Speed:1
		},
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
        return false;
    };
    document.getElementById("down").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: down");
        var sendData = {
		actionType:'move',
		value:{
			Direction:3.14,
			Speed:1
		},
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
        return false;
    };
});
