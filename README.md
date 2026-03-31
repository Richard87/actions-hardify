# GitHub Actions Hardener

A CLI tool that finds all GitHub Actions workflows in a folder and hardens them.

## Installation

```sh
go install github.com/richard87/actions-hardify@latest
```

- Restrict permissions for GITHUB_TOKEN.
- Pin actions to a full length commit SHA.
- list outdated versions, suggest upgrade to newest version
- Use github api to find versions

## Usage

```
$ export GITHUB_TOKEN=$(gh auth token)
$ actions-hardify                                                 
                      
Found 2 workflow file(s)

+--------------------------------+----------------+-------------+----------+--------+
| FILE                           | LOCATION       | PERMISSIONS | OLD      | NEW    |
+--------------------------------+----------------+-------------+----------+--------+
| .github/workflows/pr.yaml      | build          | ok          | # v6.0.2 | v6.0.2 |
| .github/workflows/pr.yaml      | build          | ok          | # v6.4.0 | v6.4.0 |
| .github/workflows/pr.yaml      | test           | ok          | # v6.0.2 | v6.0.2 |
| .github/workflows/pr.yaml      | test           | ok          | # v6.4.0 | v6.4.0 |
| .github/workflows/pr.yaml      | lint           | ok          | # v6.0.2 | v6.0.2 |
| .github/workflows/pr.yaml      | lint           | ok          | # v6.4.0 | v6.4.0 |
| .github/workflows/pr.yaml      | lint           | ok          | # v9.2.0 | v9.2.0 |
| .github/workflows/release.yaml | release-please | ok          | # v4.4.0 | v4.4.0 |
+--------------------------------+----------------+-------------+----------+--------+

Total: 8 finding(s)

✅ Workflows hardened successfully.
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

Thanks to [Step Security](https://github.com/step-security) for the inspiration behind this CLI. 