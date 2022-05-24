# nautible-kong-serverless

KEDAと組み合わせることで、HTTPリクエストベースのサービスをサーバレス化するプラグイン

## 概要

Podを利用しないときはPod数0とし、HTTPリクエストをトリガーにPodを起動するためのプラグイン。KEDAと組み合わせて利用する。

## 仕様


## 構成

TODO 図

- APIGateway
  - Kong 2.8
- Queue
  - RabbitMQ（Minikube）
  - SQS(AWS)
  - ServiceBus(Azure)
- Replica Control
  - KEDA

## フォルダ構成

```text
$HOME
├ keda
│  ├ deploy.yaml ・・・ KEDAのデプロイ用マニフェストファイル
│  └ scaledobject.yaml ・・・ サンプルScaledObject（RabbitMQ接続）
├ kong
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

### キューの準備

RabbitMQにログインして、下記設定を追加する

#### Queues

- Type: classic
- Name: serverless
- Durability: Durable
- Auto delete: No
- arguments
  - x-message-ttl: 30000

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

## kong

## kongのマニフェスト取得（DB Less）

```bash
curl https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/single/all-in-one-dbless.yaml
```

取得したYAMLを種別（Kind）ごとにファイルを分け、kongディレクトリ以下に配置（カスタムリソースはkong/crds配下）

kong/deployment.yamlのname: proxyコンテナの環境変数に下記を追加する

```yaml
        - name: KONG_PLUGINS
          value: bundled, serverless
        - name: KONG_PLUGINSERVER_NAMES
          value: serverless
        - name: KONG_PLUGINSERVER_SERVERLESS_QUERY_CMD
          value: /usr/local/bin/serverless -dump
```

## コード

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

FROM kong:2.8

COPY --from=builder /go-plugins/bin/serverless /usr/local/bin/serverless
```

## build

プラグインをローカルでビルド

```bash
cd plugin
docker build -t nautible-kong-serverless:v0.1.0 -f ./package/Dockerfile .
```

## push

プラグインをECRのパブリックリポジトリにPush

※ タグ名は

```bash
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/nautible-kong-serverless
docker tag nautible-kong-serverless:v0.1.0 public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
docker push public.ecr.aws/nautible/nautible-kong-serverless:v0.1.0
```

## ローカル（Minikubeでの実行）

### kongデプロイ

### RabbitMQデプロイ

### KEDAデプロイ

### サンプルアプリケーションデプロイ

Minikubeにサンプルアプリケーションをデプロイする。

```bash
eval $(minikube docker-env)
cd sample_consumer
docker build -t consumer:v0.1.0 -f ./package/Dockerfile .
cd manifest
kubectl apply -f .
```

### ExternalIPの設定

```bash
minikube tunnel
```

### ブラウザからアクセス

```bash
http://localhost/kong/consumer
```
