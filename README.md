## 1. 開発時のリセット方法:
```shell
docker compose down -v
./cognito/init-cognito.sh
docker compose up -d
```

## 2. gitignoreでバージョン管理から除外:
- local_pool_id.jsonはgitに含まれなくなります
- テンプレートファイル（.template）のみをコミット

## 3. 自動初期化:
- 初回起動時やlocal_pool_id.jsonが存在しない場合、自動的にテンプレートから生成されます

これでdocker compose down -vの後も、ユーザーデータが残らずクリーンな状態から始められます。