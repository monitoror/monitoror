import TileStatus from '@/enums/tileStatus'
import TileAuthor from '@/interfaces/tileAuthor'
import TileMergeRequest from '@/interfaces/tileMergeRequest'

export default interface TileBuild {
  previousStatus?: TileStatus,
  id?: string,
  branch?: string,
  mergeRequest?: TileMergeRequest,
  author?: TileAuthor,
  estimatedDuration?: number,
  startedAt?: number,
  finishedAt?: number,
}
