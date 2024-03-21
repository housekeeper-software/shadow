package controller

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"jingxi.cn/tools/shadow/internal/pkg/db"
)

type Controller struct {
	srv   *http.Server
	entry *db.Entry
}

type ResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NoResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "404, page not exists!",
	})
}

func Response(ctx *gin.Context, code int, message string) {
	if code != http.StatusOK {
		ctx.AbortWithStatusJSON(code, ResponseHeader{
			Code:    code,
			Message: message,
		})
	} else {
		ctx.JSON(code, ResponseHeader{
			Code:    code,
			Message: message,
		})
	}
}

func NewController() *Controller {
	return &Controller{
		srv:   nil,
		entry: db.NewEntry(),
	}
}

func (c *Controller) Run(httpAddr string, dir string) error {
	start := time.Now()
	err := c.entry.LoadEntry(dir)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	fmt.Println(duration)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s %+v %s %d %s %s",
				param.Method,
				param.Request.Header,
				param.Request.Proto,
				param.StatusCode,
				param.Request.UserAgent(),
				param.ErrorMessage)
		}))
	}
	v1 := router.Group("/api/v1")
	{
		v1.GET("/removeEntry", c.removeEntryHandlerFunc)
		v1.POST("/addItem", c.addItemHandlerFunc)
		v1.GET("/removeItem", c.removeItemHandlerFunc)
		v1.GET("/getIndexes", c.getIndexesHandlerFunc)
		v1.GET("/getEntry", c.getEntryHandlerFunc)
	}

	router.NoRoute(NoResponse)

	c.srv = &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
	logrus.Infof("http server listen on: %s", httpAddr)
	if err := c.srv.ListenAndServe(); err != nil {
		logrus.Errorf("gin ListenAndServe(%s) error: %+v", httpAddr, err)
		return err
	}
	return c.entry.Close()
}

func (c *Controller) Stop() {
	if c.srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.srv.Shutdown(ctx); err != nil {
		logrus.Errorf("gin Shutdown error: %+v", err)
	}
}

func (c *Controller) removeEntryHandlerFunc(ctx *gin.Context) {
	entry := ctx.Query("entry")
	if len(entry) < 1 {
		Response(ctx, http.StatusBadRequest, "Bad request")
		return
	}
	err := c.entry.RemoveEntry(entry)
	if err != nil {
		Response(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Response(ctx, http.StatusOK, "Success")
}

func (c *Controller) addItemHandlerFunc(ctx *gin.Context) {
	entry := ctx.Query("entry")
	name := ctx.Query("name")
	if len(entry) < 1 || len(name) < 1 {
		Response(ctx, http.StatusBadRequest, "Bad request")
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		Response(ctx, http.StatusBadRequest, "Bad request")
		return
	}

	// Restore the io.ReadCloser to its original state
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = c.entry.AddItem(entry, name, body)
	if err != nil {
		Response(ctx, http.StatusNotFound, err.Error())
		return
	}
	Response(ctx, http.StatusOK, "Success")
}

func (c *Controller) removeItemHandlerFunc(ctx *gin.Context) {
	entry := ctx.Query("entry")
	name := ctx.Query("name")
	id := ctx.Query("id")
	if len(entry) < 1 || len(name) < 1 {
		Response(ctx, http.StatusBadRequest, "Bad request")
		return
	}
	err := c.entry.RemoveItem(entry, name, id)
	if err != nil {
		Response(ctx, http.StatusNotFound, err.Error())
		return
	}
	Response(ctx, http.StatusOK, "Success")
}

func (c *Controller) getIndexesHandlerFunc(ctx *gin.Context) {
	j, err := c.entry.IndexToJson()
	if err != nil {
		ctx.Data(http.StatusOK, "application/json", []byte("{}"))
		return
	}
	ctx.Data(http.StatusOK, "application/json", j)
}

func (c *Controller) getEntryHandlerFunc(ctx *gin.Context) {
	entry := ctx.Query("entry")
	if len(entry) < 1 {
		Response(ctx, http.StatusBadRequest, "Bad request")
		return
	}
	j, err := c.entry.ToJson(entry)
	if err != nil {
		Response(ctx, http.StatusNotFound, err.Error())
		return
	}
	ctx.Data(http.StatusOK, "application/json", j)
}
