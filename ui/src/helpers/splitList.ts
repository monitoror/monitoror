export default function splitList(list: string): string[] {
  try {
    return list.split(',').map((item) => item.trim()).sort()
  } catch (e) {
    return [list]
  }
}
