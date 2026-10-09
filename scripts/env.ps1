$env:GOPATH = Join-Path $PSScriptRoot '..\.cache\go'
$env:GOMODCACHE = Join-Path $env:GOPATH 'pkg\mod'
$env:GOCACHE = Join-Path $PSScriptRoot '..\.cache\go-build'
$env:UV_CACHE_DIR = Join-Path $PSScriptRoot '..\.cache\uv'
$env:UV_PYTHON_INSTALL_DIR = Join-Path $PSScriptRoot '..\.cache\python'
$env:TMP = Join-Path $PSScriptRoot '..\.cache\tmp'
$env:TEMP = $env:TMP
New-Item -ItemType Directory -Force -Path $env:TMP | Out-Null
