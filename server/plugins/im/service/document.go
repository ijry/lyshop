package service

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

// SupportedDocExts lists the file extensions ExtractText can handle.
var SupportedDocExts = []string{
	".txt", ".md", ".markdown", ".text", ".log",
	".csv", ".tsv",
	".json", ".xml", ".html", ".htm",
	".docx",
	".pdf",
	".xlsx",
}

// IsSupportedDoc reports whether a filename has a supported extension.
func IsSupportedDoc(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, e := range SupportedDocExts {
		if e == ext {
			return true
		}
	}
	return false
}

// ExtractText reads raw document bytes and returns plain UTF-8 text. The
// extraction strategy is chosen by file extension. Unknown extensions are
// treated as UTF-8 plain text.
func ExtractText(filename string, data []byte) (string, error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".docx":
		return extractDocx(data)
	case ".pdf":
		return extractPDF(data)
	case ".xlsx":
		return extractXlsx(data)
	case ".csv":
		return extractDelimited(data, ',')
	case ".tsv":
		return extractDelimited(data, '\t')
	case ".html", ".htm", ".xml":
		return stripTags(string(data)), nil
	default:
		// .txt/.md/.json/.log and anything else: trust UTF-8 plain text.
		return string(data), nil
	}
}

// extractDocx pulls text from a .docx (Office Open XML) archive by reading
// word/document.xml and stripping its tags. Paragraph and break tags are
// converted to newlines so the slicer can split on structure.
func extractDocx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("无法读取 docx：%w", err)
	}
	var docXML []byte
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			docXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", err
			}
			break
		}
	}
	if docXML == nil {
		return "", fmt.Errorf("docx 缺少 word/document.xml")
	}
	xml := string(docXML)
	// Paragraph / line-break boundaries → newline before tags are stripped.
	xml = regexp.MustCompile(`(?i)</w:p>`).ReplaceAllString(xml, "\n")
	xml = regexp.MustCompile(`(?i)<w:br\s*/?>`).ReplaceAllString(xml, "\n")
	return stripTags(xml), nil
}

// extractPDF concatenates the plain text of every page.
func extractPDF(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("无法读取 pdf：%w", err)
	}
	var b strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		txt, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		b.WriteString(txt)
		b.WriteString("\n")
	}
	return b.String(), nil
}

// extractXlsx flattens every sheet into "cell\tcell" rows, one row per line.
func extractXlsx(data []byte) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("无法读取 xlsx：%w", err)
	}
	defer f.Close()
	var b strings.Builder
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		b.WriteString("# " + sheet + "\n")
		for _, row := range rows {
			b.WriteString(strings.Join(row, "\t"))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return b.String(), nil
}

// extractDelimited renders CSV/TSV rows as tab-joined lines.
func extractDelimited(data []byte, comma rune) (string, error) {
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = comma
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	var b strings.Builder
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Fall back to raw text on malformed input.
			return string(data), nil
		}
		b.WriteString(strings.Join(rec, "\t"))
		b.WriteString("\n")
	}
	return b.String(), nil
}

var (
	reTag      = regexp.MustCompile(`(?s)<[^>]*>`)
	reScript   = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	reStyle    = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	reSpaces   = regexp.MustCompile(`[ \t]+`)
	reNewlines = regexp.MustCompile(`\n{3,}`)
)

// stripTags removes HTML/XML tags (and script/style bodies) and normalises
// whitespace, yielding readable plain text.
func stripTags(s string) string {
	s = reScript.ReplaceAllString(s, " ")
	s = reStyle.ReplaceAllString(s, " ")
	s = reTag.ReplaceAllString(s, " ")
	s = htmlUnescape(s)
	s = reSpaces.ReplaceAllString(s, " ")
	s = reNewlines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// htmlUnescape decodes the handful of entities common in extracted text.
func htmlUnescape(s string) string {
	repl := strings.NewReplacer(
		"&nbsp;", " ", "&amp;", "&", "&lt;", "<", "&gt;", ">",
		"&quot;", "\"", "&#39;", "'", "&apos;", "'",
	)
	return repl.Replace(s)
}

// SliceText splits text into overlapping chunks suitable for one knowledge
// entry each. chunkSize/overlap are measured in runes (so CJK counts correctly).
// Splitting prefers paragraph (blank-line) then sentence boundaries, falling
// back to a hard rune cut so no chunk exceeds chunkSize.
func SliceText(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		chunkSize = 500
	}
	if overlap < 0 || overlap >= chunkSize {
		overlap = chunkSize / 10
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// Build paragraph units first.
	paras := splitParagraphs(text)

	var chunks []string
	var cur []rune
	flush := func() {
		if len(cur) == 0 {
			return
		}
		chunks = append(chunks, strings.TrimSpace(string(cur)))
		// Carry an overlap tail into the next chunk for context continuity.
		if overlap > 0 && len(cur) > overlap {
			cur = append([]rune{}, cur[len(cur)-overlap:]...)
		} else {
			cur = cur[:0]
		}
	}

	for _, p := range paras {
		pr := []rune(p)
		// A single paragraph larger than the budget is hard-split.
		if len(pr) > chunkSize {
			if len(strings.TrimSpace(string(cur))) > 0 {
				flush()
			}
			for start := 0; start < len(pr); start += (chunkSize - overlap) {
				end := start + chunkSize
				if end > len(pr) {
					end = len(pr)
				}
				chunks = append(chunks, strings.TrimSpace(string(pr[start:end])))
				if end == len(pr) {
					break
				}
			}
			cur = cur[:0]
			continue
		}
		// Would adding this paragraph overflow the current chunk? Flush first.
		if len(cur)+len(pr)+1 > chunkSize && len(cur) > 0 {
			flush()
		}
		if len(cur) > 0 {
			cur = append(cur, '\n', '\n')
		}
		cur = append(cur, pr...)
	}
	if len(strings.TrimSpace(string(cur))) > 0 {
		chunks = append(chunks, strings.TrimSpace(string(cur)))
	}

	// Drop empties produced by overlap/trim.
	out := chunks[:0]
	for _, c := range chunks {
		if strings.TrimSpace(c) != "" {
			out = append(out, c)
		}
	}
	return out
}

// splitParagraphs splits on blank lines; lines themselves are kept intact.
func splitParagraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	raw := regexp.MustCompile(`\n[ \t]*\n`).Split(text, -1)
	var out []string
	for _, p := range raw {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
