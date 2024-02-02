package main

import (
	"divine/kms/kms"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IssueRewardsRequest struct {
	Users   []string `json:"users"`
	Amounts []uint64 `json:"amounts"`
}

type CreateTokenRequest struct {
	User           string `json:"user"`
	ID             string `json:"id"`
	AdditionalData string `json:"id"'`
	Fee            uint64 `json:"fee"`
}

var db = make(map[string]string)

func setupRouter(k *kms.KMS) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.POST("/rewards/issue", func(c *gin.Context) {
		var requestInput IssueRewardsRequest
		c.BindJSON(&requestInput)
		if len(requestInput.Amounts) != len(requestInput.Users) {
			c.AbortWithStatus(400)
		}

		if err := k.IssueRewards(requestInput.Amounts, requestInput.Users); err != nil {
			c.AbortWithStatus(500)
		}
	})

	r.POST("/nft/create", func(c *gin.Context) {
		var requestInput CreateTokenRequest
		c.BindJSON(&requestInput)

		user, err := k.GetUser(requestInput.User)
		if err != nil {
			c.AbortWithStatus(500)
		}

		if user == nil {
			c.AbortWithStatus(400)
		}

		if err := k.CreateToken(requestInput.User, requestInput.ID, requestInput.Fee, requestInput.AdditionalData); err != nil {
			c.AbortWithStatus(500)
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	k := kms.NewKMS()
	r := setupRouter(k)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
