
window.addEventListener("load", function(evt) {
    var canvas = document.getElementById("canvas");
	canvas.height=999;
	canvas.width=999;
    var canvasControl = function() {
	var Objects = [];
	function setObjects(objs) {
		Objects = objs;
	}
	function getObjects(){
		return Objects
	}
	function setObject(uuid, obj){
		Objects[Objects.findIndex(function(element){return element.Uuid==uuid})]=obj;
	}
	function getObject(uuid){
		return Objects.find(function(element){
			return element.Uuid==uuid;
		})
	}
	function drawCanvas(canvas) {
		var ctx = canvas.getContext('2d');
		ctx.clearRect(0, 0, canvas.width, canvas.height);
		for (var index in Objects){
			var object = Objects[index];
			var x = parseInt(object.Position/10000,10);
			var y = object.Position-x*10000;
			x+=Math.round(Math.sin(object.Direction)*object.Speed)
			y+=Math.round(Math.cos(object.Direction)*object.Speed*-1)
			if (object.Uuid == myUuid){
				ctx.fillStyle="blue";
			}
			ctx.fillRect(x,y,10,10)
			object.Position = x*10000+y;
			Objects[index] = object
			ctx.fillStyle="black";
		}
	}
	
	return {getObjects:getObjects,setObject:setObject, getObject:getObject, setObjects:setObjects, drawCanvas:drawCanvas}
    };
    var simulation = canvasControl();
    var simulatingInterval
    var output = document.getElementById("output");
    var type = document.getElementById("type");
    var input = document.getElementById("input");
    var statusSpan = document.getElementById("status");
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
	simulatingInterval = setInterval(function(){simulation.drawCanvas(canvas)},100);

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
		    simulation.setObjects(Object.values(jsonData.Data.Maps.start.Objects));
		    statusSpan.innerHTML = Object.values(jsonData.Data.Maps.start.Objects).map(function(object){
		    	return "uuid : "+object.Uuid+"<br/> Ap/Dp/Hp : "+object.Ap+"/"+object.Dp+"/"+object.Hp+" pos : "+object.Position+"<br/>"
		    }).join("<br/>");
		    break;
		case 'move':
		    var targetObj = simulation.getObject(jsonData.Uuid)
		    targetObj.Speed = jsonData.Data.Speed;
		    targetObj.Direction = jsonData.Data.Direction;
		    simulation.setObject(jsonData.Uuid, targetObj);
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
    var stopTimeout;
    document.getElementById("canvas").onclick = function(evt) {
        if (!ws) {
            return false;
        }
	clearTimeout(stopTimeout);
	var myObject = simulation.getObject(myUuid)
	var originX = parseInt(myObject.Position/10000,10);
	var originY = myObject.Position-originX*10000;
	var mouseX = evt.offsetX;
	var mouseY = evt.offsetY;
	var simulationObjects = simulation.getObjects()
	for(var index in simulationObjects){
		var targetObject = simulationObjects[index]
		if (myUuid==targetObject.Uuid){
			break
		}
		var x = parseInt(targetObject.Position/10000);
		var y = parseInt(targetObject.Position%10000);
		console.log("object:",x,y)
		if(mouseX>=x && mouseX<=x+10 && mouseY>=y && mouseY<=y+10){
			var sendData = {
				actionType:'attack',
				value:{
					Target:targetObject.Uuid,
				},
				time:((new Date()).getTime()-diffTime)
			};
			ws.send(JSON.stringify(sendData));
			return
		}
	}
	var diagonal = Math.sqrt(Math.pow(mouseX-originX,2)+Math.pow(mouseY-originY,2))
	var radian = (mouseX>=originX?1:-1)*Math.acos(-1*(mouseY-originY)/diagonal)
	myObject.Direction = Math.round(radian*100)/100
	simulation.setObject(myUuid, myObject)
        print("SEND: move with mouse");
        var sendData = {
		actionType:'move',
		value:{
			Direction:Math.round(radian*100)/100,
			Speed:3
		},
		time:((new Date()).getTime()-diffTime)
	};
	ws.send(JSON.stringify(sendData));
	stopTimeout = setTimeout(function(){
		var sendData = {
			actionType:'move',
			value:{
				Direction:Math.round(radian*100)/100,
				Speed:0
			},
			time:((new Date()).getTime()-diffTime)
		};
		ws.send(JSON.stringify(sendData));
	}, diagonal/3*100)

        return false;
    };
});
