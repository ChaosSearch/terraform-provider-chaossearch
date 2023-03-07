Tag the release (e.g. git tag v0.0.11)
Push the tags to GitHub (git push --tags)
Make sure you have the GPG key imported using gpg --import <key file>:

```
GPG_TTY=$(tty) GPG_FINGERPRINT='5420953B46E6E9EC2B5D9A77C57E6C31D20B3D96' GITHUB_TOKEN='<INSERT-YOUR-GITHUB-PERSONAL-ACCESS-TOKEN-HERE>' goreleaser release --rm-dist
```
