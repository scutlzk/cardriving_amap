## 行政区域查询

#### 行政区域查询 API 服务地址

|   |   |
|---|---|
|URL|请求方式|
|https://restapi.amap.com/v3/config/district?parameters|GET|

parameters 代表的参数包括必填参数和可选参数。所有参数均使用和号字符(&)进行分隔。下面的列表枚举了这些参数及其使用规则。

#### 请求参数

|   |   |   |   |   |
|---|---|---|---|---|
|参数名|含义|规则说明|是否必须|缺省值|
|key|请求服务权限标识|用户在高德地图官网 [申请 Web 服务 API 类型 KEY](https://lbs.amap.com/dev/)|必填|无|
|keywords|查询关键字|规则：只支持单个关键词语搜索关键词支持：行政区名称、citycode、adcode<br><br>例如，在 subdistrict=2，搜索省份（例如山东），能够显示市（例如济南），区（例如历下区）<br><br>adcode 信息可参考 [城市编码表](https://lbs.amap.com/api/webservice/download) 获取|可选|无|
|subdistrict|子级行政区|规则：设置显示下级行政区级数（行政区级别包括：国家、省/直辖市、市、区/县、乡镇/街道多级数据）<br><br>可选值：0、1、2、3等数字，并以此类推<br><br>0：不返回下级行政区；<br><br>1：返回下一级行政区；<br><br>2：返回下两级行政区；<br><br>3：返回下三级行政区；<br><br>需要在此特殊说明，目前部分城市和省直辖县因为没有区县的概念，故在市级下方直接显示街道。<br><br>例如：广东-东莞、海南-文昌市|可选|1|
|page|需要第几页数据|最外层的 districts 最多会返回20个数据，若超过限制，请用 page 请求下一页数据。<br><br>例如：page=2；page=3。默认：page=1|可选|1|
|offset|最外层返回数据个数||可选|20|
|extensions|返回结果控制|此项控制行政区信息中返回行政区边界坐标点； 可选值：base、all;<br><br>base:不返回行政区边界坐标点；<br><br>all:只返回当前查询 district 的边界值，不返回子节点的边界值；<br><br>目前不能返回乡镇/街道级别的边界值|可选|base|
|filter|根据区划过滤|按照指定行政区划进行过滤，填入后则只返回该省/直辖市信息<br><br>需填入 adcode，为了保证数据的正确，强烈建议填入此参数|可选||
|callback|回调函数|callback 值是用户定义的函数名称，此参数只在 output=JSON 时有效|可选||
|output|返回数据格式类型|可选值：JSON，XML|可选|JSON|

#### 返回结果参数说明

行政区域查询的响应结果的格式由请求参数output指定。

|   |   |   |   |   |
|---|---|---|---|---|
|名称|   |   |含义|规则说明|
|status|   |   |返回结果状态值|值为0或1，0表示失败；1表示成功|
|info|   |   |返回状态说明|返回状态说明，status 为0时，info 返回错误原因，否则返回“OK”。|
|infocode|   |   |状态码|返回状态说明，10000代表正确，详情参阅 info 状态表|
|suggestion|   |   |建议结果列表||
||keywords|   |建议关键字列表||
|cities|   |建议城市列表||
|districts|   |   |行政区列表||
||district|   |行政区信息||
||citycode|城市编码||
|adcode|区域编码|街道没有独有的 adcode，均继承父类（区县）的 adcode|
|name|行政区名称||
|polyline|行政区边界坐标点|当一个行政区范围，由完全分隔两块或者多块的地块组<br><br>成，每块地的 polyline 坐标串以 \| 分隔 。<br><br>如北京 的 朝阳区|
|center|区域中心点|乡镇级别返回的center是边界线上的形点，其他行政级别返回的center不一定是中心点，若政府机构位于面内，则返回政府坐标，政府不在面内，则返回繁华点坐标。|
|level|行政区划级别|country:国家<br><br>province:省份（直辖市会在province显示）<br><br>city:市（直辖市会在province显示）<br><br>district:区县<br><br>street:街道|
|districts|下级行政区列表，包含 district 元素||

#### 服务示例

```
https://restapi.amap.com/v3/config/district?keywords=北京&subdistrict=2&key=<用户的key>
```

|参数|值|备注|必选|
|---|---|---|---|
|keywords||规则：只支持单个关键词语搜索关键词支持：行政区名称、citycode、adcode  <br>例如，在 subdistrict=2，搜索省份（例如山东），能够显示市（例如济南），区（例如历下区）|否|
|subdistrict|0                                                     1                                                     2                                                   3|规则：设置显示下级行政区级数（行政区级别包括：国家、省/直辖市、市、区/县4个级别）  <br>可选值：0、1、2、3  <br>0：不返回下级行政区；  <br>1：返回下一级行政区；  <br>2：返回下两级行政区；  <br>3：返回下三级行政区；|否|
|extensions|base                                                    all|此项控制行政区信息中返回行政区边界坐标点； 可选值：base、all;  <br>base:不返回行政区边界坐标点；  <br>all:只返回当前查询 district 的边界值，不返回子节点的边界值；|否|