import TileValueUnit from '@/enums/tileValueUnit'
import TileMetrics from '@/types/tileMetrics'
import TileState from '@/types/tileState'
import {computed, ComputedRef} from 'vue'

export default function useTileMetrics(state: ComputedRef<TileState | undefined>) {
  const metrics = computed((): TileMetrics | undefined => {
    if (state.value === undefined) {
      return
    }

    return state.value.metrics
  })

  const unit = computed((): TileValueUnit => {
    if (metrics.value === undefined) {
      return TileValueUnit.Raw
    }

    return metrics.value.unit
  })

  const values = computed((): string[] | undefined => {
    if (metrics.value === undefined) {
      return
    }

    return metrics.value.values
  })

  const displayedMetric = computed((): string | undefined => {
    if (values.value === undefined) {
      return
    }

    const UNIT_DISPLAY = {
      [TileValueUnit.Millisecond]: 'ms',
      [TileValueUnit.Ratio]: '%',
      [TileValueUnit.Number]: '',
      [TileValueUnit.Raw]: '',
    }

    let value = values.value[values.value.length - 1]
    if (unit.value === TileValueUnit.Millisecond) {
      value = Math.round(parseFloat(value)).toString()
    } else if (unit.value === TileValueUnit.Ratio) {
      value = (parseFloat(value) * 100).toFixed(2).toString()
    }

    return value + UNIT_DISPLAY[unit.value]
  })

  return {
    displayedMetric,
  }
}
