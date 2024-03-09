# gengo
gengo 不仅仅是一个 CRUD 代码生成器，更是一个 Go 项目开发脚手架。通过 SQL 脚本直接生成规范统一的业务代码，能够显著提升后台 API 开发的效率。

(备注: go版本>=1.22.1)

### 1、编辑配置文件./configs/generate.yaml
```
# mysql配置
mysql:
  address: "127.0.0.1:3306"
  username: root
  password:
  db: canary
  max_open_conn: 20
  max_idle_conn: 10
  conn_lifetime: 120000
  charset: utf8
  debug: true
generate:
  # 项目名称
  package_name: suzuku
  # 服务名称
  service_name: user
  # 接口名称
  api_name: user_list
  # mysql数据表
  table_name: vip_users
  # sql语句
  sql: "SELECT uid,lark_id,udid,firstname,lastname,gender FROM vip_users WHERE gender=1 LIMIT 10 OFFSET 10;"
  # 操作 1:插入 2:更新 3:删除(软删除) 4:查询
  action: 4
```

### 2、执行./xgen-mac [或] .\xgen-win.exe [或] ./xgen-linux命令
执行命令前代码文件结构

![Snip20240105_3.png](lark%2Fassets%2Fimages%2FSnip20240105_3.png)

### 3、生成项目代码

![Snip20240105_4.png](lark%2Fassets%2Fimages%2FSnip20240105_4.png)

### 4、设置redis配置
api_suzuku.yaml
```
redis:
  address: ["127.0.0.1:6379"]
  db: 0
  password:
  prefix: "LK:"
  single: true
```

### 5、生成测试token
main_suzuku.go
```
func init() {
	……
	// 生成测试token
	token, _ := xjwt.CreateToken(1, 1, true, 3600*24*365)
	fmt.Println(token.Token)
	xredis.Set("USER:ACCESS_TOKEN_SESSION_ID:{1}:1", token.SessionId, 24*365*time.Hour)
}

```

### 6、调用创建的api
```
## user_list
curl "http://127.0.0.1:6600/api/user/user_list?uid=1&gender=1&page=1&limit=20" \
     -H 'Cookie: jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzU5Nzc1ODAsImlhdCI6MTcwNDQ0MTU4MCwiaXNzIjoibGFyay5jb20iLCJwbGF0Zm9ybSI6MSwic2Vzc2lvbl9pZCI6IjQ1Nzg2MGQ5MjljNjg5MzdmYjJmNWVmYjA5ZWIyNjAwIiwidWlkIjoiMSJ9.RExxbaus2wJF_mdYCBUbbrCSeNUH4VopsoTNB4ATXFw'
```
