var board = document.getElementById('board');
var ctx = board.getContext('2d');

var sizeMulti = 3;
var width = 180*sizeMulti+20, height = 180*sizeMulti+20;
ctx.canvas.width = width;
ctx.canvas.height = height;

// board bg
ctx.fillStyle = '#ec9';
ctx.fillRect(0, 0, width, height);

// board lines
ctx.strokeStyle = '#320';
for (var i = 0; i < 19; i++) {
  ctx.beginPath();
  ctx.moveTo(i*10*sizeMulti+10, 10);
  ctx.lineTo(i*10*sizeMulti+10, height-10);
  ctx.stroke();
}
for (var j = 0; j < 19; j++) {
  ctx.beginPath();
  ctx.moveTo(10, j*10*sizeMulti+10);
  ctx.lineTo(width-10, j*10*sizeMulti+10);
  ctx.stroke();
}