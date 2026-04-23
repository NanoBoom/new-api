package doubao

var ModelList = []string{
	"doubao-seedance-1-0-pro-250528",
	"doubao-seedance-1-0-lite-t2v",
	"doubao-seedance-1-0-lite-i2v",
	"doubao-seedance-1-5-pro-251215",
	"doubao-seedance-2-0-260128",
	"doubao-seedance-2-0-fast-260128",
}

var ChannelName = "doubao-video"

// Doubao 图生视频的三种互斥模式（TaskSubmitReq.mode 取值）。
const (
	ModeFirstFrame     = "first_frame"
	ModeFirstLastFrame = "first_last_frame"
	ModeMultiReference = "multi_reference"
)

// ContentItem.role 取值。
const (
	RoleFirstFrame     = "first_frame"
	RoleLastFrame      = "last_frame"
	RoleReferenceImage = "reference_image"
	RoleReferenceAudio = "reference_audio"
	RoleReferenceVideo = "reference_video"
)

// ModelProfile 汇总单个模型在生视频场景下的能力与计费系数。
// 新增 / 调整模型时，只需修改下方 modelProfiles 表。
//
//   - ImageFrames:     支持 role=first_frame / last_frame
//   - ReferenceImage:  支持 role=reference_image（多模态参考生视频）
//   - AudioVideoRef:   支持 role=reference_audio / reference_video
//   - VideoInputRatio: 含视频输入时的费率折扣（含视频单价 / 不含视频单价），0 表示无折扣
type ModelProfile struct {
	ImageFrames     bool
	ReferenceImage  bool
	AudioVideoRef   bool
	VideoInputRatio float64
}

// modelProfiles 是模型能力的唯一事实来源。
var modelProfiles = map[string]ModelProfile{
	"doubao-seedance-1-0-pro-250528":  {ImageFrames: true},
	"doubao-seedance-1-0-lite-t2v":    {}, // 仅文生视频，无任何图片/音视频 role
	"doubao-seedance-1-0-lite-i2v":    {ImageFrames: true, ReferenceImage: true},
	"doubao-seedance-1-5-pro-251215":  {ImageFrames: true},
	"doubao-seedance-2-0-260128":      {ImageFrames: true, ReferenceImage: true, AudioVideoRef: true, VideoInputRatio: 28.0 / 46.0},
	"doubao-seedance-2-0-fast-260128": {ImageFrames: true, ReferenceImage: true, AudioVideoRef: true, VideoInputRatio: 22.0 / 37.0},
}

// ProfileOf 返回模型能力档案，未知模型返回零值（禁用所有特殊 role）。
func ProfileOf(modelName string) ModelProfile {
	return modelProfiles[modelName]
}

// GetVideoInputRatio 返回"含视频输入"场景下的费率折扣。
func GetVideoInputRatio(modelName string) (float64, bool) {
	if p, ok := modelProfiles[modelName]; ok && p.VideoInputRatio > 0 {
		return p.VideoInputRatio, true
	}
	return 0, false
}
