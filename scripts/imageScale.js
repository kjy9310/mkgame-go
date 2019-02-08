document.write('<input type="file" id="test"/>')

document.getElementById('test').onchange = function (evt) {
    var tgt = evt.target || window.event.srcElement,
        files = tgt.files;

    // FileReader support
    if (FileReader && files && files.length) {
        var fr = new FileReader();
        fr.onload = async function () {
            console.log(fr.result.length)
            size25 = await resizeBase64Img(fr.result,25,25);
            console.log(size25)
        }
        fr.readAsDataURL(files[0]);
    }

    // Not supported
    else {
        // fallback -- perhaps submit the input to an iframe and temporarily store
        // them on the server until the user's session ends.
    }
}
 function resizeBase64Img(file, width, height) {
    var canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    var img = document.createElement('img');
    img.src = file;
    return new Promise(function (resolve){
        img.onload = function(){
            var self = this;
            var context = canvas.getContext("2d");
            context.scale(width/self.width, height/self.height);
            context.drawImage(self, 0, 0);
            resolve(canvas.toDataURL());
        }    
    });
}