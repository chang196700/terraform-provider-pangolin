param([string]$Version = '0.1.0')
$ErrorActionPreference = 'Stop'
$workspace = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$archive = Join-Path $workspace "dist/terraform-provider-pangolin_${Version}_windows_amd64.zip"
$archiveHash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.Substring(0,16)
$testRoot = Join-Path $workspace ".cache/mirror-smoke/$Version-$archiveHash"
$mirror = Join-Path $testRoot 'mirror'
$providerDir = Join-Path $mirror 'registry.terraform.io/chang196700/pangolin'
$module = Join-Path $testRoot 'module'
New-Item -ItemType Directory -Force -Path $providerDir,$module | Out-Null
Copy-Item -LiteralPath $archive -Destination $providerDir -Force
$cliConfig = Join-Path $testRoot 'terraform.rc'
$mirrorHcl = $mirror.Replace('\','/')
@"
provider_installation {
  filesystem_mirror {
    path = "$mirrorHcl"
    include = ["registry.terraform.io/chang196700/pangolin"]
  }
}
"@ | Set-Content -LiteralPath $cliConfig -Encoding utf8NoBOM
@"
terraform {
  required_version = ">= 1.11.0"
  required_providers {
    pangolin = {
      source = "chang196700/pangolin"
      version = "$Version"
    }
  }
}
provider "pangolin" {
  endpoint = "http://127.0.0.1:1/v1"
  api_key = "offline-validation-only"
}
"@ | Set-Content -LiteralPath (Join-Path $module 'provider.tf') -Encoding utf8NoBOM
$examples = Get-ChildItem -LiteralPath (Join-Path $workspace 'examples/resources') -Recurse -Filter resource.tf |
  Sort-Object FullName | ForEach-Object { Get-Content -Raw -LiteralPath $_.FullName }
$examples -join "`n" | Set-Content -LiteralPath (Join-Path $module 'examples.tf') -Encoding utf8NoBOM
$previousConfig = $env:TF_CLI_CONFIG_FILE
$env:TF_CLI_CONFIG_FILE = $cliConfig
try {
  terraform "-chdir=$module" init -backend=false -input=false -no-color
  if ($LASTEXITCODE -ne 0) { throw 'Mirror installation failed' }
  $schemaText = terraform "-chdir=$module" providers schema -json
  if ($LASTEXITCODE -ne 0) { throw 'Provider schema handshake failed' }
  $schema = $schemaText | ConvertFrom-Json
  $provider = $schema.provider_schemas.'registry.terraform.io/chang196700/pangolin'
  $resourceCount = @($provider.resource_schemas.PSObject.Properties).Count
  $sourceCount = @($provider.data_source_schemas.PSObject.Properties).Count
  if ($resourceCount -ne 24) { throw "Expected 24 resources, received $resourceCount" }
  terraform "-chdir=$module" validate -no-color
  if ($LASTEXITCODE -ne 0) { throw 'Generated examples failed Terraform validation' }
  Write-Output "Mirror verified: $resourceCount resources, $sourceCount data sources; all resource examples valid."
} finally {
  $env:TF_CLI_CONFIG_FILE = $previousConfig
}
