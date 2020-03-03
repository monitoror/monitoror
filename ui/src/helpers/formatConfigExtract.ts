import extractFieldValue from '@/helpers/extractFieldValue'
import jsonSyntaxColor from '@/helpers/jsonSyntaxColor'
import ConfigError from '@/interfaces/configError'

export default function formatConfigExtract(configError: ConfigError): string {
  if (configError.data.configExtract === undefined) {
    return ''
  }

  const formattedConfigExtract = JSON.stringify(JSON.parse(configError.data.configExtract), null, 2)
  let html = formattedConfigExtract

  let configExtractHighlight = configError.data.configExtractHighlight
  let patternPrefix = ''

  if (configExtractHighlight === undefined && configError.data.fieldName !== undefined) {
    patternPrefix = `"${configError.data.fieldName}":\\s+`
    configExtractHighlight = extractFieldValue(configError.data.configExtract, configError.data.fieldName)
  }

  if (configExtractHighlight !== undefined) {
    const formattedConfigExtractHighlight = JSON.stringify(JSON.parse(configExtractHighlight), null, 2)
    const isHighlightMultiline = formattedConfigExtractHighlight.includes('\n')
    const multilinePrefix = isHighlightMultiline ? ' *' : ''
    const multilineSuffix = isHighlightMultiline ? ',?' : ''
    const pattern = [
      multilinePrefix,
      patternPrefix,
      formattedConfigExtractHighlight.replace(/\s+/g, '\\s+'),
      multilineSuffix,
    ].join('')
    const find = new RegExp(pattern)
    const matches = formattedConfigExtract.match(find)

    if (matches === null) {
      return html
    }

    const match = matches[0]
    const markClassAttr = isHighlightMultiline ? 'class="multiline-mark"' : ''

    html = formattedConfigExtract.replace(match, `<mark ${markClassAttr}>${match}</mark>`)
  }

  if (html.includes('</mark>')) {
    html = `<span class="has-mark">${html}</span>`
  }

  return jsonSyntaxColor(html)
}
