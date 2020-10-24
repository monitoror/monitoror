import TileStatus from '@/enums/tileStatus'
import TileAuthor from '@/types/tileAuthor'
import TileMergeRequest from '@/types/tileMergeRequest'

type TileBuild = {
  previousStatus?: TileStatus,
  id?: string,
  branch?: string,
  mergeRequest?: TileMergeRequest,
  author?: TileAuthor,
  estimatedDuration?: number,
  startedAt?: number,
  finishedAt?: number,
}

export default TileBuild
