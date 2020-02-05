//===== WEBSOCKET =====

var ws = null;
ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.addEventListener('message', function(e) {
  console.log("Received move update from server")
  var move = JSON.parse(e.data);
  console.log(move);
  // Add stone to board.
  placedStones[move.x][move.y] = move.color;
  console.log(placedStones);
});


//===== BOARD =======

var board = document.getElementById('board');
var ctx = board.getContext('2d');

var sizeMulti = 4;
var padding = 10*sizeMulti;
var width = 180*sizeMulti+padding*2, height = 180*sizeMulti+padding*2;
ctx.canvas.width = width;
ctx.canvas.height = height;

var stars = [ [3,3], [3,9], [3,15],
              [9,3], [9,9], [9,15],
              [15,3],[15,9],[15,15] ];

//=== STONES ===

var stoneRadius = padding/2;
var placedStones = [];
clear_board();

function clear_board() {
  for (var i = 0; i < 19; i++){
    placedStones[i] = [];
    for (var j = 0; j < 19; j++){
      placedStones[i][j] = null;
    }
  }
}

// Send a move request to the api
function sendMove(e) {
  /* ADD THIS TO SERVER
  var pos = getMousePos(e);
  var point = fromPixel(pos);
  if (inbounds(point)) {
    if (placedStones[point.x][point.y] == null) {
      if (blackTurn)
        placedStones[point.x][point.y] = "black";
      else
        placedStones[point.x][point.y] = "white";
      //console.log(placedStones);

      blackTurn = !blackTurn;
    }
  }
  */
  console.log("Sent move to server")
  var pos = getMousePos(e);
  var point = fromPixel(pos);
  if (inbounds(point)) {
    if (placedStones[point.x][point.y] == null) {
      ws.send(
        JSON.stringify({
          x: point.x,
          y: point.y
        })
      );
    }
  }
}

//===== OTHER =====

var blackTurn = true;

//===== DRAW ======

draw_board();
draw_lines();
draw_starpoints();

function draw_board() {
    // board bg
    ctx.fillStyle = '#ec9';
    ctx.fillRect(0, 0, width, height);
}

function draw_lines() {
  // board lines
  ctx.strokeStyle = '#320';
  for (var i = 0; i < 19; i++) {
    ctx.beginPath();
    ctx.moveTo(i*10*sizeMulti+padding, padding);
    ctx.lineTo(i*10*sizeMulti+padding, height-padding);
    ctx.stroke();
  }
  for (var j = 0; j < 19; j++) {
    ctx.beginPath();
    ctx.moveTo(padding, j*10*sizeMulti+padding);
    ctx.lineTo(width-padding, j*10*sizeMulti+padding);
    ctx.stroke();
  }
}

function draw_starpoints() {
  for (var star = 0; star < stars.length; star++) {
    var starpos = {x: stars[star][0], y: stars[star][1]};
    starpos = toPoint(starpos);
    ctx.strokeStyle = 'rgba(32,32,32,1)';
    ctx.fillStyle = "rgba(32,32,32,1)";
    ctx.beginPath();
    ctx.arc(starpos.x, starpos.y, 10*sizeMulti/10, 0, 2*Math.PI);
    ctx.stroke();
    ctx.fill();
  }
}

function draw_hover_stone(pos) {
  var point = fromPixel(pos);
  //console.log(point);
  if (inbounds(point)){
    if (placedStones[point.x][point.y] == null){
      var rounded = toPoint(point);
      //clamp
      rounded.x = clamp(rounded.x, padding, width-padding);
      rounded.y = clamp(rounded.y, padding, height-padding);

      // stone hover
      if (blackTurn) {
        ctx.strokeStyle = 'rgba(32,32,32,0.5)';
        ctx.fillStyle = 'rgba(32,32,32,0.5)';
      } else {
        ctx.strokeStyle = 'rgba(32,32,32,0.5)';
        //ctx.strokeStyle = 'rgba(255,255,255,0.5)';
        ctx.fillStyle = 'rgba(255,255,255,0.5)';
      }
      ctx.beginPath();
      ctx.arc(rounded.x, rounded.y, stoneRadius, 0, 2*Math.PI);
      ctx.stroke();
      ctx.fill();
    }
  }
}

function draw_placedStones() {
  for (var i = 0; i < 19; i++){
    for (var j = 0; j < 19; j++){
      if (placedStones[i][j] != null) {
        var pos = {x:i,y:j};
        pos = toPoint(pos);
        switch (placedStones[i][j]) {
          case "black":
            ctx.strokeStyle = 'rgba(32,32,32,1)';
            ctx.fillStyle = 'rgba(32,32,32,1)';
            break;
          case "white":
            //ctx.strokeStyle = 'rgba(255,255,255,1)';
            ctx.strokeStyle = 'rgba(32,32,32,1)';
            ctx.fillStyle = 'rgba(255,255,255,1)';
            break;
        }
        ctx.beginPath();
        ctx.arc(pos.x, pos.y, stoneRadius, 0, 2*Math.PI);
        ctx.stroke();
        ctx.fill();
      }

    }
  }
}

function draw(e) {
  var pos = getMousePos(e);

  draw_board();
  draw_lines();
  draw_starpoints();
  draw_placedStones();
  draw_hover_stone(pos);
}


function getMousePos(e) {
  var rect = board.getBoundingClientRect();
  return {
    x: e.clientX - rect.left,
    y: e.clientY - rect.top
  };
}

function inbounds(point) {
  return point.x >= 0 && point.x < 19 && point.y >= 0 && point.y < 19;
}

function clamp(num, min, max) {
  return Math.min(Math.max(num, min), max);
}

function fromPixel(pos) {
  return {
    x: Math.round(((pos.x - padding) / sizeMulti) / 10),
    y: Math.round(((pos.y - padding) / sizeMulti) / 10),
  }
}

function toPoint(pos) {
  return {
    x: ((pos.x *10) * sizeMulti) + padding,
    y: ((pos.y *10) * sizeMulti) + padding
  }
}