## 坐标转换

#### 坐标转换 API 服务地址：

|   |   |
|---|---|
|URL|请求方式|
|https://restapi.amap.com/v3/assistant/coordinate/convert?parameters|GET|

parameters 代表的参数包括必填参数和可选参数。所有参数均使用和号字符(&)进行分隔。下面的列表枚举了这些参数及其使用规则。

#### 请求参数

|   |   |   |   |   |
|---|---|---|---|---|
|参数名|含义|规则说明|是否必须|缺省值|
|key|请求服务权限标识|用户在高德地图官网 [申请 Web 服务 API 类型 KEY](https://lbs.amap.com/dev/)|必填|无|
|locations|坐标点|经度和纬度用","分割，经度在前，纬度在后，经纬度小数点后不得超过6位。多个坐标对之间用”\|”进行分隔最多支持40对坐标。|必填|无|
|coordsys|原坐标系|可选值：<br><br>gps;<br><br>mapbar;<br><br>baidu;<br><br>autonavi(不进行转换)|可选|autonavi|
|sig|数字签名|请参考 [数字签名获取和使用方法](https://lbs.amap.com/faq/quota-key/key/41181/)|可选|无|
|output|返回数据格式类型|可选值：JSON,XML|可选|JSON|

#### 返回结果参数说明

坐标转换的响应结果的格式由请求参数 output 指定。

|   |   |   |
|---|---|---|
|名称|含义|规则说明|
|status|返回状态|值为0或1<br><br>1：成功；0：失败|
|info|返回的状态信息|status 为0时，info 返回错误原；否则返回“OK”。详情参阅 [info 状态表](https://lbs.amap.com/api/webservice/guide/tools/info/)|
|locations|转换之后的坐标。若有多个坐标，则用 “;”进行区分和间隔||

#### 服务示例

```
https://restapi.amap.com/v3/assistant/coordinate/convert?locations=116.481499,39.990475&coordsys=gps&output=xml&key=<用户的key>
```

|参数|值|备注|必选|
|---|---|---|---|
|locations||坐标点,经度和纬度用“,”分割，经度在前，纬度在后，经纬度小数点后不得超过6位。多个坐标对之间用“\|”进行分隔最多支持40对坐标。|是|
|coordsys||原坐标系，可选值：gps;mapbar;baidu;autonavi(不进行转换)  <br>默认：autonavi|否|