package bad

import "myapp/internal/repository" // want `api cannot import repository directly`

var _ = repository.Repo{}
