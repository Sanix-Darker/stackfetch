## STACKFETCH

> Heavily inspired by how neofetch works in principe.

Ever got an issue from someone you can't reproduce on your side ?
But then asking to run a bunch of command line to get like, the version of his cli/package-manager/system environment, processor... ?

Well now, you're going to ask to run only one thing and get everything at once.

Meet stackfetch, with a provided programming language or just a stack, you are going to have all informations you want to help debuguing.

![screenshot](./screenshot.png)

### DEV INSTALL

```bash
go install github.com/sanix-darker/stackfetch/cmd/stackfetch@latest  # source build
# — or —
wget https://github.com/sanix-darker/stackfetch/releases/download/vX.Y.Z/stackfetch-linux-amd64
chmod +x stackfetch-* && mv stackfetch-* /usr/local/bin/stackfetch
```

### USAGE

```bash
stackfetch                         # system only (fast)
stackfetch node python docker      # add Node, Python, Docker info
stackfetch mean lamp --json        # JSON report for MEAN & LAMP stacks
```

### CONTRIBUTION

By [sanixdk](https://github.com/sanix-darker).
