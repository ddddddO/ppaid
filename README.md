# 💊PuCo
Tool that aids **P**HP**U**nit and P**CO**V

Running tests takes a lot of time, and generating test coverage also takes a significant amount of time.
This tool might be useful in such situations.

![](./assets/puco_demo.gif)

> [!IMPORTANT] 
> `vendor/bin/phpunit` and `pcov` are required.

## With this tool,
1. Select test files to run,
1. Select files for which you want to generate coverage reports (HTML),
1. You might be able to execute steps 1 and 2 easily and quickly. Probably. Probably..

## Which test files can be selected?
Files located in `tests/` directory.

## Which files can be selected to get coverage?
Files located in `src/` or `app/` directory.

## How to install?

### Homebrew

```console
brew install ddddddO/tap/puco
```

### Go
```console
go install github.com/ddddddO/puco/cmd/puco@latest
```

## Usage

```console
$ puco --help
Usage: puco [options]
puco

Options:
  -repeat
        This flag starts with data selected by the most recently executed puco.

Example:
  puco          # normal launch
  puco --repeat # launch using the most recent data
$
```

## Processing of PuCo
todo: 英語もほしい

1. 選択されたテストファイルのパスを取得(複数可)
1. 選択されたカバレッジレポートを生成したいファイルのパスを取得(複数可)
1. 2の最長一致のパス(ディレクトリパス)を計算
    - ※このディレクトリパスの配下がカバレッジ生成の対象なので、2で選択された各ファイルパスのみがカバレッジ生成対象ではないことに注意
1. 1と3と既存の`phpunit.xml`があれば、それらを元に`phpunitxml_generated_by_puco.xml`を生成
1. 実行する`php`コマンドを組み立て、実行する
1. `coverage-puco`ディレクトリ配下にカバレッジレポートが生成される

> [!WARNING]
> ※`puco`初回実行時に、`~/.config/puco.toml`という設定ファイルができます。
> この設定ファイル内のキー:`CommandToSpecifyBeforePHPCommand`にdockerコマンド越しにphpコマンドを実行するよう記載していますが、直接phpコマンドを実行したい場合は、このキーの値を`""`にしていただくか、この行ごと消してください。

## TODO
- [ ] カバレッジレポートをHTML形式以外でも出力できるようにする
- [ ] ヒストリー機能欲しい
    - 何度もファイル選択は手間。ただ、ツールで組み立てられたコマンドは表示されるので、それコピペで実行でも代替できるから後でいいかも
    - 一旦repeatフラグを実装したから、ほんとに欲しくなってからでいいかも