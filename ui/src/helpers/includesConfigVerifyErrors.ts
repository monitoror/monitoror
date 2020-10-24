import ConfigErrorId from '@/enums/configErrorId'
import ConfigError from '@/types/configError'

const CONFIG_NON_VERIFY_ERRORS = [
  ConfigErrorId.CannotBeFetched,
  ConfigErrorId.ConfigNotFound,
  ConfigErrorId.UnknownNamedConfig,
]

export default function includesConfigVerifyErrors(errors: ConfigError[]): boolean {
  return errors.filter((error) => !CONFIG_NON_VERIFY_ERRORS.includes(error.id)).length > 0
}
