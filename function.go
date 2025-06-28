package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {

	// 测试大模型功能
	//testLLM()
	ConnectMyDatabase()
	r := gin.Default()
	r.Use(CORS()) // 默认跨域
	// 示例：保留一个测试路由
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "服务运行正常")
	})

	//下面实现登录和注册功能
	{
		r.POST("/Register", func(c *gin.Context) {
			var dataReceive struct {
				Account  string `json:"account"`
				Password string `json:"password"`
				Username string `json:"username"`
			}
			errflag1 := c.BindJSON(&dataReceive)
			if errflag1 != nil {
				fmt.Println(errflag1)
				println("从请求体数据中获取json失败")
			}
			//检查是否有重账号名的账户
			var countnum int
			database.QueryRow("select  count(*) from register where username=?", dataReceive.Username).Scan(&countnum)
			if countnum != 0 {
				c.String(http.StatusOK, "账号名重复")
			} else {
				_, errflag2 := database.Exec("insert into register(account,password,username) values(?,?,?)", dataReceive.Account, dataReceive.Password, dataReceive.Username)
				if errflag2 != nil {
					fmt.Println(errflag2)
					println("数据库新建条目失败")
				}
			}
			//go asyncRegister(c)
		})

		// 用户登录
		r.POST("/login", func(c *gin.Context) {
			var datareceive struct {
				Account  string `json:"account"`
				Password string `json:"password"`
			}
			errflag1 := c.BindJSON(&datareceive)
			if errflag1 != nil {
				fmt.Println(errflag1)
				println("从请求体中获取json数据失败")
			}
			//先查看是否有该用户
			var countnum int
			database.QueryRow("select count(*) from register where account=?", datareceive.Account).Scan(&countnum)
			if countnum == 0 {
				c.String(http.StatusOK, "账户名错误（没有注册或输入错误）")
			} else {
				var password string
				database.QueryRow("select password from register where account=?", datareceive.Account).Scan(&password)
				if password == datareceive.Password {
					c.String(http.StatusOK, "登录成功")
				} else {
					c.String(http.StatusOK, "密码错误")
				}
			}

		})
	}

	// 大模型调用路由
	r.POST("/LLM", func(c *gin.Context) {
		// 调用语言模型函数
		var datareceive struct {
			SystemPrompt string                   `json:"system"`
			Messages     []map[string]interface{} `json:"messages"`
		}
		var data2send struct {
			Response string `json:"response"`
		}
		if err := c.BindJSON(&datareceive); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}
		response, err := askBaiduQianfan(datareceive.SystemPrompt, datareceive.Messages)
		if err != nil {
			fmt.Printf("API错误: %v\n", err)
		}
		data2send.Response = response
		fmt.Printf("大模型测试响应: %s", response)

		c.JSON(http.StatusOK, &data2send)
	})

	// 启动服务器
	r.Run(":8080")
}

// CORS中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
