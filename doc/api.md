# 接口文档

> GitHub 个人信息包括姓名、个人说明、邮箱、网站、头像，你需要在自己的后端中加入可以编辑个人信息的接口，包括查询和修改信息。研究 GitHub 是如何来完成这件事的，模仿 GitHub 的实现来编写自己的接口。在开始编码之前，书写一份接口文档，标注接口路径、HTTP Method、URL 参数、请求体和响应体数据类型、格式。

## 数据库：mongo

```
mongodb://localhost:27017
collection: echo-for-github
```

拿docker跑的

## 返回消息格式如下

```json
{
  "success": true,
  "message": "这是一些消息",
  "data": {
    "name": "ciro"
  }
}
```

## 关于User的若干操作

### 注册

```
api: /api/v1/user/register
method: POST
```

数据用json上传，应包括以下部分

```json
{
  "username": "ciro",
  "password": "123456879",
  "email": "1231231231@qq.com",
  "phone": "182123231523"
}
```

返回结果如下：

正确时状态码为`201`, 返回结果如下

```json
{
  "success": true,
  "message": "",
  "data": {
    "_idhex": "61a4b4c4aa209e9491c5d293"
  }
}
```

错误可能有：

```
username already exists
```

### 登录

### 查询个人信息(Profile)

emm,github好像是直接渲染在网页里了,没有看到对应的api

```
api:127.0.0.1:8080/api/v1/user/:username/profile
method: GET
```

200， 返回信息如下

> 数据类型: json

```json
{
  "success": true,
  "message": "",
  "data": {
    "name": "ciro",
    "bio": "happy!",
    "company": "bingyan",
    "location": "wuhan",
    "blog": "",
    "image": ""
  }
}
```

### 更改个人信息(Profile)

需要cookie`_gt_session`

```
api:/api/v1/users/:username/profile
method: POST
```

github的表单数据如下：

```
_method: put
authenticity_token: x8J70o5btnuv4346XlWspVAduGgL0ma50KVBg97ycyX+SKyWpTvmen0ceL/qYZ08I/CgOWRo1p40lsqU9AuXlw==
user[profile_name]: the nanme
user[profile_bio]: it's a test
user[profile_company]: BIngyan
user[profile_location]: wuhan
user[profile_blog]: www.bingyan.com
user[profile_twitter_username]: 
```

我期望的数据有：(可以以 ~~表单形式或~~ json传输)

```json
{
  "method": "put",
  "name": "ciro",
  "bio": "happy!",
  "company": "bingyan",
  "location": "wuhan",
  "blog": ""
}
```

响应为更改成功或者失败

以下为响应(json)

```json
{
  "success": true,
  "message": "",
  "data": "更改成功"
}
```

```json
{
  "success": false,
  "message": "请设置method为put",
  "data": null
}
```

### 更改头像

分为上传图片和设置头像两个阶段？

Image用编号表示

在数据库里开个image的编号和类型的存储

### 关于图片的上传

需要cookie`_gt_session`

```
api:127.0.0.1:8080/api/v1/user/cirolong/image
method:POST
```

采用原始表单上传

```html
<!doctype html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <title>Multiple file upload</title>
</head>

<body>
<h1>Upload multiple files with fields</h1>

<form action="/api/v1/user/cirolong/image" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="file"><br><br>
    <input type="submit" value="Submit">
</form>
</body>

</html>
```

返回数据格式如下：

```json
{
  "success": true,
  "message": "",
  "data": {
    "msg": "File 7pASiCOz2N.jpg uploaded successfully with fields name=%!s(MISSING) and email=%!s(MISSING)."
  }
}
```

### 获取头像图片

需要之前已经上传,不需要cookie，图片由username定位。

```
api:127.0.0.1:8080/api/v1/user/:username/image
method: GET
```

直接返回图片文件

## 关于Issues

> #### Level 2 - 实现 Issue 接口
>
>GitHub 还支持 Issue 的提交和变更，请你实现有关 Issue 的接口，要求同上。你需要支持的功能包括：
>
>- 获取 Issue 列表
>- 获取 Issue 的详细情况
>- 编辑 Issue 信息
>- 在 Issue 中发表回复
>- 关闭和开启 Issue

等等，github的issue是对reposity而言的，怎么绑定啊？

> 那就假设有一个不存在的Schrodinger仓库好了hh
>
> 直接绑定到user上面()

Model 定义

### how github do it

#### 新建

```
https://github.com/CiroLong/Go-Echo-System/issues
POST


issue[title]: test 3
issue[body]: how new a issue
```

#### get

```
请求网址: https://github.com/CiroLong/Go-Echo-System/issues/1
请求方法: GET
```

#### 评论

>
>
>1. 请求网址:
    > https://github.com/CiroLong/Go-Echo-System/issue_comments
>
>2.请求方法:
>
>
>
>   POST

```
issue: 1
comment[body]: 这是第二条回复

```

close

```
comment_and_close: 1
```

#### 修改title

```
请求网址: https://github.com/CiroLong/Go-Echo-System/issues/1
请求方法: POST
```

post表单

```
_method: put
authenticity_token: PzMbO9pRivbdLu6XWdkyVmkqQMgpIBlnmEewtWjB20MQ81GMwqKPODYxfkEcYJ2yN1omtOMOCJy7VD0OxZQfiA==
issue[title]: Test Title and 修改测试
```

### 自定义

#### GET

```
请求网址: /cirolong/issues/[编号]
请求方法: GET
```

返回标题， 创建者username

评论（body和auther

是否open

#### 创建

```
请求网址: /cirolong/issues
请求方法: POST
```

```
issue[title]: ???
issue[body]: how new a issue
```

返回创建成功，以及编号

#### 评论

```
请求网址: /cirolong/issue_comments
请求方法: POST
```

```
issue: 1   //这里标识第几个issue
comment[body]: 这是第二条回复
```

### api

#### 获取 Issue 列表

api: /cirolong/issues

method： GET

返回issue列表， 每一条包括

+ 编号

+ 发起人id

+ 标题

+ 状态

+ 评论数

#### 获取 Issue 的详细情况

参考GET issue

#### 编辑 Issue 信息

to be done

#### 在 Issue 中发表回复

参考issue的评论

#### 关闭和开启 Issue

