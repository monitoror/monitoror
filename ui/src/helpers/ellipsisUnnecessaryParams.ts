export default function ellipsisUnnecessaryParams(jsonString: string): string {
  const cleanedJsonString = jsonString.replace(/"params":[^}]*}/g, '"params": { ... }')

  return cleanedJsonString
}
