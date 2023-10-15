# kosample

## ビルド方法

```bash
go build ./cmd/kosample
```

## コンテナビルド方法

```bash
ko build --base-import-paths ./cmd/kosample --local
```

## ローカルでのdocker-compose での起動

```bash
docker compose -f deploy/local/docker-compose.yml up
```
depends_onを指定しているが、アプリの起動の方が早い場合があり、HealthCheck部分を付け加えている。

## ローカルでのkubernetesでの起動

### kindでクラスタを作成する
```bash
kind create cluster
```

### dockerイメージのロード
```bash
kind load docker-image ko.local/kosample/cmd/kosample:latest
```
Kindの公式に`imagePullPolicy`を`IfNotPresent` or `Never`にしろと書いてある。  
参考:https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster

### マニフェストのapply
```bash
kubectl apply -f db-deploy.yaml
kubectl apply -f db-svc.yaml
kubectl apply -f app-deploy.yaml
kubectl apply -f app-svc.yaml
```

### 疎通

```bash
kubectl port-forward services/app 8080:8080
```

## コンテナイメージの扱い
### Github Actionの設定

リポジトリのSettings>Actions>Generalから、WorkflowのパーミッションをRead and Writeに変更する。

### ローカルへのpull

1. Settings>Developer Settingsにアクセス
2. Personal Access Tokens>Token(classic)にアクセス
3. 右上のGenerate new Tokenからトークンを作成
4. read:packagesをチェックしてpackageからダウンロードできるようにしておく
5. 下記コマンドを実行(tokenファイルにトークン情報が入っている)
	```bash
	cat token | docker login ghcr.io -u lottotto --password-stdin
	```
6. Docker pullする


### kind上でプライベートリポジトリからpullできるようにする。

1. imagepullsecretの作成
```bash
kubectl create -n it1 secret docker-registry regcred \
--docker-server=https://ghcr.io \
--docker-username=lottotto \
--docker-password=$GITHUB_TOKEN
```
2. マニフェスト変更
```diff
   apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
  namespace: it1
spec:
  selector:
    matchLabels:
      app: app
  template:
    metadata:
      labels:
        app: app
    spec:
      containers:
	  ==
+      imagePullSecrets:
+      - name: regcred
``` 
