package app

import (
	"github.com/gin-gonic/gin"
	"github.com/yufy/crontab/internal/app/controllers"
)

func RegisterRouter(
	job *controllers.JobController,
) InitRouter {
	return func(r *gin.Engine) {
		apiv1 := r.Group("/api/v1/")
		{
			apiv1.GET("jobs", job.List)
			apiv1.POST("jobs", job.Create)
			apiv1.DELETE("jobs/:name", job.Delete)
			apiv1.PUT("jobs/:name", job.Update)
			apiv1.POST("jobs/:name", job.Kill)
		}
	}
}
