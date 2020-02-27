import TileConfig from '@/interfaces/tileConfig'

export default interface Config {
  version: number,
  columns: number,
  zoom?: number,
  tiles: TileConfig[],
}
