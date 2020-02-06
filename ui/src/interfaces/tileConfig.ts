import TileType from '@/enums/tileType'

export default interface TileConfig {
  type: TileType,
  label?: string,
  columnSpan?: number,
  rowSpan?: number,
  url?: string,
  stateKey: string,
  tiles?: TileConfig[],
}
