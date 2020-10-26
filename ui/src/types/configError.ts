import ConfigErrorId from '@/enums/configErrorId'

type ConfigError = {
  id: ConfigErrorId,
  message: string,
  data: {
    configExtract?: string,
    configExtractHighlight?: string,
    value?: string,
    fieldName?: string,
    expected?: string,
  },
}

export default ConfigError
