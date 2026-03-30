# Contributing to Actions Hardener

Thanks for your interest in contributing!

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -am 'Add my change'`)
6. Push to the branch (`git push origin feature/my-change`)
7. Open a Pull Request

## Development

```sh
# Build
go build ./...

# Run tests
go test ./...

# Run locally
go run main.go --dry-run -d /path/to/repo
```

## Reporting Issues

Please open a GitHub issue with:
- A description of the problem
- Steps to reproduce
- Expected vs actual behavior

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
