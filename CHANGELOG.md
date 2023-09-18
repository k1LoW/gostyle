# Changelog

## [v0.12.0](https://github.com/k1LoW/gostyle/compare/v0.11.0...v0.12.0) - 2023-09-18
### New Features 🎉
- Add typealiases analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/56
- Add useq analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/58
### Other Changes
- Add reference style credits by @k1LoW in https://github.com/k1LoW/gostyle/pull/57

## [v0.11.0](https://github.com/k1LoW/gostyle/compare/v0.10.6...v0.11.0) - 2023-09-16
### New Features 🎉
- varnames analyzer supports excluding variables from analysis. by @k1LoW in https://github.com/k1LoW/gostyle/pull/52
- Add recvtype analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/54

## [v0.10.6](https://github.com/k1LoW/gostyle/compare/v0.10.5...v0.10.6) - 2023-09-16
### Fix bug 🐛
- Fix repetition bug by @k1LoW in https://github.com/k1LoW/gostyle/pull/49
- Skip test function by @k1LoW in https://github.com/k1LoW/gostyle/pull/50

## [v0.10.5](https://github.com/k1LoW/gostyle/compare/v0.10.4...v0.10.5) - 2023-09-16
### New Features 🎉
- Separate each analyzer for stand-alone use by @k1LoW in https://github.com/k1LoW/gostyle/pull/47

## [v0.10.4](https://github.com/k1LoW/gostyle/compare/v0.10.3...v0.10.4) - 2023-09-16
### New Features 🎉
- Update nilslices ( Check distinction between a nil slice and a non-nil ) by @k1LoW in https://github.com/k1LoW/gostyle/pull/45

## [v0.10.3](https://github.com/k1LoW/gostyle/compare/v0.10.2...v0.10.3) - 2023-09-15

## [v0.10.2](https://github.com/k1LoW/gostyle/compare/v0.10.1...v0.10.2) - 2023-09-15
### Fix bug 🐛
- `Test*` function exists in `*_test` package by @k1LoW in https://github.com/k1LoW/gostyle/pull/42

## [v0.10.1](https://github.com/k1LoW/gostyle/compare/v0.10.0...v0.10.1) - 2023-09-14

## [v0.10.0](https://github.com/k1LoW/gostyle/compare/v0.9.0...v0.10.0) - 2023-09-14
### New Features 🎉
- Add getter analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/40

## [v0.9.0](https://github.com/k1LoW/gostyle/compare/v0.8.0...v0.9.0) - 2023-09-13
### New Features 🎉
- Add nilslices analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/38

## [v0.8.0](https://github.com/k1LoW/gostyle/compare/v0.7.1...v0.8.0) - 2023-09-13
### New Features 🎉
- Add repetition analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/35

## [v0.7.1](https://github.com/k1LoW/gostyle/compare/v0.7.0...v0.7.1) - 2023-09-13
### New Features 🎉
- Add .gostyle.yml by @k1LoW in https://github.com/k1LoW/gostyle/pull/34
### Fix bug 🐛
- Fix lookup func of types.Config by @k1LoW in https://github.com/k1LoW/gostyle/pull/32
- Fix typo by @k1LoW in https://github.com/k1LoW/gostyle/pull/33

## [v0.7.0](https://github.com/k1LoW/gostyle/compare/v0.6.0...v0.7.0) - 2023-09-13
### New Features 🎉
- Add varnames analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/27
### Fix bug 🐛
- Support *ast.RangeStmt (underscores, mixedcaps) by @k1LoW in https://github.com/k1LoW/gostyle/pull/28
- Revert "Add `init` command to generate .gostyle.yml" by @k1LoW in https://github.com/k1LoW/gostyle/pull/29
- Use source by @k1LoW in https://github.com/k1LoW/gostyle/pull/30

## [v0.6.0](https://github.com/k1LoW/gostyle/compare/v0.5.0...v0.6.0) - 2023-09-12
### Breaking Changes 🛠
- Support config file by @k1LoW in https://github.com/k1LoW/gostyle/pull/23
### New Features 🎉
- Add `init` command to generate .gostyle.yml by @k1LoW in https://github.com/k1LoW/gostyle/pull/25

## [v0.5.0](https://github.com/k1LoW/gostyle/compare/v0.4.0...v0.5.0) - 2023-09-12
### New Features 🎉
- Support exclude words in mixedcaps analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/20
- Check only where the name is defined. by @k1LoW in https://github.com/k1LoW/gostyle/pull/21
- Add underscores analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/22
### Other Changes
- Check renamed package name by @k1LoW in https://github.com/k1LoW/gostyle/pull/19

## [v0.4.0](https://github.com/k1LoW/gostyle/compare/v0.3.2...v0.4.0) - 2023-09-12
### Breaking Changes 🛠
- By default, generated code is not inspected. by @k1LoW in https://github.com/k1LoW/gostyle/pull/17

## [v0.3.2](https://github.com/k1LoW/gostyle/compare/v0.3.1...v0.3.2) - 2023-09-12
### New Features 🎉
- Reporter can be used generally. by @k1LoW in https://github.com/k1LoW/gostyle/pull/15

## [v0.3.1](https://github.com/k1LoW/gostyle/compare/v0.3.0...v0.3.1) - 2023-09-12
### New Features 🎉
- Add analyzer for receiver names by @k1LoW in https://github.com/k1LoW/gostyle/pull/13

## [v0.3.0](https://github.com/k1LoW/gostyle/compare/v0.2.0...v0.3.0) - 2023-09-11
### New Features 🎉
- Support `-disable` option by @k1LoW in https://github.com/k1LoW/gostyle/pull/8
- Add pkgnames analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/10
- Add message prefix by @k1LoW in https://github.com/k1LoW/gostyle/pull/11

## [v0.2.0](https://github.com/k1LoW/gostyle/compare/v0.1.0...v0.2.0) - 2023-09-11
### New Features 🎉
- Support inline ignore directives in comments by @k1LoW in https://github.com/k1LoW/gostyle/pull/7

## [v0.1.0](https://github.com/k1LoW/gostyle/commits/v0.1.0) - 2023-09-11
### New Features 🎉
- Add MixedCaps analyzer by @k1LoW in https://github.com/k1LoW/gostyle/pull/3
- Add `-all` option to ifacenames analyzer. by @k1LoW in https://github.com/k1LoW/gostyle/pull/4
### Other Changes
- Bump golang.org/x/text from 0.3.3 to 0.3.8 by @dependabot in https://github.com/k1LoW/gostyle/pull/1
- mv `mixedcaps` to analyzer/guide/ by @k1LoW in https://github.com/k1LoW/gostyle/pull/5
