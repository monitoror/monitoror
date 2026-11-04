export default function getQueryParamValue(
  queryParamName: string,
  defaultValue?: string,
): string | undefined {
  const queryParams = window.location.search.substr(1).split('&')

  let value = defaultValue
  const valueQueryParam = queryParams.find((queryParam: string) => {
    return new RegExp(`^${queryParamName}=`).test(queryParam)
  })
  if (valueQueryParam) {
    value = valueQueryParam.substr(valueQueryParam.indexOf('=') + 1)
    try {
      value = decodeURIComponent(value.replace(/\+/g, ' '))
    } catch (e) {
      // Malformed encoding: keep the raw value
    }
  }

  return value
}
