# Mattermost Wiki Plugin

This is a plugin trying to bring the ability to create documentation in mattermost.
The docs will be in .md format and saved in the mattermost database.

## Getting Started
Clone the repository:
```
git clone https://github.com/CyberPeace-Institute/mattermost-plugin-wiki wiki-plugin
```

Note that this project uses [Go modules](https://github.com/golang/go/wiki/Modules). Be sure to locate the project outside of `$GOPATH`.

Build your plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.example.my-plugin.tar.gz
```

## Misc

This project was based on:
- https://github.com/mattermost/mattermost-plugin-playbooks/
- https://github.com/mattermost/mattermost-plugin-starter-template
