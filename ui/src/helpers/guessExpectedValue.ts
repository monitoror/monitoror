import parsedExtractFieldValue from '@/helpers/parsedExpectedValue'
import {findBestMatch} from 'string-similarity'

export default function guessExpectedValue(
  configExtract: string,
  fieldName: string,
  expectedValues: string[],
): string | undefined {
  const currentValue = parsedExtractFieldValue(configExtract, fieldName)

  if (currentValue === undefined) {
    return
  }

  const bestMatch = findBestMatch(currentValue, expectedValues).bestMatch

  return bestMatch.target
}
