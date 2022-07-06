# プラグインの作成

## 構成

```text
plugin
├ cmd・・・ エントリーポイント
├ manifests ・・・ マニフェストファイル
├ package・・・ パッケージング用コード
└ pkg ・・・ コード本体
   ├ health_check・・・ ヘルスチェック処理
   └ pubsub・・・ Publish処理
```

### Dockerfile

プラグインをkongイメージに配置したカスタムイメージを作成するためのDockerfileを用意する。

```docker
FROM kong/go-plugin-tool:latest-alpine-latest as builder
ENV GO111MODULE=on

RUN mkdir /go-plugins
COPY plugin/go.mod /go-plugins/
COPY plugin/pkg/ /go-plugins/pkg/
RUN cd /go-plugins && \
    go build -o /go-plugins/bin/serverless pkg/main.go

FROM kong:2.8.1

COPY --from=builder /go-plugins/bin/serverless /usr/local/bin/serverless
```

なお、kongのベースイメージに更新が入った際は、本Dockerfileの FROM:kongのバージョンを更新する。

## build

プラグインをローカルでビルド(バージョンは都度変更)

```bash
cd plugin
docker build -t nautible-kong-serverless:<タグ> -f ./package/Dockerfile .
```

## push

プラグインをECRのパブリックリポジトリにPush

```bash
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/nautible-kong-serverless
docker tag nautible-kong-serverless:v0.1.0 public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
docker push public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
```

※ Pushする前にレジストリに対して認証おく必要あり [参考](https://docs.aws.amazon.com/ja_jp/AmazonECR/latest/userguide/getting-started-cli.html)