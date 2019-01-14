package model

type MediaConfig struct {
	SaveTmp     bool    // 保留合成中生成临时文件
	UseOldFile  bool    // 当目标文件已存在时，不再生成新的文件
	IsTs        bool    // 最终文件是否为ts文件
	Rate        int     // 帧率
	Size        string  // 生成视频文件的分辨率
	Width       float64 // 生成视频文件的宽
	Height      float64 // 生成视频文件的高
	Scale       float64 // 生成视频文件的宽高比
	Profile     string  // 配置属性profile
	Level       string  // 配置属性level
	Pix         string  // 配置属性fix_fmt
	IsAddSuffix bool    // 是否在结尾追加推广片段
}

type DownloadConfig struct {
	CoverStyle    string  // 覆盖类型
	CoverQrcode   bool    // 是否覆盖水印图片
	CoverImageKey string  // 覆盖水印图objectKey
	Transparent   string  // 覆盖水印图名都
	CoverLocation string  // 覆盖水印位置
	XInstance     string  // x轴偏移量
	YInstance     string  // y轴偏移量
	WaterWS       float64 // 水印图相对于原图的宽度比例
	WaterHS       float64 // 水印图相对于原图的高度比例
}
