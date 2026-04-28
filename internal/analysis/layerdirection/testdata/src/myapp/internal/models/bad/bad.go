package bad

import "myapp/internal/repository" // want `models cannot import repository, service, or api`

var _ = repository.Repo{}
