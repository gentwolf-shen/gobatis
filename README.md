# gobatis 

模仿MyBatis开发的一个go版本，实现了一些基本功能，请查看 sample.xml。

使用类库:
* XML解析库：github.com/beevik/etree
* 数据库：github.com/jmoiron/sqlx
* 日志库：go.uber.org/zap

## 1. 安装

```shell script
go get github.com/gentwolf-shen/gobatis
```

* 建议使用go mod

## 2. 初始化
```go
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gentwolf-shen/gobatis"
)
```

```go
// 数据库配置请查看 github.com/jmoiron/sqlx
config := gobatis.DbConfig{
    Driver:             "mysql",
    Dsn:                "username:password@tcp(mysqlhost:3306)/dbname?charset=utf8",
    MaxOpenConnections: 5,
    MaxIdleConnections: 1,
    MaxLifeTime:        3600,
    MaxIdleTime:        1800,
}

// gobatis对象
app := gobatis.NewGoBatis(config)

// 设置xml文件
filename := "/path/to/your/mapper.xml"
err := app.LoadFromFile("User", filename)
if err != nil {
    panic(err)
}
```

## 3. 添加

```go
lastInsertId, err := app.Insert("User.Insert", map[string]interface{}{
    "username":   "test-user",
    "password":   "pwd-111111",
    "status":     1,
    "createTime": time.Now().Unix(),
})
fmt.Println(err)
fmt.Println(lastInsertId)
```

## 4. 更新

```go
rowsAffected, err := app.Update("User.Update", map[string]interface{}{
    "id":       1,
    "username": "new-user",
})
fmt.Println(err)
fmt.Println(rowsAffected)
```

## 5. 查询一条记录

```go
// 用户结构
type (
	User struct {
		Id         int64
		Username   string
		Password   string
		Status     uint8
		CreateTime int64 `db:"create_time"`
	}
)
```

```go
var user User
err := app.QueryObject(&user, "User.Query", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(user)
```

## 6. 查询多条记录

```go
var users []User
err := app.QueryObjects(&users, "User.List", map[string]interface{}{"limit": 10})
fmt.Println(err)
fmt.Println(users)
```

## 7. 查询一个字段

```go
var username string
err := app.QueryObject(&username, "User.QueryUsername", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(username)
```

## 8. 删除

```go
rowsAffected, err := app.Delete("User.Delete", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(rowsAffected)
```
