package mindmapservice

import (
	"context"
	"fmt"
	"forge/pkg/log/zlog"
	"github.com/unidoc/unioffice/v2/document"
	"github.com/unidoc/unioffice/v2/presentation"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func ParseFile(ctx context.Context, fh *multipart.FileHeader) (text string, err error) {
	mime := fileHeaderMime(fh)

	switch mime {
	case "application/pdf":
		// PDF 处理
		text, err = extractPDF(fh)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to extract pdf: %v", err)
			return
		}
	case
		"application/msword",                                                      // .doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": // .docx
		// Word 处理
		text, err = extractWord(fh)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to extract word: %v", err)
			return
		}
	case
		"application/vnd.ms-powerpoint",                                             // .ppt
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": // .pptx
		// PPT 处理
		text, err = extractPPT(fh)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to extract PPT: %v", err)
			return
		}

	default:
		// 其他类型（或报错、或放行）
		zlog.CtxErrorf(ctx, "Unsupported file type for parsing")
		return "", fmt.Errorf(`unknown mime: "%s"`, mime)
	}
	return

}

// 返回形如 "image/jpeg"、"application/zip" 的 MIME 类型，出错时返回空串
func fileHeaderMime(fh *multipart.FileHeader) string {
	rc, err := fh.Open()
	if err != nil {
		return ""
	}
	defer rc.Close()

	// 只取文件头 512 字节即可
	buf := make([]byte, 512)
	n, _ := io.ReadFull(rc, buf)
	if n == 0 {
		return ""
	}
	return http.DetectContentType(buf[:n])
}

func extractPDF(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder

	fmt.Printf("PDF to text extraction:\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum) //文本操作对象
		if err != nil {
			return "", err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		text, err := ex.ExtractText() //文本
		if err != nil {
			return "", err
		}
		//拼接
		textBuilder.WriteString("------------------------------\n")
		textBuilder.WriteString(fmt.Sprintf("Page %d:\n", pageNum))
		textBuilder.WriteString(fmt.Sprintf("\"%s\"\n", text))
		textBuilder.WriteString("------------------------------")

	}

	return textBuilder.String(), nil
}

func extractWord(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	doc, err := document.Read(f, fh.Size) //word文件对象
	if err != nil {
		return "", err
	}

	var allTest strings.Builder

	// To extract the text and work with the formatted info in a simple fashion, you can use:
	extracted := doc.ExtractText()
	for _, e := range extracted.Items {
		allTest.WriteString(e.Text)
	}
	return allTest.String(), nil
}

func extractPPT(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	ppt, err := presentation.Read(f, fh.Size)
	if err != nil {
		return "", err
	}
	pt := ppt.ExtractText()
	var allText strings.Builder
	for _, slide := range pt.Slides { //每个  slide  代表一张幻灯片
		for _, item := range slide.Items { //当前这页ppt中的文本项的列表
			allText.WriteString(item.Text)
			allText.WriteString("\n")
		}
	}
	return allText.String(), nil
}
