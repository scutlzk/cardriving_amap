# amap — 高德地图 Web 服务 Go 封装 & 驾车分析工具

基于高德地图开放平台 Web 服务 API 的 Go 封装库，同时提供一个交互式 Web 应用，支持在地图上可视化分析驾车到达多个目的地的耗时。

![Alt text](%E5%B1%80%E9%83%A8%E6%88%AA%E5%8F%96_20260301_234628.png)
## 功能概览

### Go 库（根包 `amap`）

| 函数 | 文件 | 说明 |
|---|---|---|
| `DrivingDuration` / `GetDrivingDetail` | `driving.go` | 驾车路线规划（v5），返回最短用时与距离 |
| `GeoCode` | `geocode.go` | 地理编码：地址 → 经纬度 |
| `ReGeoCode` | `geocode.go` | 逆地理编码：经纬度 → 结构化地址 |
| `QueryDistrict` | `district.go` | 行政区域查询，支持递归返回下级行政区 |
| `ConvertCoordinate` | `coordinate.go` | 坐标转换（GPS / 百度 / 图吧 → 高德 GCJ-02） |

### Web 应用（`server/`）

- 交互式高德地图，支持配置多个目的地
- **单点查询**：点击地图任意位置，计算到所有目的地的驾车耗时
- **区域批量分析**：在地图上画圈，自动在圈内均匀采样 10 个点并批量计算
- **热力覆盖图**：根据历史数据自动渲染颜色覆盖层（绿色=近，红色=远）
- 所有计算结果以结构化 JSON 格式追加写入 `results.json`

## 项目结构

```
amap/
├── driving.go          # 驾车路线规划
├── geocode.go          # 地理编码 / 逆地理编码
├── district.go         # 行政区域查询
├── coordinate.go       # 坐标转换
├── config.yaml         # 配置文件（API Key + 目的地列表）
├── example/
│   └── main.go         # 各 API 命令行示例
├── server/
│   ├── main.go         # Web 服务入口
│   └── index.html      # 地图前端页面（嵌入二进制）
├── go.mod
└── go.sum
```

## 快速开始

### 前置条件

- Go 1.22+
- 高德地图开放平台 Key（Web 服务 API 类型 + JS API 类型）

### 1. 配置

编辑 `config.yaml`：

```yaml
js_api:
  key: "你的 JS API Key"
  security_key: "你的安全密钥"

web_api:
  key: "你的 Web 服务 API Key"

destinations:
  - "上海市浦东新区陆家嘴"
  - "上海市徐汇区漕河泾开发区"
```

- `js_api`：用于前端地图渲染（[申请 JS API Key](https://lbs.amap.com/dev/)）
- `web_api`：用于后端路线规划、地理编码等接口
- `destinations`：目的地地址数组，支持任意数量

### 2. 启动 Web 服务

```bash
go run ./server/
```

启动后访问控制台输出的地址（默认 `http://localhost:6052`）。

### 3. 使用

- **单点分析**：直接点击地图，底部弹出到各目的地的驾车耗时
- **区域分析**：点击左上角「画圈分析」，在地图上拖拽画圈，自动计算圈内 10 个采样点到各目的地的耗时
- **撤销**：点击「撤销」清除当前画圈和标记
- 每次计算结果自动追加到项目根目录的 `results.json`

### 4. 命令行示例

```bash
go run ./example/
```

依次演示驾车规划、逆地理编码、地理编码、行政区域查询、坐标转换。

## API 文档

高德地图接口文档保存在项目中供参考：

| 文件 | 对应接口 |
|---|---|
| `carapi.md` | 驾车路线规划 v5 |
| `positionapi.md` | 地理编码 / 逆地理编码 v3 |
| `districtapi.md` | 行政区域查询 v3 |
| `coordinateapi.md` | 坐标转换 v3 |

## 依赖

- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML 配置解析
- [高德地图 JS API 2.0](https://lbs.amap.com/api/javascript-api-v2/summary) — 前端地图渲染
