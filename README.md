# code-strength

リポジトリ内ディレクトリごとのAI要求レベル定義を生成するGo/Cobra CLIです。

## 使い方

対話型で検索・複数選択する場合:

```sh
go run ./cmd/code-strength generate --root .
```

TTYでは、入力した文字で候補がリアルタイムに絞り込まれます。上下キーで移動し、Spaceで複数選択、Enterで確定します。パイプ入力などTTY以外では、検索語と候補番号を入力する簡易UIに切り替わります。

非対話で生成する場合:

```sh
go run ./cmd/code-strength generate --root . --production services/api --production docker/production
```

生成済みの要求水準を確認する場合:

```sh
go run ./cmd/code-strength level internal/scanner/scanner.go
# absolute paths are also supported
go run ./cmd/code-strength level "$(pwd)/internal/scanner/scanner.go"
```

ディレクトリまたはファイルパスを指定すると、最も近い親ディレクトリの `development` / `production` を1行で返します。
該当する定義がない場合は `unknown` を返します。

## Codex Plugin

AIエージェントが実装前に要求水準を確認できるSkillを、配布用Pluginとして `plugins/code-strength-ai` に含めています。PluginディレクトリをCodexのPlugin配置先へコピーすると利用できます。

追加の除外対象は `--exclude` を繰り返して指定できます。既定の出力先はリポジトリルートの `.ai-requirements.yml` です。既存ファイルは保持せず、毎回全量を書き換えます。

選択したディレクトリ以下は `production` として出力され、それ以外は `development` になります。AI向け起点ファイルの管理や依存関係の判断はCLIの対象外です。
