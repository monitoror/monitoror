import Config from '@/interfaces/config'
import ConfigError from '@/interfaces/configError'

export default interface ConfigBag {
  config?: Config,
  errors?: ConfigError[],
}
