# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.4] - 2023-03-16

### Bug Fix

- Fix race condition when using multiple goroutines to inject the error data

## [0.0.3] - 2021-11-18

### Changed

- Change `ErrorDefinition.maskMessage` type to `*string` 
- Change `ErrorDefinition.formatter` type to `*MessageFormatter` 
- Change `ErrorDefinition.maskFormatter` type to `*MaskFormatter` 

## [0.0.2] - 2021-10-26

### Added

- Add error wrapper interface as replacement to use direct struct, thus reducing [nil is not nil](https://yourbasic.org/golang/gotcha-why-nil-error-not-equal-nil) problem
- Add context injection for storing error data
- Add formatter function for error message and mask message
- Add error categorization
- Add `Cast()` and `Convert()` util function
- Add error code string as an addition to current numerical error code. The difference is, numerical error code is an anonymized form of the error that can be sent to user, while error code string is a static string, usually same as variable name, and can be used by developer for metrics tagging and easier error finding in addition to stack traces
- Add `ErrorDefinition.Masked*()` functions, removing `NewMaskedError()` function usage
- Add stack trace trimmer via `errwrap.DefaultPackagePrefix` and `errwrap.DefaultStackTraceMode`

### Changed

- Change package name, from `errors` to `errwrap`
- Change error wrapper struct visibility to private, instead use the error wrapper interface
- Change `NewError()` parameters

### Removed

- Remove `ErrorWrapper.WithData()` function, error data is now injected via error wrapper instantiation
- Remove `WithMaskedError()` function, use `ErrorDefinition.Masked*()` functions to mark the error definition as masked error.
- Remove `HTTPErrorCode` field and replaced with error categorization to generalize usage, so the error wrapper is not bound to HTTP only


## [0.0.1] - 2021-09-14

### Added

- Initial version of ErrWrap

[Unreleased]: https://github.com/rapidashorg/errwrap/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/rapidashorg/errwrap/releases/tag/v0.0.1
[0.0.2]: https://github.com/rapidashorg/errwrap/releases/tag/v0.0.2
[0.0.3]: https://github.com/rapidashorg/errwrap/releases/tag/v0.0.3
