// Package controller — swag_video.go
//
// 本文件仅作为 swaggo/swag 注解的宿主，函数体为空，实际路由处理由
// controller.RelayTask / controller.RelayTaskFetch / controller.VideoProxy 承担。
// 在 controller 层保留空壳函数，可让业务代码完全不受 swagger 注解污染。
package controller

import (
	"github.com/gin-gonic/gin"
)

// VideoErrorResponse 视频相关接口的错误响应结构。
// 与实际 relay/proxy 路径返回的 `{"error": {...}}` 保持一致。
type VideoErrorResponse struct {
	Error VideoErrorPayload `json:"error"`
}

// VideoErrorPayload 错误详情。
type VideoErrorPayload struct {
	Message string `json:"message" example:"task_id is required"`
	Type    string `json:"type" example:"invalid_request_error"`
}

// VideoGenerationsRequest `/v1/video/generations` 和 `/v1/videos` 的请求体。
// 对应运行时 relaycommon.TaskSubmitReq（relay/common/relay_info.go），保持文档与实现一致。
// 未在此列出的 `metadata` 子字段会被各渠道 adaptor 透传到上游（例如 Doubao 支持
// resolution/ratio/seed/camera_fixed/generate_audio/watermark/callback_url 等）。
type VideoGenerationsRequest struct {
	Prompt   string         `json:"prompt" binding:"required" example:"宇航员在月面慢跑，背景是地球"`
	Model    string         `json:"model,omitempty" example:"doubao-seedance-2-0-260128"`
	Mode     string         `json:"mode,omitempty" example:"multi_reference"` // Doubao: first_frame / first_last_frame / multi_reference
	Image    string         `json:"image,omitempty"`                          // 兼容单图；images 为空时会自动复制到 images[0]
	Images   []string       `json:"images,omitempty"`                         // 参考图 URL 列表
	Audios   []string       `json:"audios,omitempty"`                         // 参考音频 URL 列表（仅 AudioVideoRef 模型生效）
	Videos   []string       `json:"videos,omitempty"`                         // 参考视频 URL 列表（仅 AudioVideoRef 模型生效）
	Size     string         `json:"size,omitempty" example:"720x1280"`
	Duration int            `json:"duration,omitempty" example:"5"`
	Seconds  string         `json:"seconds,omitempty" example:"5"`     // 字符串秒数；>0 时覆盖 duration
	Metadata map[string]any `json:"metadata,omitempty"`                // 透传给上游渠道的扩展字段
}

// ============================================================================
// OpenAI 兼容视频 API
// ============================================================================

// VideoCreate godoc
//
//	@Summary		创建视频生成任务 (OpenAI 兼容)
//	@Description	OpenAI 兼容的视频生成接口，提交一个异步任务。
//	@Description	参考：https://platform.openai.com/docs/api-reference/videos/create
//	@Description	Sora 渠道同时接受 `multipart/form-data`（支持 `input_reference` 文件上传）与 JSON。
//	@Tags			Video
//	@Accept			json
//	@Accept			mpfd
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		VideoGenerationsRequest	true	"视频生成请求参数"
//	@Success		200		{object}	dto.OpenAIVideo			"任务对象（OpenAI 视频对象格式）"
//	@Failure		400		{object}	VideoErrorResponse		"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse		"未授权"
//	@Failure		403		{object}	VideoErrorResponse		"无权限"
//	@Failure		500		{object}	VideoErrorResponse		"服务器内部错误"
//	@Router			/v1/videos [post]
func VideoCreate(c *gin.Context) {}

// VideoFetch godoc
//
//	@Summary		查询视频任务状态 (OpenAI 兼容)
//	@Description	根据 task_id 查询 OpenAI 兼容视频生成任务的状态与结果。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			task_id	path		string			true	"Task ID"
//	@Success		200		{object}	dto.OpenAIVideo	"任务对象"
//	@Failure		401		{object}	VideoErrorResponse	"未授权"
//	@Failure		404		{object}	VideoErrorResponse	"任务不存在"
//	@Failure		500		{object}	VideoErrorResponse	"服务器内部错误"
//	@Router			/v1/videos/{task_id} [get]
func VideoFetch(c *gin.Context) {}

// VideoRemix godoc
//
//	@Summary		视频重混 (OpenAI 兼容)
//	@Description	基于已有 video_id 的视频提交重混任务。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			video_id	path		string					true	"源视频 ID"
//	@Param			request		body		VideoGenerationsRequest	true	"重混请求参数"
//	@Success		200			{object}	dto.OpenAIVideo			"任务对象"
//	@Failure		400			{object}	VideoErrorResponse		"请求参数错误"
//	@Failure		401			{object}	VideoErrorResponse		"未授权"
//	@Failure		500			{object}	VideoErrorResponse		"服务器内部错误"
//	@Router			/v1/videos/{video_id}/remix [post]
func VideoRemix(c *gin.Context) {}

// VideoProxyContent godoc
//
//	@Summary		拉取视频内容
//	@Description	根据 task_id 将上游视频内容透传给客户端。支持用户会话或 API token 鉴权。
//	@Description	成功时响应体为视频二进制流，`Content-Type` 由上游决定（通常为 `video/mp4`）。
//	@Tags			Video
//	@Security		ApiKeyAuth
//	@Produce		application/octet-stream
//	@Param			task_id	path		string			true	"Task ID"
//	@Success		200		{file}		binary			"视频二进制流"
//	@Failure		400		{object}	VideoErrorResponse	"task_id 缺失或任务未完成"
//	@Failure		401		{object}	VideoErrorResponse	"未授权"
//	@Failure		403		{object}	VideoErrorResponse	"URL 被 SSRF 策略拦截"
//	@Failure		404		{object}	VideoErrorResponse	"任务不存在"
//	@Failure		502		{object}	VideoErrorResponse	"上游获取视频失败"
//	@Router			/v1/videos/{task_id}/content [get]
func VideoProxyContent(c *gin.Context) {}

// ============================================================================
// 自研视频 API（/v1/video/generations）
// ============================================================================

// VideoGenerations godoc
//
//	@Summary		生成视频
//	@Description	调用视频生成接口生成视频，支持多种上游：
//	@Description	- 可灵 AI (Kling): https://app.klingai.com/cn/dev/document-api/apiReference/commonInfo
//	@Description	- 即梦 (Jimeng): https://www.volcengine.com/docs/85621/1538636
//	@Description	- 豆包 (Doubao-Seedance): https://www.volcengine.com/docs/82379/1520757
//	@Description	请求体字段对应 relaycommon.TaskSubmitReq；`metadata` 里的字段会被各渠道 adaptor 按需透传到上游。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		VideoGenerationsRequest	true	"视频生成请求参数"
//	@Success		200		{object}	dto.OpenAIVideo			"任务对象（Doubao/Sora 等渠道统一 OpenAI Video 对象格式；部分渠道可能返回兼容的 task_id 包装）"
//	@Failure		400		{object}	VideoErrorResponse		"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse		"未授权"
//	@Failure		403		{object}	VideoErrorResponse		"无权限"
//	@Failure		500		{object}	VideoErrorResponse		"服务器内部错误"
//	@Router			/v1/video/generations [post]
func VideoGenerations(c *gin.Context) {}

// VideoGenerationsTaskId godoc
//
//	@Summary		查询视频任务
//	@Description	根据任务 ID 查询视频生成任务的状态与结果。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			task_id	path		string					true	"Task ID"
//	@Success		200		{object}	dto.VideoTaskResponse	"任务状态和结果"
//	@Failure		400		{object}	VideoErrorResponse			"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse			"未授权"
//	@Failure		404		{object}	VideoErrorResponse			"任务不存在"
//	@Failure		500		{object}	VideoErrorResponse			"服务器内部错误"
//	@Router			/v1/video/generations/{task_id} [get]
func VideoGenerationsTaskId(c *gin.Context) {}

// ============================================================================
// 可灵 AI (Kling) 官方兼容路由
// ============================================================================

// KlingText2VideoRequest 可灵文生视频请求体。
type KlingText2VideoRequest struct {
	ModelName      string              `json:"model_name,omitempty" example:"kling-v1"`
	Prompt         string              `json:"prompt" binding:"required" example:"A cat playing piano in the garden"`
	NegativePrompt string              `json:"negative_prompt,omitempty" example:"blurry, low quality"`
	CfgScale       float64             `json:"cfg_scale,omitempty" example:"0.7"`
	Mode           string              `json:"mode,omitempty" example:"std"`
	CameraControl  *KlingCameraControl `json:"camera_control,omitempty"`
	AspectRatio    string              `json:"aspect_ratio,omitempty" example:"16:9"`
	Duration       string              `json:"duration,omitempty" example:"5"`
	CallbackURL    string              `json:"callback_url,omitempty" example:"https://your.domain/callback"`
	ExternalTaskId string              `json:"external_task_id,omitempty" example:"custom-task-001"`
}

// KlingCameraControl 可灵运镜控制。
type KlingCameraControl struct {
	Type   string             `json:"type,omitempty" example:"simple"`
	Config *KlingCameraConfig `json:"config,omitempty"`
}

// KlingCameraConfig 可灵运镜参数。
type KlingCameraConfig struct {
	Horizontal float64 `json:"horizontal,omitempty" example:"2.5"`
	Vertical   float64 `json:"vertical,omitempty" example:"0"`
	Pan        float64 `json:"pan,omitempty" example:"0"`
	Tilt       float64 `json:"tilt,omitempty" example:"0"`
	Roll       float64 `json:"roll,omitempty" example:"0"`
	Zoom       float64 `json:"zoom,omitempty" example:"0"`
}

// KlingImage2VideoRequest 可灵图生视频请求体。
type KlingImage2VideoRequest struct {
	ModelName      string              `json:"model_name,omitempty" example:"kling-v2-master"`
	Image          string              `json:"image" binding:"required" example:"https://h2.inkwai.com/bs2/upload-ylab-stunt/se/ai_portal_queue_mmu_image_upscale_aiweb/3214b798-e1b4-4b00-b7af-72b5b0417420_raw_image_0.jpg"`
	Prompt         string              `json:"prompt,omitempty" example:"A cat playing piano in the garden"`
	NegativePrompt string              `json:"negative_prompt,omitempty" example:"blurry, low quality"`
	CfgScale       float64             `json:"cfg_scale,omitempty" example:"0.7"`
	Mode           string              `json:"mode,omitempty" example:"std"`
	CameraControl  *KlingCameraControl `json:"camera_control,omitempty"`
	AspectRatio    string              `json:"aspect_ratio,omitempty" example:"16:9"`
	Duration       string              `json:"duration,omitempty" example:"5"`
	CallbackURL    string              `json:"callback_url,omitempty" example:"https://your.domain/callback"`
	ExternalTaskId string              `json:"external_task_id,omitempty" example:"custom-task-002"`
}

// KlingText2VideoGenerations godoc
//
//	@Summary		可灵文生视频
//	@Description	调用可灵 AI 文生视频接口提交异步任务。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		KlingText2VideoRequest	true	"文生视频请求参数"
//	@Success		200		{object}	dto.VideoTaskResponse	"任务状态和结果"
//	@Failure		400		{object}	VideoErrorResponse			"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse			"未授权"
//	@Failure		500		{object}	VideoErrorResponse			"服务器内部错误"
//	@Router			/kling/v1/videos/text2video [post]
func KlingText2VideoGenerations(c *gin.Context) {}

// KlingImage2VideoGenerations godoc
//
//	@Summary		可灵图生视频
//	@Description	调用可灵 AI 图生视频接口提交异步任务。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		KlingImage2VideoRequest	true	"图生视频请求参数"
//	@Success		200		{object}	dto.VideoTaskResponse	"任务状态和结果"
//	@Failure		400		{object}	VideoErrorResponse			"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse			"未授权"
//	@Failure		500		{object}	VideoErrorResponse			"服务器内部错误"
//	@Router			/kling/v1/videos/image2video [post]
func KlingImage2VideoGenerations(c *gin.Context) {}

// KlingText2videoTaskId godoc
//
//	@Summary		可灵任务查询 (文生视频)
//	@Description	根据 task_id 查询可灵文生视频任务状态。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			task_id	path		string					true	"Task ID"
//	@Success		200		{object}	dto.VideoTaskResponse	"任务状态和结果"
//	@Failure		401		{object}	VideoErrorResponse			"未授权"
//	@Failure		404		{object}	VideoErrorResponse			"任务不存在"
//	@Router			/kling/v1/videos/text2video/{task_id} [get]
func KlingText2videoTaskId(c *gin.Context) {}

// KlingImage2videoTaskId godoc
//
//	@Summary		可灵任务查询 (图生视频)
//	@Description	根据 task_id 查询可灵图生视频任务状态。
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			task_id	path		string					true	"Task ID"
//	@Success		200		{object}	dto.VideoTaskResponse	"任务状态和结果"
//	@Failure		401		{object}	VideoErrorResponse			"未授权"
//	@Failure		404		{object}	VideoErrorResponse			"任务不存在"
//	@Router			/kling/v1/videos/image2video/{task_id} [get]
func KlingImage2videoTaskId(c *gin.Context) {}

// ============================================================================
// 即梦 (Jimeng) 官方兼容路由
// ============================================================================

// JimengAction godoc
//
//	@Summary		即梦视频接口 (官方兼容)
//	@Description	即梦官方 API 统一入口，通过 query 参数 `Action` + `Version` 路由到提交或查询动作。
//	@Description	- `Action=CVSync2AsyncSubmitTask&Version=2022-08-31` 提交任务
//	@Description	- `Action=CVSync2AsyncGetResult&Version=2022-08-31` 查询结果
//	@Description	参考：https://www.volcengine.com/docs/85621/1538636
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Action	query		string	true	"操作类型"	Enums(CVSync2AsyncSubmitTask, CVSync2AsyncGetResult)
//	@Param			Version	query		string	true	"API 版本"	example(2022-08-31)
//	@Param			request	body		map[string]interface{}	true	"即梦原生请求体（字段视 Action 而定）"
//	@Success		200		{object}	dto.TaskDto	"任务状态和结果"
//	@Failure		400		{object}	VideoErrorResponse	"请求参数错误"
//	@Failure		401		{object}	VideoErrorResponse	"未授权"
//	@Router			/jimeng/ [post]
func JimengAction(c *gin.Context) {}

