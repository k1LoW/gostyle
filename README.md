# gostyle

**THIS IS A WORK IN PROGRESS AND PROOF OF CONCEPT AND STUDY WORK**

## Usage

```console
$ go vet -vettool=`which gostyle`
```

## Analyzers

### [Effective Go](https://go.dev/doc/effective_go)

- [ifacenames](analyzer/effective/ifacenames) ... https://go.dev/doc/effective_go#interface-names

### [Go Style](https://google.github.io/styleguide/go/)

- [**Style Guide**](https://google.github.io/styleguide/go/guide)
  - [mixedcaps](analyzer/guide/mixedcaps) ... https://google.github.io/styleguide/go/guide#mixed-caps
- [**Go Style Decisions**](https://google.github.io/styleguide/go/decisions)
  - [pkgnames](analyzer/decisions/pkgnames) ... https://google.github.io/styleguide/go/decisions#package-names

## Ignore Directive

- `//lint:ignore`
- `//nolint:all`
- `//nostyle:all`
- `//nostyle:[analyzer name]` (e.g. `//nostyle:mixedcaps`)

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
