# gobatis 

模仿MyBatis开发的一个go版本，实现了一些基本功能，并做了一些简化，支持 trim，where，set，if，foreach 标签，具体使用请查看 sample.xml。

使用类库:
* XML解析库：github.com/beevik/etree
* 数据库：github.com/jmoiron/sqlx
* 日志库：go.uber.org/zap
* 表达式解析库：github.com/antonmedv/expr

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

func main() {
	// 设置默认日志
	SetDefaultLogger()

	// 设置自定义日志，支持如：go.uber.org/zap，
	// SetCustomLogger(logger)

	// ...
}
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
if err := app.LoadFromFile("User", "/path/to/your/mapper.xml"); err != nil {
    panic(err)
}
```

## 3. 添加数据

* 传入参数类型为 map[string]interface{}，以下均相同。

```xml
<insert id="Insert">
    INSERT INTO app_user(username,password,status,create_time)
    VALUES(#{username},#{password},#{status},#{createTime})
</insert>
```
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
```xml
<update id="Update">
    UPDATE app_user
    <set>
        <if test="password != nil and password != ''">password = #{password}</if>
        <if test="status > 0">status = #{status}</if>
        <if>update_time = #{updateTime}</if>
    </set>
    WHERE id = #{id}
</update>
```
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
		UpdateTime int64 `db:"udpate_time"`
	}
)
```
```xml
<select id="Query">
    SELECT id,username,status,create_time,update_time FROM app_user WHERE id = #{id}
</select>
```
```go
var user User
err := app.QueryObject(&user, "User.Query", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(user)
```

## 6. 查询多条记录
```xml
<select id="List">
    SELECT id,username,status,create_time,update_time FROM app_user
    <where>
        <if test="status != nil and status > 0">status = #{status}</if>
        <if test="username != nil and username != ''">username LIKE CONCAT('%', #{username}, '%')</if>
    </where>
    ORDER BY id DESC
    <if test="limit != nil and limit > 0">LIMIT #{limit}</if>
</select>
```
```go
var users []User
err := app.QueryObjects(&users, "User.List", map[string]interface{}{"limit": 10})
fmt.Println(err)
fmt.Println(users)
```

## 7. 查询一个字段
```xml
<select id="QueryUsername">
    SELECT username FROM app_user WHERE id = #{id}
</select>
```
```go
var username string
err := app.QueryObject(&username, "User.QueryUsername", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(username)
```

## 8. 删除
```xml
<delete id="Delete">
    DELETE FROM app_user WHERE id = #{id}
</delete>
```
```go
rowsAffected, err := app.Delete("User.Delete", map[string]interface{}{"id": 1})
fmt.Println(err)
fmt.Println(rowsAffected)
```

## 9. foreach
```xml
<update id="UpdateForeach">
    UPDATE app_user SET
    <foreach collection="values" separator="," index="index" item="item">
       ${index} = #{item}
    </foreach>
    where id = #{id}
</update>
```
```go
rowsAffected, err := app.Update("User.UpdateForeach", map[string]interface{}{
    "id":     3,
    "values": map[string]interface{}{"password": "pwd-333333", "status": 2, "update_time": time.Now().Unix()},
})
fmt.Println(err)
fmt.Println(rowsAffected)
```


## 10. 设置自定义logger

第三方logger需要实现的接口，github上的大部分库已支持。
㮂具体的实现代码请查看 logger.go。

```go
type ILogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}
```

```go
import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	SetCustomLogger(logger.Sugar())
	
	// ...
}
```