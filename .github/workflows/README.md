# To Run the Workflows locally

It is useful to use [`act`](https://nektosact.com/installation/index.html) to run the workflows locally before pushing them to the repository.

> __Note__: `GITHUB_TOKEN` is required.

```bash
act -s GITHUB_TOKEN="$(gh auth token)"
```

## Running the Workflows via GitHub CLI

If you have installed `gh` then `act` GitHub CLI extension can be installed as follows:

```bash
# Install the extension
gh extension install https://github.com/nektos/gh-act

# Run the workflow
gh act -P ubuntu-latest=catthehacker/ubuntu:act-latest
```
