import ConfigErrorId from '@/enums/configErrorId'

export default interface ConfigError {
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
