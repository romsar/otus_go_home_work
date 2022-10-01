package mocks

//go:generate find . ! -name gen.go -type f -o -type d -maxdepth 1 -mindepth 1 -exec rm -rf {} +
//nolint:lll
//go:generate go run github.com/vektra/mockery/v2 --all --keeptree --case underscore --dir "../" --output "./" --outpkg "mocks"
