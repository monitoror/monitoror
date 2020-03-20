import extractFieldValue from '@/helpers/extractFieldValue'
import guessExpectedFieldName from '@/helpers/guessExpectedFieldName'
import jsonSyntaxColor from '@/helpers/jsonSyntaxColor'
import splitList from '@/helpers/splitList'
import ConfigError from '@/interfaces/configError'

export default function formatConfigExtract(configError: ConfigError): string {
  let configExtract = configError.data.configExtract
  if (configExtract === undefined) {
    return ''
  }

  // Avoid parse issue due to wrong escape sequences (e.g: \s is invalid in JSON)
  // This will allow us to format the JSON, then we will revert this in the result string
  configExtract = configExtract.replace(/\\/g, '\\\\')

  let highlight = true
  if (configError.data.fieldName && configError.data.expected) {
    const guessedFieldName = guessExpectedFieldName(configError.data.fieldName, splitList(configError.data.expected))
    highlight = guessedFieldName === configError.data.fieldName
  }

  try {
    JSON.parse(configExtract)
  } catch (err) {
    return jsonSyntaxColor(configExtract.replace(/\\\\/g, '\\'))
  }

  const formattedConfigExtract = JSON.stringify(JSON.parse(configExtract), null, 2)
  let html = formattedConfigExtract.replace(/\\\\/g, '\\')

  if (highlight) {
    let configExtractHighlight = configError.data.configExtractHighlight
    let patternPrefix = ''

    if (configExtractHighlight === undefined && configError.data.fieldName !== undefined) {
      patternPrefix = `"${configError.data.fieldName}":\\s+`
      configExtractHighlight = extractFieldValue(configExtract, configError.data.fieldName)
    }

    if (configExtractHighlight !== undefined) {
      let formattedConfigExtractHighlight = configExtractHighlight
      try {
        formattedConfigExtractHighlight = JSON.stringify(JSON.parse(configExtractHighlight), null, 2)
      } catch (err) {
        // Escape unescaped char to be a string used as regex
        // \w => \\w
        if (formattedConfigExtractHighlight.startsWith('\\')) {
          formattedConfigExtractHighlight = '\\' + formattedConfigExtractHighlight
        }
      }
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
        return jsonSyntaxColor(html)
      }

      const match = matches[0]
      const matchRegexPrefix = match.startsWith('\\') ? '\\' : ''
      const markClassAttr = isHighlightMultiline ? 'class="multiline-mark"' : ''

      html = formattedConfigExtract.replace(matchRegexPrefix + match, `<mark ${markClassAttr}>${match}</mark>`)
    }

    if (html.includes('</mark>')) {
      html = `<span class="has-mark">${html}</span>`
    }
  }

  return jsonSyntaxColor(html)
}
