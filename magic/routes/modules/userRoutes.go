package routes

import (
	"forestwhisperGo/magic/controller"
)

func init() { Register(func(r *gin.Engine){
group := r.Group("/api/v1/user/")
{
// v1.GET("/users", controller.getUserByUsername)
// v1.GET("/me", controller.getCurrentUserInfo)
v1.POST("/login", controller.Login)
v1.POST("/regter", controller.Register)
v1.POST("/email/code/regter", controller.sendEmailCodeRegister)
// v1.POST("/users/del", controller.deleteUserByUid)
// v1.POST("/me/upd", controller.updateUserEmail)
// v1.POST("/passwd", controller.changePassword)
}})}
