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

	r := gin.Default()
	r.Use(CORS()) // 默认跨域
	// 示例：保留一个测试路由
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "服务运行正常")
	})
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
	r.Run(":8050")
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
