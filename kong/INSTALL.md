# Kong

## HELMチャートからのインストール

Helm V3 でのインストール

```bash
helm repo add kong https://charts.konghq.com
helm repo update
kubectl create namespace kong
helm install serverless-kong kong/kong -n kong --values ./kong/values.yaml
```

## values.yaml

導入するKong GatewayはHelmのvalues.yamlで下記変更を加えている。

環境変数にpluginの設定を追加

```yaml
env:
  # add plugin configuration
  plugins: 'bundled, serverless'
  pluginserver_names: 'serverless'
  pluginserver_serverless_query_cmd: '/usr/local/bin/serverless -dump'
```

コンテナをカスタマイズ版に変更

```yaml
image:
  repository: 'public.ecr.aws/nautible/nautible-kong-serverless'
  tag: 'v0.1.5'
```

Daprの有効化

```yaml
podAnnotations:
  # add dapr configuration
  dapr.io/enabled: 'true'
  dapr.io/app-id: 'serverless'
  dapr.io/app-port: '8000'
  dapr.io/enable-api-logging: 'true'
```

## 更新

イメージ更新時のMinikubeへの反映

```bash
helm upgrade serverless-kong kong/kong -n kong --values ./kong/values.yaml
```
