import extractFieldValue from '@/helpers/extractFieldValue'
import ConfigError from '@/types/configError'

export default function getTileDocUrl(error: ConfigError): string | undefined {
  const tileType = extractFieldValue(error.data.configExtract as string, 'type')

  if (tileType === undefined) {
    return
  }

  const url = 'https://monitoror.com/documentation/#tile-' + JSON.parse(tileType).toLowerCase()

  return url
}
