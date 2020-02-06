import TileStatus from '@/enums/tileStatus'

const ORDERED_TILE_STATUS = [
  TileStatus.Unknown,
  TileStatus.Success,
  TileStatus.Canceled,
  TileStatus.Warning,
  TileStatus.Failed,
  TileStatus.Queued,
  TileStatus.Running,
]

export default function mostImportantStatus(status1: TileStatus, status2: TileStatus): TileStatus {
  return ORDERED_TILE_STATUS.indexOf(status1) < ORDERED_TILE_STATUS.indexOf(status2) ? status2 : status1
}
