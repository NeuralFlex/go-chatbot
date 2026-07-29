package utils

import "strings"

const attachmentPrefix = "Attached file: "
const csvFence = "\"\"\"\n"

func BuildMessageWithFile(filename, csvContent, userPrompt string) string {
	var b strings.Builder
	b.WriteString(userPrompt)
	b.WriteString("\n\n")
	b.WriteString(attachmentPrefix)
	b.WriteString(filename)
	b.WriteString("\n")
	b.WriteString(csvFence)
	b.WriteString(csvContent)
	b.WriteString("\n")
	b.WriteString(csvFence)
	return b.String()
}

func SplitAttachment(content string) (filename, prompt string, ok bool) {
	markerStart := strings.LastIndex(content, attachmentPrefix)
	if markerStart == -1 {
		return "", content, false
	}
	prompt = strings.TrimSuffix(content[:markerStart], "\n\n")
	rest := content[markerStart+len(attachmentPrefix):]

	nameEnd := strings.IndexByte(rest, '\n')
	if nameEnd == -1 {
		return "", content, false
	}
	filename, rest = rest[:nameEnd], rest[nameEnd+1:]

	rest, ok = strings.CutPrefix(rest, csvFence)
	if !ok {
		return "", content, false
	}
	closeIdx := strings.LastIndex(rest, csvFence)
	if closeIdx == -1 {
		return "", content, false
	}

	if closeIdx+len(csvFence) != len(rest) {
		return "", content, false
	}
	return filename, prompt, true
}
