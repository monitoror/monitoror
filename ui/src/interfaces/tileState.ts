import TileStatus from '@/enums/tileStatus'
import TileType from '@/enums/tileType'
import TileBuild from '@/interfaces/tileBuild'
import TileValue from '@/interfaces/tileValue'

export default interface TileState {
  type: TileType,
  status: TileStatus,
  label?: string,
  message?: string,
  value?: TileValue,
  build?: TileBuild
}
