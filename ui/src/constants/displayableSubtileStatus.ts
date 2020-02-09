import TileStatus from '@/enums/tileStatus'

const DISPLAYABLE_SUBTILE_STATUS = [
  TileStatus.Warning,
  TileStatus.Failed,
  TileStatus.Queued,
  TileStatus.Running,
  TileStatus.Canceled,
  TileStatus.ActionRequired,
]

export default DISPLAYABLE_SUBTILE_STATUS
