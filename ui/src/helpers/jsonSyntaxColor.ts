export default function jsonSyntaxColor(jsonString: string): string {
  const coloredJson = jsonString
    .replace(/:\s+"(.*)"/g, ': <span class="code-string">"$1"</span>')
    .replace(/:\s+([.\d]+)/g, ': <span class="code-number">$1</span>')

  return coloredJson
}
