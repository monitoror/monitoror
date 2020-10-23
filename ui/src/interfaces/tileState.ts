import TileStatus from '@/enums/tileStatus'
import TileType from '@/enums/tileType'
import TileBuild from '@/interfaces/tileBuild'
import TileMetrics from '@/interfaces/tileMetrics'

export default interface TileState {
  type: TileType,
  status: TileStatus,
  label?: string,
  message?: string,
  metrics?: TileMetrics,
  build?: TileBuild
}
