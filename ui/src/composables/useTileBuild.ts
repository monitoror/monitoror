import TileStatus from '@/enums/tileStatus'
import TileType from '@/enums/tileType'
import TileAuthor from '@/interfaces/tileAuthor'
import TileBuild from '@/interfaces/tileBuild'
import TileMergeRequest from '@/interfaces/tileMergeRequest'
import {RootState} from '@/store'
import {addSeconds, differenceInSeconds, format, formatDistance} from 'date-fns'
import {ComputedRef, computed} from 'vue'
import {useStore} from 'vuex'

export default function useTileBuild(type: ComputedRef<TileType>, status: ComputedRef<TileStatus | undefined>, build: ComputedRef<TileBuild | undefined>) {
  const store = useStore<RootState>()

  const now = computed((): Date => {
    return store.state.now
  })

  const mergeRequestLabelPrefix = computed((): string | undefined => {
    if (mergeRequest.value === undefined) {
      return
    }

    let mergeRequestPrefix = 'MR'
    if (type.value === TileType.GitHubPullRequest) {
      mergeRequestPrefix = 'PR'
    }

    return mergeRequestPrefix + '#' + mergeRequest.value.id
  })

  const branch = computed((): string | undefined => {
    if (build.value === undefined) {
      return
    }

    return build.value.branch
  })

  const mergeRequest = computed((): TileMergeRequest | undefined => {
    if (build.value === undefined) {
      return
    }

    return build.value.mergeRequest
  })

  const startedAt = computed((): Date | undefined => {
    if (build.value === undefined || build.value.startedAt === undefined) {
      return
    }

    return new Date(build.value.startedAt)
  })

  const finishedAt = computed((): Date | undefined => {
    if (build.value === undefined || build.value.finishedAt === undefined) {
      return
    }

    return new Date(build.value.finishedAt)
  })

  const duration = computed((): number | undefined => {
    if (startedAt.value === undefined) {
      return
    }

    return differenceInSeconds(now.value, startedAt.value)
  })

  const estimatedDuration = computed((): number | undefined => {
    if (build.value === undefined) {
      return
    }

    return build.value.estimatedDuration
  })

  const progress = computed((): number | undefined => {
    if (duration.value === undefined || estimatedDuration.value === undefined) {
      return
    }

    const progress = duration.value / estimatedDuration.value * 100

    return progress
  })

  const progressTime = computed((): string | undefined => {
    if (progress.value === undefined || estimatedDuration.value === undefined || duration.value === undefined) {
      return
    }

    const totalSeconds = Math.abs(Math.round((estimatedDuration.value - duration.value)))

    const overtimePrefix = (progress.value > 100 ? 'Overtime: +' : '')
    const date = addSeconds(new Date(0), totalSeconds)
    const dateFormat = totalSeconds > 3600 ? 'hh:mm:ss' : 'mm:ss'

    return overtimePrefix + format(date, dateFormat)
  })

  const isOvertime = computed((): boolean => {
    if (progressTime.value === undefined) {
      return false
    }

    return progressTime.value.includes('+')
  })

  const progressBarStyle = computed((): { transform: string } | undefined => {
    if (progress.value === undefined) {
      return
    }

    const progressPercentage = Math.min(progress.value, 100)

    return {
      transform: `translateX(${-100 + progressPercentage}%)`,
    }
  })

  const finishedSince = computed((): string | undefined => {
    if (finishedAt.value === undefined) {
      return
    }

    return formatDistance(finishedAt.value, now.value) + ' ago'
  })

  const author = computed((): TileAuthor | undefined => {
    if (build.value === undefined) {
      return
    }

    return build.value.author
  })

  const showAuthor = computed((): boolean => {
    return author.value !== undefined && status.value === TileStatus.Failed
  })

  return {
    now,
    mergeRequestLabelPrefix,
    branch,
    mergeRequest,
    startedAt,
    finishedAt,
    duration,
    estimatedDuration,
    progress,
    progressTime,
    isOvertime,
    progressBarStyle,
    finishedSince,
    author,
    showAuthor,
  }
}
