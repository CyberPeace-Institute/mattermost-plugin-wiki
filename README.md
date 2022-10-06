# Mattermost Wiki Plugin

This is a plugin trying to bring the ability to create documentation in mattermost.
The docs will be in .md format and saved in the mattermost database.

## Getting Started
Clone the repository:
```
git clone https://github.com/CyberPeace-Institute/mattermost-plugin-wiki wiki-plugin
```

Note that this project uses [Go modules](https://github.com/golang/go/wiki/Modules). Be sure to locate the project outside of `$GOPATH`.

### Build your plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/cpi.wiki-0.1.0.tar.gz
```

**or**

### Build the plugin is through the Dockerfile:
```
docker build -t wiki . && docker run --name mattermost-wiki wiki
docker cp mattermost-wiki:/app/dist/cpi.wiki-0.1.0.tar.gz .
```

Then no matter how you created the **cpi.wiki-0.1.0.tar.gz** file you can upload it and [manually install](https://mattermost.gitbook.io/plugin-confluence/setup/install-the-mattermost-plugin#alternative-install-via-manual-upload) it.

## Misc

This project was based on:
- https://github.com/mattermost/mattermost-plugin-playbooks/
- https://github.com/mattermost/mattermost-plugin-starter-template
