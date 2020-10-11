# skk-cli

![](https://raw.githubusercontent.com/mattn/skk-cli/main/misc/screenshot.png)

## Usage

```
Usage of skk-cli:
  -V    Print the version
  -d value
        path to SKK-JISYO.L
  -json
        JSON mode
```

You can specify multiple `-d`.

```
$ skk-cli -d /path/to/SKK-JISYO.L -d /path/to/SKK-JISYO.emoji.utf8
```

## Installation

```
go get github.com/mattn/skk-cli
```

SKK-JISYO.L must be located at `~/.config/skk-cli/SKK-JISYO.L`. For Windows: `%APPDATA%\skk-cli\SKK-JISYO.L`

## License

MIT

## Author

Yasuhiro Matsumoto
