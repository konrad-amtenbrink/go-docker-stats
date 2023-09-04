<!-- markdownlint-configure-file {
  "MD013": {
    "code_blocks": false,
    "tables": false
  },
  "MD033": false,
  "MD041": false
} -->
<div align="center">
<h1>go-docker-stats</h1>

Comparison of distroless vs. scratch for go docker builds.

> **FYI**: Feel free to open PRs as I am sure the Dockerfiles can be improved (@joelrose).

[Comparison](#comparison) •
[Results](#results)

</div>

## Comparison
We compare docker images sizes based on a super simple http server only having [httprouter](https://github.com/julienschmidt/httprouter) and [zerolog](https://github.com/rs/zerolog)
as dependencies.

## Results
| Base image   | Image size (in MB) |
|--------------|--------------------|
| `scratch`    |    7.01            |
| `distroless` |    9.25            |
