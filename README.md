# dasel

[![Gitbook](https://badges.aleen42.com/src/gitbook_1.svg)](https://daseldocs.tomwright.me)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub All Releases Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)](https://github.com/TomWright/dasel/releases/latest)
[![Homebrew tag (latest by date)](https://img.shields.io/homebrew/v/dasel)](https://formulae.brew.sh/formula/dasel)

Dasel(data-selector 的缩写)允许您使用选择器字符串查询和修改数据结构。

和 [jq](https://github.com/stedolan/jq) / [yq](https://github.com/kislyuk/yq) 比起来，dasel 支持 JSON, YAML, TOML, XML 和 CSV。

## One tool to rule them all

不要再为了使用不同的数据格式而学习新工具了。

无论数据格式如何，Dasel 都使用标准的选择器语法。这意味着，一旦您学会了如何使用 dasel，您就可以立即查询/修改任何受支持的数据类型，
而无需任何额外的工具或工作。

![Update Kubernetes Manifest](update_kubernetes.gif)

## Table of contents
* [Dasel](#dasel)
* [One tool to rule them all](#one-tool-to-rule-them-all)
* [Quickstart](#quickstart)
* [Issue vs discussion](#issue-vs-discussion)
* [Features](#features)
* [Table of contents](#table-of-contents)
* [Documentation](#documentation)
* [Playground](#playground)
* [Benchmarks](#benchmarks)

## Quickstart

Dasel is available on [homebrew](https://daseldocs.tomwright.me/installation#homebrew), [ASDF](https://daseldocs.tomwright.me/installation#asdf), [scoop](https://daseldocs.tomwright.me/installation#scoop), [docker](https://daseldocs.tomwright.me/installation#docker), [Nix](https://daseldocs.tomwright.me/installation#nix) or as [compiled binaries](https://daseldocs.tomwright.me/installation#manual) from the [latest release](https://github.com/TomWright/dasel/releases/latest).

```bash
brew install dasel
```

You can also install a [development version](https://daseldocs.tomwright.me/installation#development-version) with:
```bash
go install github.com/tomwright/dasel/cmd/dasel@master
```

For more information see the [installation documentation](https://daseldocs.tomwright.me/installation).

### Select

```bash
echo '{"name": "Tom"}' | dasel -r json '.name'
"Tom"
```

See [select documentation](https://daseldocs.tomwright.me/usage/select).

### 将 json 转为 yaml

```bash
echo '{"name": "Tom"}' | dasel -r json -w yaml
name: Tom
```

查看 [select documentation](https://daseldocs.tomwright.me/usage/select).

### Put

```bash
echo '{"name": "Tom"}' | dasel put string -r json '.email' 'contact@tomwright.me'
{
  "email": "contact@tomwright.me",
  "name": "Tom"
}
```

See [put documentation](https://daseldocs.tomwright.me/usage/put).

### Delete

```bash
echo '{
  "email": "contact@tomwright.me",
  "name": "Tom"
}' | dasel delete -p json '.email'
{
  "name": "Tom"
}
```

See [delete documentation](https://daseldocs.tomwright.me/usage/delete).

## Issue vs Discussion

I have enabled [discussions](https://github.com/TomWright/dasel/discussions) on this repository.

I am aware there may be some confusion when deciding where you should communicate when reporting issues, asking questions or raising feature requests so this section aims to help us align on that.

Please [raise an issue](https://github.com/TomWright/dasel/issues) if:
- You find a bug.
- You have a feature request and can clearly describe your request.

Please [open a discussion](https://github.com/TomWright/dasel/discussions) if:
- You have a question.
- You're not sure how to achieve something with dasel.
- You have an idea but don't quite know how you would like it to work.
- You have achieved something cool with dasel and want to show it off.
- Anything else!

## Features
- [Query/select data from structured data files](https://daseldocs.tomwright.me/usage/select).
- [Update data in structured data files](https://daseldocs.tomwright.me/usage/put).
- [Create data files](https://daseldocs.tomwright.me/usage/put#create-documents-from-scratch).
- [Supports multiple data formats/types](https://daseldocs.tomwright.me/usage/supported-file-types).
- [Convert between data formats/types](https://daseldocs.tomwright.me/notes/converting-between-formats).
- Uses a [standard query/selector syntax](https://daseldocs.tomwright.me/selectors/introduction) across all data formats.
- Zero runtime dependencies.
- [Available on Linux, Mac and Windows](https://daseldocs.tomwright.me/installation).
- Available to [import and use in your own projects](https://daseldocs.tomwright.me/use-as-a-go-package).
- [Run via Docker](https://daseldocs.tomwright.me/installation#docker).
- [Faster than jq/yq](#benchmarks).
- [Self update](https://daseldocs.tomwright.me/installation/update).

## Documentation

The official dasel docs can be found at [daseldocs.tomwright.me](https://daseldocs.tomwright.me).

## Playground

You can test out dasel commands using the [playground](https://dasel.tomwright.me).

Source code for the playground can be found at  [github.com/TomWright/daselplayground](https://github.com/TomWright/daselplayground).

## Benchmarks

In my tests dasel has been up to 3x faster than jq and 15x faster than yq.

See the [benchmark directory](./benchmark/README.md).
