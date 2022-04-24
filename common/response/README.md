## 统一返回的JSON数据格式

统一格式
```golang
{
  "ret": 数字,        //响应码, 0 成功，非0失败
  "msg": 字符串,       //错误信息，"success" 成功，"xxxx" 错误信息提示
  "data": HashMap     //返回数据，放在键值对中
}
```

### 1. 列表数据
```golang
{
  "ret": 0,
  "msg": "success",
  "data": {
    "items": [
      {
        "id": "1",
        "name": "刘德华",
        "intro": "毕业于师范大学数学系，热爱教育事业，执教数学思维6年有余"
      }
    ]
  }
}
```

### 2.分页数据
```golang
{
  "ret": 0,
  "msg": "success",
  "data": {
    "total": 17,
    "rows": [
      {
        "id": "1",
        "name": "刘德华",
        "intro": "毕业于师范大学数学系，热爱教育事业，执教数学思维6年有余"
      }
    ]
}
```
### 3.无返回数据

```golang
{
  "ret": 0,
  "msg": "success",
  "data": null
}
```

### 4.请求失败
```golang
{
  "ret": 20001,           //错误码
  "msg": "错误具体信息",
  "data": null
}
```