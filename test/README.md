## http package

### 指定したパスにアクセスしたときにJSONを返すことを確認するテスト
- [x] /api/posts/{id}にアクセスしたときにJSONを返すことを確認すること
- [x] /api/posts/{id}にアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [x] /api/posts/createにアクセスしたときにJSONを返すことを確認すること
- [x] /api/posts/createにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [x] /api/posts/{id}/updateにアクセスしたときにJSONを返すことを確認すること
- [x] /api/posts/{id}/updateにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [x] /api/posts/{id}/deleteにアクセスしたときにJSONを返すことを確認すること
- [x] /api/posts/{id}/deleteにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [ ] /api/users/{id}にアクセスしたときにJSONを返すことを確認する
- [ ] /api/users/{id}にアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [x] /api/users/{id}/updateにアクセスしたときにJSONを返すことを確認すること
- [x] /api/users/{id}/updateにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する
- [x] /api/users/{id}/deleteにアクセスしたときにJSONを返すことを確認すること
- [x] /api/users/{id}/deleteにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する

### UserControllerのcreateまたはloginにアクセスしたときに正しいJSONを返すことを確認する
- [x] /api/users/createにアクセスしたときにJSONを返すことを確認すること
- [x] /api/users/loginにアクセスしたときにJSONを返すことを確認すること


## service package

### データベースからIDに一致しているデータを取り出すことを確認する
- [x] IDが１のデータを取得できることを確認すること

### データベースに保存されている投稿データを取得出来ることを確認する
- [x] 全ての投稿データを取得できていることを確認すること
- [x] データが空のときにエラーメッセージのJSONを返すことを確認すること

### データベースにデータを保存できることを確認する
- [x] データベースに指定した形式でデータを登録できることを確認すること

### データを編集できることを確認する
- [x] Post.Nameに「Google」、Post.URLに「https://www.google.com」と変更を加えたときに、データが更新されることを確認する

### データを削除できることを確認する
- [x] PostモデルのIDが1のデータを削除できることを確認すること

### 選択したカラムのデータのみを取得できることを確認する
- [x] PostモデルのNameカラムのみを取得できることを確認する
- [x] PostモデルのURLカラムのみを取得できることを確認する



## auth package

### jwtトークンをミドルウェアーで認証することを確認すること
- `/api/users/login, /api/users/create`でリクエストがきた時はJWTがなければ認証をかけないことを確認すること
- `/api/users/login, /api/users/create`でリクエストがきたときにJWTがあればエラーをだすことを確認すること
- headerが正しい形式ではないときにエラーを出すことを確認すること
- jwtをトークンが一致しないときはエラーを出すことを確認すること
- jwt認証が通ったときは指定のハンドラーにデータを渡すことを確認すること