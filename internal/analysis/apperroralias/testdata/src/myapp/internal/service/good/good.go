package good

import apperrors "myapp/internal/errors"

func Foo() error { return apperrors.New("ok") }
