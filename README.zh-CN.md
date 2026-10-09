# Pangolin Terraform Provider

面向 Pangolin 1.24.0 的社区 Provider，地址为 `chang196700/pangolin`。支持 Terraform 1.11 及以上版本，使用 Go 1.25 和 Plugin Framework 1.19 开发。

实现 24 类资源，覆盖组织、站点、公网/私网资源、Target、访问策略、用户与角色、OIDC、密钥、告警、健康检查和 AI 网关。单项和列表数据源另外覆盖域名、DNS、证书、Blueprint、邀请、用户设备、模型目录。全部 312 个 API 操作的映射或排除原因见 [接口覆盖清单](docs/api-coverage.md)。

**1.24.0 存在实际 API 限制。** 部分 OpenAPI 声明没有挂载到 API Key 接口，部分处理程序依赖浏览器会话。因此包含 24 类资源模型不代表未修改的 1.24.0 能运行所有生命周期。采用前请查看 [兼容性表](docs/compatibility.md) 和 [权限清单](docs/permissions.md)；模拟测试不能证明上游缺失接口可用。

## 使用

安装示例和完整说明见 [英文 README](README.md)；Registry 尚未完成发布时，使用本地镜像安装。`endpoint` 必须包含 `/v1`，密钥可通过 `PANGOLIN_API_KEY` 提供。每个资源显式填写 `org_id`，不同组织不会隐式共用默认组织。

`PANGOLIN_API_KEY` 需要完整的 `apiKeyId.apiKeySecret`。组织密钥资源的 `api_key` 属性只有 secret 部分，使用时需要和该资源的 `id` 拼接。

```hcl
provider "pangolin" {
  endpoint = "https://api.example.com/v1"
}

resource "pangolin_site" "office" {
  org_id = "your-org-id"
  name   = "Office"
  type   = "newt"
}
```

资源参数、导入格式和示例见 [资源文档](docs/resources)。组织级 API Key 按实际操作授权；组织、全局身份提供方、组织 API Key 管理需要 Root Key。部分接口依赖商业版功能；Root Key 也无法使用上游未开放或仅供浏览器会话调用的接口。

## 已明确的行为与限制

- 授权集合由 Terraform 完整管理；删除集合中的条目会撤销对应授权，系统强制管理员权限除外。
- 公共资源内联 `access` 与共享 `resource_policy_id` 互斥。1.24.0 的白名单接口可能直接修改共享策略，不能安全地声称支持独立覆盖；共享认证配置统一放在策略资源里。
- Newt 的密钥、组织 API Key 等仅在创建时返回。刷新会保留已保存的值，导入无法恢复明文。Sensitive 只是隐藏输出，密钥仍保存在 Terraform state 中。
- 组织用户创建仅支持 OIDC；删除资源只移出组织。删除角色要指定接收用户的 `reassign_role_id`。
- 虚拟 API Key 的普通更新不会延长有效期。健康检查、监控、SSO 等根据实际服务端版本和权限提供能力。
- 不将重启、发送邀请或邮件、Blueprint 应用、日志导出、2FA 设置流程伪装成 Terraform 持久资源。域名缺少创建和删除接口，因此只提供查询。

## 构建与验证

```powershell
. ./scripts/env.ps1
go test ./...
go vet ./...
$env:TF_ACC = '1'
$env:TF_ACC_TERRAFORM_PATH = (Get-Command terraform).Source
go test ./internal/provider -run TestAcc -count=1 -timeout 30m
go run ./cmd/release -version 0.1.0
./scripts/verify-mirror.ps1
```

验收使用真实 Terraform CLI 对接本地 HTTP 模拟服务，验证生命周期、空计划、导入、远端删除恢复、策略切换及密钥行为。默认不访问真实 Pangolin。显式启用真实测试的方式见英文 README；真实实例联调结果不能由模拟测试替代。

发布流程见 [发布说明](docs/releasing.md)。Linux、macOS、Windows 的 amd64/arm64 安装包统一输出到 `dist/`，发布前必须签名校验和并在 Registry 注册公钥。
