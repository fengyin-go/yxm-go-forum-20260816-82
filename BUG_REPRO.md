# 回帖计数与统计一致性问题

## Bug 是什么
回帖内容校验、按主题取回帖、删除回帖时的计数更新和全局统计之间存在组合性错误，
导致空白回帖被接受、回帖列表取错主题、删除后计数反而上涨、全局回帖数统计错误。

## 如何触发
启动服务后创建一个主题，然后：

```bash
curl -X POST 'http://127.0.0.1:8080/api/topics/{id}/replies' -d '{"content":"   "}'
curl 'http://127.0.0.1:8080/api/topics/{id}/replies'
curl 'http://127.0.0.1:8080/api/stats'
curl -X DELETE 'http://127.0.0.1:8080/api/replies/{reply_id}'
```

## 错误信息
空白内容回帖应被拒绝却创建成功；回帖列表返回了别的主题下的回帖；全局统计的回帖数与
主题详情不一致；删除回帖后回帖数反而变大。对应测试 `TestReplyCountAndStatsConsistency`
失败。
