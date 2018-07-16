package model

type DownFile struct {
	Path        string
	ChapterName string
	LessonName  string
	Name        string
}

type InputArgs struct {
	OutputPath string /** 输出目录 */
	LocalPath  string /** 输入的目录或文件路径 */
	LogoPath   string /** 水印图片名称*/
}

type OSSConfig struct {
	KeyId      string // accesskeyid
	KeySecret  string // accesskeysecret
}