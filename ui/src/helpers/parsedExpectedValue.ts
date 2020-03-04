import extractFieldValue from '@/helpers/extractFieldValue'

export default function parsedExtractFieldValue(jsonString: string, fieldName: string): string | undefined {
  const fieldValue = extractFieldValue(jsonString, fieldName)

  if (fieldValue === undefined) {
    return
  }

  return JSON.parse(fieldValue)
}
