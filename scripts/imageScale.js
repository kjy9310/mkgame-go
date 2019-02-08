document.write('<input type="file" id="test"/>')

document.getElementById('test').onchange = function (evt) {
    var tgt = evt.target || window.event.srcElement,
        files = tgt.files;

    // FileReader support
    if (FileReader && files && files.length) {
        var fr = new FileReader();
        fr.onload = function () {
            console.log(fr.result.length)
            resizeBase64Img(fr.result,2500,2500).then(function (val){
                var img = document.createElement('img')
                img.src=val
                document.body.appendChild(img);

            })
        }
        fr.readAsDataURL(files[0]);
    }

    // Not supported
    else {
        // fallback -- perhaps submit the input to an iframe and temporarily store
        // them on the server until the user's session ends.
    }
}
async function resizeBase64Img(file, width, height) {
    var canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    var img = document.createElement('img');
    img.src = file;
    
    var x = await new Promise (function (resolve){
        img.onload = function(){
            var self = this;
            var context = canvas.getContext("2d");
            context.scale(width/self.width, height/self.height);
            context.drawImage(self, 0, 0);
            resolve(canvas.toDataURL());
        }   
    } );
    // console.log(x);
    return x;
}