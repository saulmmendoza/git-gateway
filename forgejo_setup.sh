# Create forgejo.go from github.go
cp api/github.go api/forgejo.go
sed -i 's/GitHub/Forgejo/g' api/forgejo.go
sed -i 's/github/forgejo/g' api/forgejo.go
sed -i 's/gh /fj /g' api/forgejo.go
