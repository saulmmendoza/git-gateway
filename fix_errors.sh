sed -i 's/handleError(unauthorizedError(err.Error()), w, r)/handleError(unauthorizedError("%s", err.Error()), w, r)/g' api/bitbucket.go
sed -i 's/handleError(unauthorizedError(err.Error()), w, r)/handleError(unauthorizedError("%s", err.Error()), w, r)/g' api/forgejo.go
sed -i 's/handleError(unauthorizedError(err.Error()), w, r)/handleError(unauthorizedError("%s", err.Error()), w, r)/g' api/github.go
sed -i 's/handleError(unauthorizedError(err.Error()), w, r)/handleError(unauthorizedError("%s", err.Error()), w, r)/g' api/gitlab.go
