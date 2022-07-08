# RabbitMQ

## HELMチャートからのインストール

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install keda-queue bitnami/rabbitmq
```

パスワード確認

```bash
echo "Password      : $(kubectl get secret --namespace default keda-queue-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)"
```

## ログイン

ポートフォワードでログイン画面を表示する

kubectl port-forward -n default svc/keda-queue-rabbitmq 15672:15672

```text
user: user
password: 上記で確認したパスワード
```

## キューの準備

RabbitMQにログインして、下記設定を追加する

### Queues

- Type: classic
- Name: serverless
- Durability: Durable
- Auto delete: No
- arguments
  - x-message-ttl: 30000 (Number)

### Exchanges

- Name: serverless
- Type: funout
- Durability: Durable
- Auto delete: No
- Intenal: No

作成後、serverlessキューをBindする
