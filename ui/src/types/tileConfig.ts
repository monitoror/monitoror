import TileType from '@/enums/tileType'

type TileConfig = {
  id: string,
  stateKey: string,
  type: TileType,
  label?: string,
  columnSpan?: number,
  rowSpan?: number,
  url?: string,
  tiles?: TileConfig[],
  initialMaxDelay?: number,
}

export default TileConfig
