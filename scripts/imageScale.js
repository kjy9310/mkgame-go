document.write('<input type="file" id="test"/><button onclick="Fclick()">image</button>')

function Fclick() {
    var tgt = document.getElementById('test'),
        files = tgt.files;

    // FileReader support
    if (FileReader && files && files.length) {
        var fr = new FileReader();
        fr.onload = async function () {
            console.log(fr.result.length)
            const a = await resizeBase64Img(fr.result,250,250);
            var img = document.createElement('img')
            img.src=a;
            document.body.appendChild(img);
        }
        fr.readAsDataURL(files[0]);
    }

    // Not supported
    else {
        alert("unsupported browser!! :S");
    }
}
async function resizeBase64Img(file, width, height) {
    var canvas = document.createElement("canvas");
    var img = document.createElement('img');
    img.src = file;
    var x = await new Promise (function (resolve){
        img.onload = function(){
            var self = this;
            var context = canvas.getContext("2d");
            var w_ratio = width/self.width
            var h_ratio = height/self.height
            if (width===0){
                canvas.width = img.width*h_ratio;
                w_ratio = h_ratio;
            }else{
                canvas.width = width;
            }
            if (height==0){
                canvas.height = img.height*w_ratio;
                h_ratio = w_ratio
            }else{
                canvas.height = height
            }
            context.scale(w_ratio, h_ratio);
            context.drawImage(self, 0, 0);
            resolve(canvas.toDataURL());
        }   
    } );
    // console.log(x);
    return x;
}