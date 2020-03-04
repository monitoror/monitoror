import TileConfig from '@/interfaces/tileConfig'

export default interface Config {
  version: string,
  columns: number,
  zoom?: number,
  tiles: TileConfig[],
}
