import TileStatus from '@/enums/tileStatus'
import TileState from '@/interfaces/tileState'

export default function getSubTilePreviousOrStatus(subTileState?: TileState): TileStatus {
  if (subTileState === undefined) {
    return TileStatus.Unknown
  }

  let subTileStatus = subTileState.status

  if ([TileStatus.Queued, TileStatus.Running].includes(subTileState.status) && subTileState.build !== undefined) {
    subTileStatus = subTileState.build.previousStatus as TileStatus
  }

  return subTileStatus
}
