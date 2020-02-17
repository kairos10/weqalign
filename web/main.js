/****************************************************************/
//center of rotation: 2 stars(a, b), before and after rotation
function Point(x, y) {
   return {x: x, y: y}
}
function centerOfRot(a1, a2, b1, b2) {
   var midPointA = new Object(), midPointB = new Object()
   midPointA.x = (a1.x + a2.x) / 2
   midPointA.y = (a1.y + a2.y) / 2
   midPointB.x = (b1.x + b2.x) / 2
   midPointB.y = (b1.y + b2.y) / 2
   var slopeXA, slopeXB
   slopeXA = (a1.x-a2.x) / (a2.y-a1.y)
   slopeXB = (b1.x-b2.x) / (b2.y-b1.y)

   var x, y
   x = (midPointB.y - slopeXB*midPointB.x - midPointA.y + slopeXA*midPointA.x) / (slopeXA - slopeXB)
   y = midPointA.y - slopeXA*midPointA.x + slopeXA*x
   return { x: x, y: y }
}
/****************************************************************/

alert('main.js: wip')
