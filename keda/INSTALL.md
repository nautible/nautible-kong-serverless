# INSTALL

## KEDA

```bash
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
kubectl create namespace keda
helm install keda kedacore/keda --namespace keda
```

## ScaledObject

### ローカル環境

```bash
kubectl apply -f keda/scaledobject_local.yaml
```

### AWS環境

```bash
kubectl apply -f keda/scaledobject_aws.yaml
```
