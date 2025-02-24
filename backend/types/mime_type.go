package types

type (
	MIMEType = string
)

const (
	MIME_AAC    MIMEType = "audio/aac"
	MIME_ABW    MIMEType = "application/x-abiword"
	MIME_APNG   MIMEType = "image/apng"
	MIME_ARC    MIMEType = "application/x-freearc"
	MIME_AVIF   MIMEType = "image/avif"
	MIME_AVI    MIMEType = "video/x-msvideo"
	MIME_AZW    MIMEType = "application/vnd.amazon.ebook"
	MIME_BIN    MIMEType = "application/octet-stream"
	MIME_BMP    MIMEType = "image/bmp"
	MIME_BZ     MIMEType = "application/x-bzip"
	MIME_BZ2    MIMEType = "application/x-bzip2"
	MIME_CDA    MIMEType = "application/x-cdf"
	MIME_CSH    MIMEType = "application/x-csh"
	MIME_CSS    MIMEType = "text/css"
	MIME_CSV    MIMEType = "text/csv"
	MIME_DOC    MIMEType = "application/msword"
	MIME_DOCX   MIMEType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MIME_EOT    MIMEType = "application/vnd.ms-fontobject"
	MIME_EPUB   MIMEType = "application/epub+zip"
	MIME_GZ     MIMEType = "application/gzip"
	MIME_GIF    MIMEType = "image/gif"
	MIME_HTML   MIMEType = "text/html"
	MIME_ICON   MIMEType = "image/vnd.microsoft.icon"
	MIME_ICS    MIMEType = "text/calendar"
	MIME_JAR    MIMEType = "application/java-archive"
	MIME_JPEG   MIMEType = "image/jpeg"
	MIME_JS     MIMEType = "text/javascript"
	MIME_JSON   MIMEType = "application/json"
	MIME_JSONLD MIMEType = "application/ld+json"
	MIME_MIDI   MIMEType = "audio/midi"
	MIME_MP3    MIMEType = "audio/mpeg"
	MIME_MP4    MIMEType = "video/mp4"
	MIME_MPEG   MIMEType = "video/mpeg"
	MIME_MPKG   MIMEType = "application/vnd.apple.installer+xml"
	MIME_ODP    MIMEType = "application/vnd.oasis.opendocument.presentation"
	MIME_ODS    MIMEType = "application/vnd.oasis.opendocument.spreadsheet"
	MIME_ODT    MIMEType = "application/vnd.oasis.opendocument.text\n"
	MIME_OGA    MIMEType = "audio/ogg"
	MIME_OGV    MIMEType = "video/ogg"
	MIME_OGX    MIMEType = "application/ogg"
	MIME_OPUS   MIMEType = "audio/opus"
	MIME_OTF    MIMEType = "font/otf"
	MIME_PNG    MIMEType = "image/png"
	MIME_PDF    MIMEType = "application/pdf"
	MIME_PHP    MIMEType = "application/x-httpd-php"
	MIME_PPT    MIMEType = "application/vnd.ms-powerpoint"
	MIME_PPTX   MIMEType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MIME_RAR    MIMEType = "application/vnd.rar"
	MIME_RTF    MIMEType = "application/rtf"
	MIME_SH     MIMEType = "application/x-sh"
	MIME_SVG    MIMEType = "image/svg+xml"
	MIME_TAR    MIMEType = "application/x-tar"
	MIME_TIFF   MIMEType = "image/tiff"
	MIME_TS     MIMEType = "video/mp2t"
	MIME_TTF    MIMEType = "font/ttf"
	MIME_TXT    MIMEType = "text/plain"
	MIME_VSD    MIMEType = "application/vnd.visio"
	MIME_WAV    MIMEType = "audio/wav"
	MIME_WEBA   MIMEType = "audio/webm"
	MIME_WEBM   MIMEType = "video/webm"
	MIME_WEBP   MIMEType = "image/webp"
	MIME_WOFF   MIMEType = "font/woff"
	MIME_WOFF2  MIMEType = "font/woff2"
	MIME_XHTML  MIMEType = "application/xhtml+xml"
	MIME_XLS    MIMEType = "application/vnd.ms-excel"
	MIME_XLSX   MIMEType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MIME_XML    MIMEType = "application/xml"
	MIME_XUL    MIMEType = "application/vnd.mozilla.xul+xml"
	MIME_ZIP    MIMEType = "application/zip"
	MIME_7Z     MIMEType = "application/x-7z-compressed"
)

var (
	extMimeType = map[string]MIMEType{
		".aac":    MIME_AAC,
		".apng":   MIME_APNG,
		".arc":    MIME_ARC,
		".avif":   MIME_AVIF,
		".avi":    MIME_AVI,
		".azm":    MIME_AZW,
		".bin":    MIME_BIN,
		".bmp":    MIME_BMP,
		".bz":     MIME_BZ,
		".bz2":    MIME_BZ2,
		".cda":    MIME_CDA,
		".csv":    MIME_CSV,
		".doc":    MIME_DOC,
		".docx":   MIME_DOCX,
		".eot":    MIME_EOT,
		".epub":   MIME_EPUB,
		".gz":     MIME_GZ,
		".gif":    MIME_GIF,
		".ico":    MIME_ICON,
		".jpeg":   MIME_JPEG,
		".jpg":    MIME_JPEG,
		".json":   MIME_JSON,
		".jsonld": MIME_JSONLD,
		".mid":    MIME_MIDI,
		".midi":   MIME_MIDI,
		".mp3":    MIME_MP3,
		".mp4":    MIME_MP4,
		".mpeg":   MIME_MPEG,
		".odp":    MIME_ODP,
		".ods":    MIME_ODS,
		".odt":    MIME_ODT,
		".oga":    MIME_OGA,
		".ogv":    MIME_OGV,
		".ogx":    MIME_OGX,
		".opus":   MIME_OPUS,
		".otf":    MIME_OTF,
		".png":    MIME_PNG,
		".pdf":    MIME_PDF,
		".ppt":    MIME_PPT,
		".pptx":   MIME_PPTX,
		".rar":    MIME_RAR,
		".rtf":    MIME_RTF,
		".svg":    MIME_SVG,
		".tar":    MIME_TAR,
		".tif":    MIME_TIFF,
		".tiff":   MIME_TIFF,
		".ts":     MIME_TS,
		".ttf":    MIME_TTF,
		".txt":    MIME_TXT,
		".vsd":    MIME_VSD,
		".wav":    MIME_WAV,
		".weba":   MIME_WEBA,
		".webm":   MIME_WEBM,
		".webp":   MIME_WEBP,
		".woff":   MIME_WOFF,
		".woff2":  MIME_WOFF2,
		".xls":    MIME_XLS,
		".xlsx":   MIME_XLSX,
		".xml":    MIME_XML,
		".zip":    MIME_ZIP,
		".7z":     MIME_7Z,
	}
)

func Ext2MimeType(ext string) MIMEType {
	res, ok := extMimeType[ext]
	if ok {
		return res
	}
	// 默认 text
	return MIME_TXT
}
