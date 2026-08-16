# 用户注册与用户统计问题

## Bug 是什么
用户名规范化、用户名唯一性检查、用户列表排序和全局用户统计之间存在组合性错误，
导致空白用户名被接受、重复用户名可注册、用户列表按注册时间升序返回、全局用户数统计错误。

## 如何触发
启动服务后：

```bash
curl -X POST 'http://127.0.0.1:8080/api/users' -d '{"username":"   "}'
curl -X POST 'http://127.0.0.1:8080/api/users' -d '{"username":"alice"}'
curl -X POST 'http://127.0.0.1:8080/api/users' -d '{"username":"alice"}'
curl 'http://127.0.0.1:8080/api/users'
curl 'http://127.0.0.1:8080/api/stats'
```

## 错误信息
空白用户名被接受；重复用户名未返回冲突；用户列表顺序与注册时间倒序约定不符；全局统计
中的用户数与实际不符。对应测试 `TestUserRegistrationAndStats` 失败。
