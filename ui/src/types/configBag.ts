import Config from '@/types/config'
import ConfigError from '@/types/configError'

type ConfigBag = {
  config?: Config,
  errors?: ConfigError[],
}

export default ConfigBag
