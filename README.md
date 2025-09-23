# gostyle

[![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/gostyle.svg)](https://pkg.go.dev/github.com/k1LoW/gostyle) [![build](https://github.com/k1LoW/gostyle/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/gostyle/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/gostyle/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/gostyle/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/gostyle/time.svg)

`gostyle` is a set of analyzers for coding styles.

## Disclaimer

`gostyle` **IS NOT** [Go Style](https://google.github.io/styleguide/go/).

"Go Style" in [Google Style Guides](https://google.github.io/styleguide) is a great style and we will actively refer to it, but we cannot implement the same rules perfectly, and we may extend the rules.

`gostyle` **IS NOT STANDARD**.

`gostyle` **IS** to help you maintain **YOUR** Go project coding **STYLE**.

## Usage

### As a Standalone CLI

```console
$ gostyle run ./...
```

### As a vet tool

```console
$ go vet -vettool=`which gostyle` ./...
```

### On GitHub Actions

**:octocat: GitHub Actions for gostyle is [here](https://github.com/k1LoW/gostyle-action) !!**

``` yaml
# .github/workflows/ci.yml
on:
  push:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
[...]
      -
        uses: k1LoW/gostyle-action@v1
        with:
          config-file: .gostyle.yml
[...]
```

## Analyzers

**Although not perfect**, it provides analyzers based on helpful styles.

### [Effective Go](https://go.dev/doc/effective_go)

- [ifacenames](analyzer/effective/ifacenames) ... based on https://go.dev/doc/effective_go#interface-names

> ["Effective Go"](https://go.dev/doc/effective_go) by The Go Authors is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)

### [Go Style](https://google.github.io/styleguide/go/) in Google Style Guides

- [**Guide**](https://google.github.io/styleguide/go/guide)
  - [mixedcaps](#mixedcaps) ... based on https://google.github.io/styleguide/go/guide#mixed-caps
- [**Decisions**](https://google.github.io/styleguide/go/decisions)
  - [funcfmt](#funcfmt) ... based on https://google.github.io/styleguide/go/decisions#function-formatting
  - [getters](#getters) ... based on https://google.github.io/styleguide/go/decisions#getters
  - [nilslices](#nilslices) ... based on https://google.github.io/styleguide/go/decisions#nil-slices
  - [pkgnames](#pkgnames) ... based on https://google.github.io/styleguide/go/decisions#package-names
  - [recvnames](#recvnames) ... based on https://google.github.io/styleguide/go/decisions#receiver-names
  - [recvtype](#recvtype) ... based on https://google.github.io/styleguide/go/decisions#receiver-type
  - [repetition](#repetition) ... based on https://google.github.io/styleguide/go/decisions#repetition
  - [typealiases](#typealiases) ... based on https://google.github.io/styleguide/go/decisions#type-aliases
  - [underscores](#underscores) ... based on https://google.github.io/styleguide/go/decisions#underscores
  - [useany](#useany) ... based on https://google.github.io/styleguide/go/decisions#use-any
  - [useq](#useq) ... based on https://google.github.io/styleguide/go/decisions#use-q
  - [varnames](#varnames) ... based on https://google.github.io/styleguide/go/decisions#variable-names

> ["Google Style Guides"](https://google.github.io/styleguide/) by Google is licensed under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/)

### [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) in Go wiki

- [contexts](#contexts) ... based on https://go.dev/wiki/CodeReviewComments#contexts
- [dontpanic](#dontpanic) ... based on https://go.dev/wiki/CodeReviewComments#dont-panic
- [errorstrings](#errorstrings) ... based on https://go.dev/wiki/CodeReviewComments#error-strings
- [handlerrors](#handlerrors) ... based on https://go.dev/wiki/CodeReviewComments#handle-errors

## Disabling and Ignoring

### Disable analyzer ( vet tool only )

Use `-[analyser name].disable` flag.

```console
$ go vet -vettool=`which gostyle` -mixedcaps.disable # Disable mixedcaps analyzer only
```

### Ignore directive

- `//lint:ignore`
- `//nolint:all`
- `//nostyle:all`
- `//nostyle:[analyzer name]` (e.g. `//nostyle:mixedcaps`)

## Configuration

`gostyle` can be configured like [golangci-lint](https://golangci-lint.run/usage/configuration/).

``` console
$ gostyle init
.gostyle.yml is generated
# As a Standalone CLI
$ gostyle run --config=.gostyle.yml ./...
# As a vet tool
$ go vet -vettool=`which gostyle` -gostyle.config=$PWD/.gostyle.yml ./...
```

> [!NOTE]
> If no configuration file is specified, `gostyle` will automatically search for `.gostyle.yml` or `.gostyle.yaml` in the Git root directory.

```yaml
# .gostyle.yml
analyzers:
  disable:
    # Disable specific analyzers.
    - analyzer-name
# All available settings of specific analyzers.
analyzers-settings:
  # See the dedicated "analyzers-settings" documentation section.
  option: value
```

### `analyzers-settings:`

#### contexts

```yaml
analyzers-settings:
  contexts:
    include-generated: false # include generated codes (default: false)
    exclude-test: true       # exclude test files (default: false)
```

#### dontpanic

```yaml
analyzers-settings:
  dontpanic:
    include-generated: false # include generated codes (default: false)
    exclude-test: true       # exclude test files (default: false)
```

#### errorstrings

```yaml
analyzers-settings:
  errorstrings:
    include-generated: false # include generated codes (default: false)
    exclude-test: true       # exclude test files (default: false)
```

#### funcfmt

```yaml
analyzers-settings:
  funcfmt:
    include-generated: false # include generated codes (default: false)
```

#### getters

```yaml
analyzers-settings:
  getters:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - GetViaHTTP
```

#### handlerrors

( **NOT** handl**ee**rrors )

```yaml
analyzers-settings:
  handlerrors:
    include-generated: false # include generated codes (default: false)
    exclude-test: true       # exclude test files (default: false)
```

#### ifacenames

```yaml
analyzers-settings:
  ifacenames:
    include-generated: false # include generated codes (default: false)
    all: true                # all interface names with the -er suffix are required (default: false)
```

#### mixedcaps

```yaml
analyzers-settings:
  mixedcaps:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - DBTX
      - EXPECT
```

#### nilslices

```yaml
analyzers-settings:
  nilslices:
    include-generated: false # include generated codes (default: false)
```

#### pkgnames

```yaml
analyzers-settings:
  pkgnames:
    include-generated: false # include generated codes (default: false)
```

#### recvnames

```yaml
analyzers-settings:
  recvnames:
    include-generated: false # include generated codes (default: false)
    max: 3                   # max length of receiver name (default: 2)
```

#### recvtype

```yaml
analyzers-settings:
  recvnames:
    include-generated: false # include generated codes (default: false)
```

#### repetition

```yaml
analyzers-settings:
  repetition:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - limitStr
```

#### typealiases

```yaml
analyzers-settings:
  typealiases:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - TmpAliasHeader
```

#### underscores

```yaml
analyzers-settings:
  underscores:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - DBTX
      - EXPECT
```

#### useany

```yaml
analyzers-settings:
  useany:
    include-generated: false # include generated codes (default: false)
```

#### useq

```yaml
analyzers-settings:
  useq:
    include-generated: false # include generated codes (default: false)
```

#### varnames

```yaml
analyzers-settings:
  varnames:
    include-generated: false  # include generated codes (default: false)
    small-scope-max: 5        # max lines for small scope (default: 7)
    small-varname-max: 3      # max length of variable name for small scope (default: -1)
    medium-scope-max: 10      # max lines for medium scope (default: 15)
    medium-varname-max: 5     # max length of variable name for medium scope (default: -1)
    large-scope-max: 15       # max lines for large scope (default: 25)
    large-varname-max: 7      # max length of variable name for large scope (default: -1)
    very-large-varname-max: 9 # max length of variable name for very large scope (default: -1)
    exclude:                  # exclude words
      - hostname
```

## Install

**go install:**

```console
$ go install github.com/k1LoW/gostyle@latest
```

**deb:**

``` console
$ export GOSTYLE_VERSION=X.X.X
$ curl -o gostyle.deb -L https://github.com/k1LoW/gostyle/releases/download/v$GOSTYLE_VERSION/gostyle_$GOSTYLE_VERSION-1_amd64.deb
$ dpkg -i gostyle.deb
```

**RPM:**

``` console
$ export GOSTYLE_VERSION=X.X.X
$ yum install https://github.com/k1LoW/gostyle/releases/download/v$GOSTYLE_VERSION/gostyle_$GOSTYLE_VERSION-1_amd64.rpm
```

**apk:**

``` console
$ export GOSTYLE_VERSION=X.X.X
$ curl -o gostyle.apk -L https://github.com/k1LoW/gostyle/releases/download/v$GOSTYLE_VERSION/gostyle_$GOSTYLE_VERSION-1_amd64.apk
$ apk add gostyle.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/gostyle
```

**[aqua](https://aquaproj.github.io/):**

```console
$ aqua g -i k1LoW/gostyle
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/gostyle/releases)
