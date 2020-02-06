import TileStatus from '@/enums/tileStatus'
import TileState from '@/interfaces/tileState'

export default function getSubTilePreviousOrStatus(subTileState?: TileState): TileStatus {
  if (subTileState === undefined) {
    return TileStatus.Unknown
  }

  let subTileStatus = subTileState.status

  if ([TileStatus.Queued, TileStatus.Running].includes(subTileState.status)) {
    subTileStatus = subTileState.previousStatus as TileStatus
  }

  return subTileStatus
}
