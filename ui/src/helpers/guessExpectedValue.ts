import parsedExtractFieldValue from '@/helpers/parsedExpectedValue'
import StringSimilarity from 'string-similarity'

export default function guessExpectedValue(
  configExtract: string,
  fieldName: string,
  expectedValues: string[],
): string | undefined {
  const currentValue = parsedExtractFieldValue(configExtract, fieldName)

  if (currentValue === undefined) {
    return
  }

  const bestMatch = StringSimilarity.findBestMatch(currentValue, expectedValues).bestMatch

  return bestMatch.target
}
