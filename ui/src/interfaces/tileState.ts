import TileStatus from '@/enums/tileStatus'
import TileValueUnit from '@/enums/tileValueUnit'
import TileAuthor from '@/interfaces/tileAuthor'

export default interface TileState {
  label?: string,
  status: TileStatus,
  previousStatus?: TileStatus,
  message?: string,
  values?: number[],
  unit?: TileValueUnit,
  author?: TileAuthor,
  duration?: number,
  estimatedDuration?: number,
  startedAt?: number,
  finishedAt?: number,
}
