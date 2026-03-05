1、自定义图层 CustomLayer
1
1.1 准备
成为开发者并创建 key
为了正常调用 API ，请先注册成为高德开放平台开发者，并申请 web 平台（JS API）的 key 和安全密钥，点击 具体操作。

提示
你在2021年12月02日以后申请的 key 需要配合你的安全密钥一起使用。

2
1.2 创建一个自定义图层
自定义图层是完全由开发者指定绘制方法的图层。该图层可以是 canvas、svg、甚至可以是 dom 组成的图层。 JS API 能够实现自定义图层与高德地图的同步平移和缩放，并调用开发者的render方法进行图层的重绘。 自定义图层使用AMap.CustomLayer来进行构造，构造函数接受两个参数，第一个参数是作为图层的 dom 画布，第二个参数是图层的相关属性设定，与通用图层属性相同。以下为自定义图层的使用方法：

//创建 canvas
var canvas = document.createElement('canvas');

//将 canvas 宽高设置为地图实例的宽高
canvas.width = map.getSize().width;
canvas.height = map.getSize().height;

//创建一个自定义图层
var customLayer = new AMap.CustomLayer(canvas, {
  zIndex: 12,
  zooms: [3, 18] //设置可见级别，当地图级别在3到18之间时，此图层可见
});

//将自定义图层添加到地图 
map.add(customLayer);
JavaScript
3
1.3 自定义渲染方法
可使用 render 方法自定义图层渲染，你应该更新绘制时使用的容器内像素位置，来重新绘制图层内容。像素位置是由经纬度坐标转换而来，通常使用 map. lnglatToContainer 方法进行转换。

提示
render方法在自定义图层初次绘制、地图移动与缩放结束时调用。

//使用 canvas 在地图中心点绘制一个圆形
customLayer.render = () => {
  //获取地图中心点位置
  var center = map.getCenter(); //获取当前地图中心点经纬度坐标
  var pos = map.lngLatToContainer(center); //将地图经纬度坐标转为地图容器像素坐标
  var r = 20;

  //使用 canvas 绘制圆形
  var ctx = canvas.getContext("2d");
  ctx.fillStyle = "#08f";
  ctx.strokeStyle = "#fff";
  ctx.beginPath();
  ctx.moveTo(pos.x + r, pos.y);
  ctx.arc(pos.x, pos.y, r, 0, 2 * Math.PI);
  ctx.lineWidth = 3;
  ctx.closePath();
  ctx.stroke();
  ctx.fill();
};