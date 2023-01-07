package mediatype

var AAC, _ = New("audio/aac", "", "aac")
var ACSM, _ = New("application/vnd.adobe.adept+xml", "Adobe Content Server Message", "acsm")
var AIFF, _ = New("audio/aiff", "", "aiff")
var AVI, _ = New("video/x-msvideo", "", "avi")
var AVIF, _ = New("image/avif", "", "avif")
var Binary, _ = New("application/octet-stream", "", "")
var BMP, _ = New("image/bmp", "Bitmap Image File", "bmp")
var CBZ, _ = New("application/vnd.comicbook+zip", "Comic Book Archive", "cbz")
var CSS, _ = New("text/css", "Cascading Style Sheets", "css")
var Divina, _ = New("application/divina+zip", "Digital Visual Narratives", "divina")
var DivinaManifest, _ = New("application/divina+json", "Digital Visual Narratives", "json")
var EPUB, _ = New("application/epub+zip", "EPUB", "epub")
var GIF, _ = New("image/gif", "", "gif")
var GZ, _ = New("application/gzip", "", "gz")
var HTML, _ = New("text/html", "Hypertext Markup Language", "html")
var JavaScript, _ = New("text/javascript", "JavaScript", "js")
var JPEG, _ = New("image/jpeg", "", "jpeg")
var JSON, _ = New("application/json", "JSON", "json")
var JXL, _ = New("image/jxl", "JPEG XL", "jxl")
var LCPLicenseDocument, _ = New("application/vnd.readium.lcp.license.v1.0+json", "LCP License", "lcpl")
var LCPProtectedAudiobook, _ = New("application/audiobook+lcp", "LCP Protected Audiobook", "lcpa")
var LCPProtectedPDF, _ = New("application/pdf+lcp", "LCP Protected PDF", "lcpdf")
var LCPStatusDocument, _ = New("application/vnd.readium.license.status.v1.0+json", "LCP Status Document", "")
var LPF, _ = New("application/lpf+zip", "Lightweight Packaging Format", "lpf")
var MP3, _ = New("audio/mpeg", "", "mp3")
var MPEG, _ = New("video/mpeg", "", "mpeg")
var NCX, _ = New("application/x-dtbncx+xml", "Navigation Control File", "ncx")
var OGG, _ = New("audio/ogg", "", "oga")
var OGV, _ = New("video/ogg", "", "ogv")
var OPDS1, _ = New("application/atom+xml;profile=opds-catalog", "", "")
var OPDS1Entry, _ = New("application/atom+xml;type=entry;profile=opds-catalog", "", "")
var OPDS2, _ = New("application/opds+json", "", "")
var OPDS2Publication, _ = New("application/opds-publication+json", "", "")
var OPDSAuthentication, _ = New("application/opds-authentication+json", "", "")
var OPUS, _ = New("audio/opus", "", "opus")
var OTF, _ = New("font/otf", "OpenType Font", "otf")
var PDF, _ = New("application/pdf", "PDF", "pdf")
var PNG, _ = New("image/png", "Portable Network Graphics", "png")
var ReadiumAudiobook, _ = New("application/audiobook+zip", "Readium Audiobook", "audiobook")
var ReadiumAudiobookManifest, _ = New("application/audiobook+json", "Readium Audiobook", "json")
var ReadiumWebpub, _ = New("application/webpub+zip", "Readium Web Publication", "webpub")
var ReadiumWebpubManifest, _ = New("application/webpub+json", "Readium Web Publication", "json")
var SMIL, _ = New("application/smil+xml", "Synchronized Multimedia Integration Language", "smil")
var SVG, _ = New("image/svg+xml", "Scalable Vector Graphics", "svg")
var Text, _ = New("text/plain", "Text", "txt")
var TIFF, _ = New("image/tiff", "", "tiff")
var TTF, _ = New("font/ttf", "TrueType Font", "ttf")
var W3CWPUBManifest, _ = New("application/x.readium.w3c.wpub+json", "Web Publication", "json") // non-existent
var WAV, _ = New("audio/wav", "", "wav")
var WEBMAudio, _ = New("audio/webm", "", "webm")
var WEBMVideo, _ = New("video/webm", "", "webm")
var WEBP, _ = New("image/webp", "", "webp")
var WOFF, _ = New("font/woff", "", "woff")
var WOFF2, _ = New("font/woff2", "", "woff2")
var XHTML, _ = New("application/xhtml+xml", "", "xhtml")
var XML, _ = New("application/xml", "Xtensible Markup Language", "xml")
var ZAB, _ = New("application/x.readium.zab+zip", "Zipped Audio Book", "zab") // non-existent
var ZIP, _ = New("application/zip", "ZIP Archive", "zip")