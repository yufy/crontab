package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorhill/cronexpr"
	"github.com/yufy/crontab/internal/app/services"
	"github.com/yufy/crontab/internal/pkg/model"
	"go.uber.org/zap"
)

type JobController struct {
	logger  *zap.Logger
	service services.JobService
}

func NewJobController(logger *zap.Logger, service services.JobService) *JobController {
	return &JobController{
		logger:  logger.With(zap.String("type", "controllers.Job")),
		service: service,
	}
}

func (c *JobController) List(ctx *gin.Context) {
	c.logger.Info("list all jobs")

	jobs, err := c.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func (c *JobController) Create(ctx *gin.Context) {
	c.logger.Info("create a new crontab job")

	job := new(model.Job)
	if err := ctx.ShouldBindJSON(job); err != nil {
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		} else {
			trans, _ := ctx.Value("trans").(ut.Translator)
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": verrs.Translate(trans)})
		}
		return
	}
	if _, err := cronexpr.Parse(job.Expr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	oldJob, err := c.service.Save(ctx.Request.Context(), job)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "msg": "success", "data": oldJob})
}

func (c *JobController) Delete(ctx *gin.Context) {
	c.logger.Info("delete a exist job")

	name := ctx.Param("name")

	oldJob, err := c.service.Delete(ctx.Request.Context(), name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": oldJob})
}

func (c *JobController) Update(ctx *gin.Context) {
	c.logger.Info("update job attribute")

	name := ctx.Param("name")
	job := new(model.Job)
	job.Name = name

	if err := ctx.ShouldBindJSON(job); err != nil {
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		} else {
			trans, _ := ctx.Value("trans").(ut.Translator)
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": verrs.Translate(trans)})
		}
		return
	}
	if _, err := cronexpr.Parse(job.Expr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	oldJob, err := c.service.Save(ctx.Request.Context(), job)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": oldJob})
}

func (c *JobController) Kill(ctx *gin.Context) {
	c.logger.Info("kill a running job")

	name := ctx.Param("name")

	if err := c.service.Kill(ctx.Request.Context(), name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"code": 0, "msg": "success"})
}
