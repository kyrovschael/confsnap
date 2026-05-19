# confsnap

Snapshot and diff config file states across deploys to track configuration drift.

---

## Installation

```bash
go install github.com/youruser/confsnap@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/confsnap.git && cd confsnap && go build ./...
```

---

## Usage

Take a snapshot of your config files before a deploy:

```bash
confsnap snapshot --paths /etc/nginx,/etc/app/config.yaml --label pre-deploy-v1.2
```

After the deploy, take another snapshot and diff the two:

```bash
confsnap snapshot --paths /etc/nginx,/etc/app/config.yaml --label post-deploy-v1.2
confsnap diff pre-deploy-v1.2 post-deploy-v1.2
```

Example output:

```
--- /etc/app/config.yaml (pre-deploy-v1.2)
+++ /etc/app/config.yaml (post-deploy-v1.2)
@@ -4,7 +4,7 @@
-  timeout: 30s
+  timeout: 60s
```

List all saved snapshots:

```bash
confsnap list
```

---

## Snapshots

Snapshots are stored locally in `~/.confsnap/` by default. Override with `--store /path/to/dir`.

---

## License

MIT © youruser