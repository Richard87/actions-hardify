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
actions-hardify                                                 
Found 1 workflow file(s)

+------------------------------------------------+-------------+----------+--------+
| LOCATION                                       | PERMISSIONS | OLD      | NEW    |
+------------------------------------------------+-------------+----------+--------+
| pr.yaml > build                                | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > lint                                 | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > lint                                 | ok          | # v6.4.0 | v6.4.0 |
| pr.yaml > lint > golangci-lint                 | ok          | # v9.2.0 | v9.2.0 |
| pr.yaml > test                                 | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > test                                 | ok          | # v6.4.0 | v6.4.0 |
| pr.yaml > verify-code-generation               | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > verify-code-generation               | ok          | # v6.4.0 | v6.4.0 |
| pr.yaml > report-swagger-changes               | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > report-swagger-changes               | ok          | # v6.4.0 | v6.4.0 |
| pr.yaml > report-swagger-changes > Add comment | ok          | # v8     | v8     |
| pr.yaml > validate-radixconfig                 | ok          | # v6.0.2 | v6.0.2 |
| pr.yaml > validate-radixconfig                 | ok          | # v2.0.2 | v2.0.2 |
+------------------------------------------------+-------------+----------+--------+

Total: 13 finding(s)

✅ Workflows hardened successfully.
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

Thanks to [Step Security](https://github.com/step-security) for the inspiration behind this CLI. 