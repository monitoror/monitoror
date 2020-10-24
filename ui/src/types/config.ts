import TileConfig from '@/types/tileConfig'

type Config = {
  version: string,
  columns: number,
  zoom?: number,
  tiles: TileConfig[],
}

export default Config
