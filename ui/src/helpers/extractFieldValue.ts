export default function extractFieldValue(jsonString: string, fieldName: string): string | undefined {
  try {
    return JSON.stringify(JSON.parse(jsonString)[fieldName])
  } catch (e) {
    return
  }
}
