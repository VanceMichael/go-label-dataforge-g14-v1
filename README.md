# DataForge 公共数据工场运营平台

DataForge 是面向资源提供方、审核员、数据产品开发者和场景运营者的纯后端平台，覆盖资源存证登记、公示赋码发证，以及目录授权、沙箱租约、模型测试和产品发布两条相互依赖的生命周期。

## 运行

```bash
GOTOOLCHAIN=local go run ./cmd/server
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

默认使用 SQLite 文件数据库，启动时执行版本化 migration。服务提供会话登录、角色鉴权、请求 ID、审计、outbox 和可恢复 worker。

## 主要 API

`POST /v1/sessions` 登录，`DELETE /v1/sessions/{id}` 注销；`POST /v1/resources` 创建资源，`POST /v1/resources/{id}/submit` 提交审核，`POST /v1/reviews/{id}/decision` 审核；`POST /v1/authorizations` 申请授权，`POST /v1/authorizations/{id}/approve` 批准；`POST /v1/leases` 创建沙箱租约，`POST /v1/leases/{id}/runs` 运行模型；`POST /v1/products` 创建产品，`POST /v1/products/{id}/publish` 发布。

## 数据与恢复

数据库表包括 tenants、users、sessions、data_resources、resource_versions、registration_reviews、authorization_requests、sandbox_leases、sandbox_runs、data_products、product_releases、audit_events、outbox_jobs。事务边界覆盖登记、授权批准、租约回收和产品发布。worker 通过数据库租约、重试次数和退避恢复未完成 outbox、审核超时、租约到期与发布任务。
