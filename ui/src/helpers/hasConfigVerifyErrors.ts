import ConfigErrorId from '@/enums/configErrorId'
import ConfigError from '@/interfaces/configError'

const CONFIG_NON_VERIFY_ERRORS = [
  ConfigErrorId.CannotBeFetched,
  ConfigErrorId.ConfigNotFound,
]

export default function hasConfigVerifyErrors(errors: ConfigError[]): boolean {
  return errors.filter((error) => !CONFIG_NON_VERIFY_ERRORS.includes(error.id)).length > 0
}
