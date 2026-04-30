package badbare

import "myapp/internal/errors" // want `must use alias .apperrors. \(bare import shadows stdlib errors\)`

func Foo() error { return errors.New("nope") }
