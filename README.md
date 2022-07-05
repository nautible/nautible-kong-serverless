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
│  ├ INSTALL.md ・・・ KEDAインストール説明ドキュメント
│  ├ scaledobject_aws.yaml ・・・ AWS用ScaledObjectデプロイ用マニフェストファイル
│  └ scaledobject_local.yaml ・・・ ローカル用ScaledObjectデプロイ用マニフェストファイル（RabbitMQ接続）
├ kong
│  ├ INSTALL.md ・・・ Kong Gatewayインストール説明ドキュメント
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
└ README.md
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

[plugin/README.md](./plugin/README.md)を参照

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

plugin/manifests/pubsub.yamlにRabbitMQの接続情報を記載しているので、パスワード部分のみ現在のRabbitMQのパスワードに変更して下記を実行する。

```bash
kubectl apply -f plugin/manifests/.
```

## デプロイしているコンテナを更新

```bash
helm upgrade serverless-kong kong/kong -n kong --values ./kong/values.yaml
```

### ExternalIPの設定

下記コマンド実行（sudoパスワードを聞かれた際は入力する）

```bash
minikube tunnel
```

### ブラウザからアクセス

```bash
http://localhost/kong/
```
