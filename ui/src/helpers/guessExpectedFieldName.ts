import StringSimilarity from 'string-similarity'

export default function guessExpectedFieldName(
  fieldName: string,
  expectedValues: string[],
): string | undefined {
  const bestMatch = StringSimilarity.findBestMatch(fieldName, expectedValues).bestMatch.target

  if (fieldName === bestMatch) {
    return
  }

  return bestMatch
}
