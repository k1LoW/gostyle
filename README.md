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

- **Style Guide**
  - [mixedcaps](analyzer/guide/mixedcaps) ... https://google.github.io/styleguide/go/guide#mixed-caps


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

