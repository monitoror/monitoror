import TileStatus from '@/enums/tileStatus'

const DISPLAYABLE_SUBTILE_STATUS = [
  TileStatus.Canceled,
  TileStatus.Warning,
  TileStatus.Failed,
  TileStatus.Queued,
  TileStatus.Running,
]

export default DISPLAYABLE_SUBTILE_STATUS
