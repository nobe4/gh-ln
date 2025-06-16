<div align="center">
  <img width="300" src="https://github.com/nobe4/gh-ln/blob/main/docs/logo.png" /> <br>
  <sub>Logo by <a href="https://www.instagram.com/malohff">@malohff</a></sub>
</div>

# `gh-ln`

> Link files between repositories.

> [!IMPORTANT]
> This project is under heavy development.

This action creates a _link_ between files in various places. When the source is
updated, the destination is as well.

It works by using the GitHub API to read files and create Pull Requests where an
update is needed. You can specify the source, destination, and schedule for the
synchronization.

> [!TIP]
> The authentication for this can be rather tricky, make sure you read
> [authentication](/docs/authentication.md) to get familiar with the various
> methods.

## Quickstart

1. Install `gh-ln`

  ```
  gh extension install nobe4/gh-ln
  ```

1. Create a config file in `.ln-config.yaml`.

    E.g. [`ln-config.yaml`](.ln-config.yaml)

1. Run

  TODO

To use in Actions, see [nobe4/action-ln](https://github.com/nobe4/action-ln).

## Further readings

- [Authentication](/docs/authentication.md)
- [Configuration](/docs/configuration.md)
