# gostyle

**THIS IS A WORK IN PROGRESS AND PROOF OF CONCEPT AND STUDY WORK**

`gostyle` is a set of analyzers for coding styles.

## Disclaimer

`gostyle` **IS NOT** [Go Style](https://google.github.io/styleguide/go/).

"Go Style" is a great style and we will actively refer to it, but we cannot implement the same rules exactly, and we may extend the rules.

## Usage

```console
$ go vet -vettool=`which gostyle`
```

## Analyzers

**Although not perfect**, it provides analyzers based on helpful styles.

### [Effective Go](https://go.dev/doc/effective_go)

- [ifacenames](analyzer/effective/ifacenames) ... based on https://go.dev/doc/effective_go#interface-names

### [Go Style](https://google.github.io/styleguide/go/)

- [**Guide**](https://google.github.io/styleguide/go/guide)
  - [mixedcaps](analyzer/guide/mixedcaps) ... based on https://google.github.io/styleguide/go/guide#mixed-caps
- [**Decisions**](https://google.github.io/styleguide/go/decisions)
  - [nilslices](analyzer/decisions/nilslices) ... based on https://google.github.io/styleguide/go/decisions#nil-slices
  - [pkgnames](analyzer/decisions/pkgnames) ... based on https://google.github.io/styleguide/go/decisions#package-names
  - [recvnames](analyzer/decisions/recvnames) ... based on https://google.github.io/styleguide/go/decisions#receiver-names
  - [repetition](analyzer/decisions/repetition) ... based on https://google.github.io/styleguide/go/decisions#repetition
  - [underscores](analyzer/decisions/underscores) ... based on https://google.github.io/styleguide/go/decisions#underscores
  - [varnames](analyzer/decisions/varnames) ... based on https://google.github.io/styleguide/go/decisions#variable-names

## Disabling and Ignoring

### Disable analyzer

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

`gostyle` can be configured by `-gostyle.config=$PWD/.gostyle.yml` like [golangci-lint](https://golangci-lint.run/usage/configuration/).

``` console
$ go vet -vettool=`which gostyle` -gostyle.config=$PWD/.gostyle.yml ./...
```

```yaml
# All available settings of specific analyzers.
analyzers-settings:
  # See the dedicated "analyzers-settings" documentation section.
  option: value
analyzers:
  disable:
    # Disable specific analyzers.
    - analyzer-name
```

### `analyzers-settings:`

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
```

#### repetition

```yaml
analyzers-settings:
  repetition:
    include-generated: false # include generated codes (default: false)
    exclude:                 # exclude words
      - limitStr
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

**manually:**

Download binary from [releases page](https://github.com/k1LoW/gostyle/releases)
