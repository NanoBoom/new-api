package doubao

import "strings"

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

// ModelProfile 汇总单个模型在生视频场景下的能力档案。
// 新增 / 调整模型能力时，只需修改下方 modelProfiles 表；
// 计费档位单独维护在 modelPriceRatios 表中。
//
//   - ImageFrames:    支持 role=first_frame / last_frame
//   - ReferenceImage: 支持 role=reference_image（多模态参考生视频）
//   - AudioVideoRef:  支持 role=reference_audio / reference_video
type ModelProfile struct {
	ImageFrames    bool
	ReferenceImage bool
	AudioVideoRef  bool
}

// modelProfiles 是模型能力的唯一事实来源。
var modelProfiles = map[string]ModelProfile{
	"doubao-seedance-1-0-pro-250528":  {ImageFrames: true},
	"doubao-seedance-1-0-lite-t2v":    {}, // 仅文生视频，无任何图片/音视频 role
	"doubao-seedance-1-0-lite-i2v":    {ImageFrames: true, ReferenceImage: true},
	"doubao-seedance-1-5-pro-251215":  {ImageFrames: true},
	"doubao-seedance-2-0-260128":      {ImageFrames: true, ReferenceImage: true, AudioVideoRef: true},
	"doubao-seedance-2-0-fast-260128": {ImageFrames: true, ReferenceImage: true, AudioVideoRef: true},
}

// ProfileOf 返回模型能力档案，未知模型返回零值（禁用所有特殊 role）。
func ProfileOf(modelName string) ModelProfile {
	return modelProfiles[modelName]
}

// PriceTier 是定价表的主键：分辨率（归一化为小写） + 是否含视频输入。
type PriceTier struct {
	Resolution string
	HasVideo   bool
}

// modelPriceRatios 以每个模型的「基准档位」单价为基准，
// 其他档位的 ratio = 上游实际单价 / 基准单价。
// 基准档位（720p 不含视频）ratio 恒为 1.0，因此无需登记；命中基准档位
// 时 GetPriceRatio 返回 (0, false)，调用方按渠道基准价直接计费。
//
// 基准单价约定：
//   - doubao-seedance-2-0-260128      基准单价 46
//   - doubao-seedance-2-0-fast-260128 基准单价 37
//
// 注：480p 与 720p 在 Doubao 官方定价中同档；若未来拆分，只需在本表按
// 实际单价增改对应 PriceTier 条目。
var modelPriceRatios = map[string]map[PriceTier]float64{
	"doubao-seedance-2-0-260128": {
		// 不含视频输入：480p / 720p / 未指定 均为基准档位（ratio 1.0，表中省略）
		{Resolution: "1080p", HasVideo: false}: 51.0 / 46.0,
		// 含视频输入
		{Resolution: "", HasVideo: true}:       28.0 / 46.0,
		{Resolution: "480p", HasVideo: true}:   28.0 / 46.0,
		{Resolution: "720p", HasVideo: true}:   28.0 / 46.0,
		{Resolution: "1080p", HasVideo: true}:  31.0 / 46.0,
	},
	"doubao-seedance-2-0-fast-260128": {
		{Resolution: "", HasVideo: true}:     22.0 / 37.0,
		{Resolution: "480p", HasVideo: true}: 22.0 / 37.0,
		{Resolution: "720p", HasVideo: true}: 22.0 / 37.0,
	},
}

// GetPriceRatio 按模型、分辨率、是否含视频输入查询价位系数。
// 未命中（即基准档位或未声明档位）返回 (0, false)，调用方不附加 OtherRatio。
func GetPriceRatio(modelName, resolution string, hasVideo bool) (float64, bool) {
	tiers, ok := modelPriceRatios[modelName]
	if !ok {
		return 0, false
	}
	r, ok := tiers[PriceTier{Resolution: strings.ToLower(resolution), HasVideo: hasVideo}]
	if !ok || r <= 0 {
		return 0, false
	}
	return r, true
}
