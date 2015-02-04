package sts

import "errors"

// Callbacks should retuen this error, whatever the real reason is.
var ErrAuthFailed = errors.New("Auth failed")

func RegisterPasswordAuth() {

}

func RegisterPublicKeyAuth() {

}
