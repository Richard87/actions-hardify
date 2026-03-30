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

## TODO:

- BUG: `error: parsing radix-acr-cleanup/charts/radix-acr-cleanup/templates/deployment.yaml: yaml: line 5: did not find expected node content`
- Add License and contributor file

## Acknowledgements

Thanks to [Step Security](https://github.com/step-security) for the inspiration behind this CLI. 