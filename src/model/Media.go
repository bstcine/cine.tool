package model

type MediaConfig struct {
	SaveTmp bool          // 保留合成中生成临时文件
	UseOldFile bool       // 当目标文件已存在时，不再生成新的文件
	IsTs bool             // 最终文件是否为ts文件
	Rate int              // 帧率
	Size string           // 生成视频文件的分辨率
	Width float64         // 生成视频文件的宽
	Height float64        // 生成视频文件的高
	Scale float64         // 生成视频文件的宽高比
	Profile string        // 配置属性profile
	Level string          // 配置属性level
	Pix string            // 配置属性fix_fmt
}