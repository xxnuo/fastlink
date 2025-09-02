# fastlink

Quickly move files/directories to a certain location and create soft links at the original location

## Usage

### Help

```
fastlink [-m] [--move] [-k] [--keep] <source> [<destination>]
```

- `-m` or `--move` [ confilict with `-k` or `--keep` ] is the quick move mode flag, if not provided, will copy the original file to the destination location and then remove the original file. If provided, will skip copying and directly move the original file to the destination location. This is useful when you disk space is limited. Default is false. Can be set by `move` in config file.
- `-k` or `--keep` is the keep original file mode flag, if not provided, will remove the original file after copying. And not create soft link at the original location. This is useful when you only want to copy the original file. Default is false. Can be set by `keep` in config file.
- `source` is the file or directory to move
- `destination` is the location to move the file or directory to, if not provided, will read config file from `~/.config/fastlink/config.yaml`

### Examples

Move to `~/.config/fastlink/config.yaml#destination`:

```bash
fastlink /home/user/documents
```

Error: rescursive move is not allowed:

```bash
fastlink /home/user/documents /home/user/documents/backup
```

Keep original file and not to create soft link at the original location:

```bash
fastlink -k /home/user/documents
```

## Libraries

- Use [https://github.com/spf13/fileflow](https://github.com/spf13/fileflow) to copy/move files.
