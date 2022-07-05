# nautible-kong-serverless

KEDAと組み合わせることで、HTTPリクエストベースのサービスをサーバレス化するプラグイン

## 概要

Podを利用しないときはPod数0とし、HTTPリクエストをトリガーにPodを起動するためのプラグイン。KEDAと組み合わせて利用する。

## 構成

- APIGateway
  - Kong 2.8
- Queue
  - RabbitMQ（Minikube）
  - SQS(AWS)
  - ServiceBus(Azure ※今後対応予定)
  - Pub/Sub(GoogleCloud ※今後対応予定)
- Pod AutoScaler
  - KEDA

## フォルダ構成

```text
$HOME
├ keda
│  ├ INSTALL.md ・・・ インストール説明ドキュメント
│  ├ scaledobject_aws.yaml ・・・ AWS用ScaledObjectデプロイ用マニフェストファイル
│  └ scaledobject_local.yaml ・・・ ローカル用ScaledObjectデプロイ用マニフェストファイル（RabbitMQ接続）
├ kong
│  ├ INSTALL.md ・・・ インストール説明ドキュメント
│  └ values.yaml ・・・ HELM設定ファイル
├ plugin
│  ├ cmd ・・・ プラグインのエントリーポイント
│  ├ manifests ・・・ サンプルマニフェストファイル
│  ├ package ・・・ Dockerfile
│  ├ pkg ・・・ カスタマイズ処理の本体
│  ├ go.mod ・・・ 導入モジュール
│  └ go.sum・・・ 依存モジュールのパスやバージョン
├ sample
│  ├ configmap.yaml ・・・ マウントするHTMLファイルを定義したConfigmap
│  ├ deployment.yaml ・・・ サンプルアプリケーション（Nginx）のDeployment
│  ├ service.yaml ・・・ サンプルアプリケーションのService
│  └ README.md ・・・ サンプルアプリケーションの説明
├ LICENSE
├ README.md
└ skaffold.yaml
```

## RabbitMQの準備（Minikube）

インストール

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install keda-queue bitnami/rabbitmq
```

パスワード確認

```bash
echo "Password      : $(kubectl get secret --namespace default keda-queue-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)"
```

### ログイン

ポートフォワードでログイン画面を表示する

kubectl port-forward -n default svc/keda-queue-rabbitmq 15672:15672

```text
user: user
password: 上記で確認したパスワード
```

### キューの準備

RabbitMQにログインして、下記設定を追加する

#### Queues

- Type: classic
- Name: serverless
- Durability: Durable
- Auto delete: No
- arguments
  - x-message-ttl: 30000 (Number)

#### Exchanges

- Name: serverless
- Type: funout
- Durability: Durable
- Auto delete: No
- Intenal: No

作成後、serverlessキューをBindする

## リポジトリ作成

ECRパブリックリポジトリにプラグイン用のリポジトリを作成する。（Terraformによる作成を推奨）

```text
nautible-kong-serverless
```

## KEDAの導入

[keda/INSTALL.md](./keda/INSTALL.md)を参照


## Kongの導入

[kong/INSTALL.md](./kong/INSTALL.md)を参照

## プラグイン作成

pluginディレクトリ配下にプラグイン本体のコードを作成する。

plugin
├ cmd ・・・ エントリーポイント
├ manifests ・・・ マニフェストファイル
├ package ・・・ パッケージング用コード（Dockerfile）
└ pkg ・・・ コード本体

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

## build

プラグインをローカルでビルド(バージョンは都度変更)

```bash
cd plugin
docker build -t nautible-kong-serverless:v0.1.0 -f ./package/Dockerfile .
```

## push

プラグインをECRのパブリックリポジトリにPush

```bash
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/nautible-kong-serverless
docker tag nautible-kong-serverless:v0.1.0 public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
docker push public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
```

※ Pushする前にレジストリに対して認証おく必要あり [参考](https://docs.aws.amazon.com/ja_jp/AmazonECR/latest/userguide/getting-started-cli.html)

## ローカル（Minikube）での実行

### サンプルアプリケーションデプロイ

```bash
kubectl apply -f sample/.
```

### KEDAのScaledObjectを導入

scaledobject_local.yamlのhostにRabbitMQの接続情報を記載しているので、パスワード部分のみ現在のRabbitMQのパスワードに変更して下記を実行する。

```bash
kubectl apply -f keda/scaledobject_local.yaml
```

### Kong Plugin設定を導入

pubsub.yamlにRabbitMQの接続情報を記載しているので、パスワード部分のみ現在のRabbitMQのパスワードに変更して下記を実行する。

```bash
kubectl apply -f plugin/manifests/.
```

## skaffold

```bash
skaffold dev
```


### ExternalIPの設定

下記コマンド実行後、

```bash
minikube tunnel
```

### ブラウザからアクセス

```bash
http://localhost/kong/consumer
```
