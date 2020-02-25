import TileConfig from '@/interfaces/tileConfig'

export default interface Config {
  columns?: number,
  zoom?: number,
  tiles?: TileConfig[],
  errors?: string[],
}
