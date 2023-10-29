package route

import (
	admin2 "fusionsms/application/admin"
	controller2 "fusionsms/application/controller"
	visit2 "fusionsms/application/visit"
	"github.com/gin-gonic/gin"
)

func Core(service *gin.Engine) {
	// 发信接口
	{
		service.POST("/local", controller2.LocalTest)
		service.POST("/send", controller2.Send)
	}
	admin := service.Group("/admin")
	{
		user := admin.Group("/user")
		user.POST("/list", admin2.List)
		user.POST("/changeBalance", admin2.ChangeBalance)
		user.GET("/detail/:uid", admin2.GetDetail)
		user.POST("/detail/:uid", admin2.ChangeDetail)
		log := admin.Group("/log")
		{
			log.POST("/get", admin2.GetLogList)
		}
		setting := admin.Group("/setting")
		{
			setting.GET("/getSent", admin2.GetSent)
			setting.POST("/setSent", admin2.SetSent)
		}
		exchange := admin.Group("/exchange")
		{
			exchange.POST("/get", admin2.GetCodeList)
			exchange.POST("/add", admin2.AddExchangeCode)
		}
	}
	visit := service.Group("/user")
	{
		visit.POST("/login", visit2.Login)
		visit.POST("/reg", visit2.Reg)
		visit.POST("/forgot", visit2.Forgot)
		visit.GET("/logout", controller2.Logout)
		visit.POST("/emailCaptcha", visit2.EmailCaptcha)
		visit.POST("/forgotCaptcha", visit2.ForgotCaptcha)
	}
	controller := service.Group("/controller")
	{
		controller.POST("/addSign", controller2.AddSign)
		controller.GET("/getSignList", controller2.GetSignList)
		controller.POST("/resetKey", controller2.ResetKey)
		controller.GET("/getSign/:sid", controller2.GetSign)
		controller.GET("/getSignContent/:sid", controller2.GetSignContent)
		controller.POST("/addSignTemplate", controller2.AddSignTemplate)
		controller.POST("/deleteSignTemplate", controller2.DeleteTemplate)
	}
	exchange := controller.Group("/exchange")
	{
		exchange.GET("/get", controller2.GetExchangeList)
		exchange.POST("/exchange", controller2.ExchangeCode)
	}
	setting := controller.Group("/setting")
	{
		setting.GET("/whiteList", controller2.WhiteList)
		setting.GET("/blackPhone", controller2.BlackPhone)
		setting.GET("/remainingLimitPrompt", controller2.RemainingLimitPrompt)
		setting.GET("/getLimit", controller2.Limit)
		setting.POST("/whiteList", controller2.WhiteListSave)
		setting.POST("/remainingLimitPrompt", controller2.RemainingLimitPromptSave)
		setting.POST("/blackPhone", controller2.BlackPhoneSave)
		setting.POST("/addLimit", controller2.AddLimit)
		setting.POST("/editLimit", controller2.EditLimit)
		setting.POST("/deleteLimit", controller2.DeleteLimit)
	}
	log := controller.Group("/log")
	{
		log.POST("/get", controller2.GetLogList)
		//log.POST("/filter", controller2.FilterLogList)
	}
	api := controller.Group("/Api")
	{
		api.GET("/getDaySendNum", controller2.DaySendNum)
		api.GET("/getDaySendErrorNum", controller2.DaySendErrorNum)
		api.GET("/getDaySendSuccessNum", controller2.DaySendSuccessNum)
		api.GET("/getPaddingSignNum", controller2.PaddingSignNum)
		api.GET("/getPaddingTemplateNum", controller2.PaddingTemplateNum)
		api.GET("/getSurplusQuotaNum", controller2.SurplusQuotaNum)
		api.GET("/getControllerIndex", controller2.ControllerIndex)
		api.GET("/getChartData", controller2.ChartData)
	}
	// 回执类
	{
		api.POST("/ReceiveSignStatus", controller2.ReceiveSignStatus)
		api.POST("/ReceiveTemplateStatus", controller2.ReceiveTemplateStatus)
		api.POST("/ReceiveSentStatus", controller2.ReceiveSentStatus)
	}
}
