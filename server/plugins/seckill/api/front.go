package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	seckillservice "github.com/ijry/lyshop/server/plugins/seckill/service"
)

// RegisterFrontRoutes 注册前端路由
func RegisterFrontRoutes(g *gin.RouterGroup) {
	// 秒杀商品列表
	g.GET("/seckill/products", listSeckillProducts)
	g.GET("/seckill/products/:id", getSeckillProductDetail)

	// 秒杀活动列表
	g.GET("/seckill/activities", listSeckillActivities)
}

// listSeckillProducts 秒杀商品列表
func listSeckillProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Query("activity_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := seckillservice.ListFrontProducts(c.Request.Context(), activityID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  products,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// getSeckillProductDetail 秒杀商品详情
func getSeckillProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := seckillservice.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": product})
}

// listSeckillActivities 秒杀活动列表
func listSeckillActivities(c *gin.Context) {
	activities, err := seckillservice.GetActiveActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": activities})
}
