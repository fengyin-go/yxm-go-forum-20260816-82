# 主题列表搜索与分页问题

## Bug 是什么
主题列表接口在关键词匹配、featured 过滤、翻页起点和分页大小上限四个环节都存在错误，
多个错误叠加后表现为搜索结果缺失、第二页内容错位、超大分页返回空列表。

## 如何触发
启动服务后依次请求：

```bash
curl 'http://127.0.0.1:8080/api/topics?keyword=hello&page=1&size=20&featured=false'
curl 'http://127.0.0.1:8080/api/topics?page=2&size=1'
curl 'http://127.0.0.1:8080/api/topics?page=1&size=200'
```

## 错误信息
第一个请求按 `hello` 搜索标题 `Hello World` 返回空；第二个请求第二页返回了错误主题；
第三个请求分页大小超过上限后返回空列表，而不是夹紧到上限。对应测试
`TestTopicListSearchAndPagination` 失败。
