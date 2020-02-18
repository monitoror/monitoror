import TileType from '@/enums/tileType'

export default interface TileConfig {
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
