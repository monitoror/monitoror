import includesConfigVerifyErrors from '@/helpers/includesConfigVerifyErrors'
import ConfigError from '@/types/configError'
import {computed} from 'vue'
import {useStore} from 'vuex'

export default function useErrors() {
  const store = useStore()

  const errors = computed((): ConfigError[] => {
    return store.state.errors
  })

  const hasErrors = computed((): boolean => {
    return errors.value.length > 0
  })

  const hasConfigVerifyErrors = computed((): boolean => {
    return includesConfigVerifyErrors(errors.value)
  })

  return {
    errors,
    hasErrors,
    hasConfigVerifyErrors,
  }
}
