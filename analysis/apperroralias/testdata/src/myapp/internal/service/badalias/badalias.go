package badalias

import apierrors "myapp/internal/errors" // want `must use alias .apperrors., not .apierrors.`

func Foo() error { return apierrors.New("nope") }
