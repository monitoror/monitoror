import {addSeconds, differenceInSeconds, format, formatDistance} from 'date-fns'
import Vue from 'vue'
import {Prop} from 'vue-property-decorator'

import TileStatus from '@/enums/tileStatus'
import TileType from '@/enums/tileType'
import TileAuthor from '@/interfaces/tileAuthor'
import TileBuild from '@/interfaces/tileBuild'
import TileConfig from '@/interfaces/tileConfig'
import TileState from '@/interfaces/tileState'

export default abstract class AbstractMonitororTile extends Vue {
  /*
   * Props
   */

  @Prop()
  protected config!: TileConfig

  /*
   * Computed
   */

  get type(): TileType {
    return this.config.type
  }

  get stateKey(): string {
    return this.config.stateKey
  }

  get theme(): string {
    return this.$store.getters.theme.toString().toLowerCase()
  }

  get now(): Date {
    return this.$store.state.now
  }

  get state(): TileState | undefined {
    if (!this.$store.state.tilesState.hasOwnProperty(this.stateKey)) {
      return
    }

    return this.$store.state.tilesState[this.stateKey]
  }

  get label(): string | undefined {
    if (this.config.label) {
      return this.config.label
    }

    if (this.state === undefined) {
      return
    }

    return this.state.label
  }

  get build(): TileBuild | undefined {
    if (this.state === undefined) {
      return
    }

    return this.state.build
  }

  get branch(): string | undefined {
    if (this.build === undefined) {
      return
    }

    return this.build.branch
  }

  get status(): string | undefined {
    if (this.state === undefined) {
      return
    }

    return this.state.status
  }

  get previousStatus(): string | undefined {
    if (this.build === undefined) {
      return
    }

    return this.build.previousStatus
  }

  get isQueued(): boolean {
    return this.status === TileStatus.Queued
  }

  get isRunning(): boolean {
    return this.status === TileStatus.Running
  }

  get isSucceeded(): boolean {
    if (this.isQueued || this.isRunning) {
      return this.previousStatus === TileStatus.Success
    }

    return this.status === TileStatus.Success
  }

  get isFailed(): boolean {
    if (this.isQueued || this.isRunning) {
      return this.previousStatus === TileStatus.Failed
    }

    return this.status === TileStatus.Failed
  }

  get isWarning(): boolean {
    if (this.isQueued || this.isRunning) {
      return this.previousStatus === TileStatus.Warning
    }

    return this.status === TileStatus.Warning
  }

  get startedAt(): Date | undefined {
    if (this.build === undefined || this.build.startedAt === undefined) {
      return
    }

    return new Date(this.build.startedAt)
  }

  get finishedAt(): Date | undefined {
    if (this.build === undefined || this.build.finishedAt === undefined) {
      return
    }

    return new Date(this.build.finishedAt)
  }

  get duration(): number | undefined {
    if (this.startedAt === undefined) {
      return
    }

    return differenceInSeconds(this.now, this.startedAt)
  }

  get estimatedDuration(): number | undefined {
    if (this.build === undefined) {
      return
    }

    return this.build.estimatedDuration
  }

  get progress(): number | undefined {
    if (this.duration === undefined || this.estimatedDuration === undefined) {
      return
    }

    const progress = this.duration / this.estimatedDuration * 100

    return progress
  }

  get progressTime(): string | undefined {
    if (!this.progress || this.estimatedDuration === undefined || this.duration === undefined) {
      return
    }

    const totalSeconds = Math.abs(Math.round((this.estimatedDuration - this.duration)))

    const overtimePrefix = (this.progress > 100 ? 'Overtime: +' : '')
    const date = addSeconds(new Date(0), totalSeconds)
    const dateFormat = totalSeconds > 3600 ? 'hh:mm:ss' : 'mm:ss'

    return overtimePrefix + format(date, dateFormat)
  }

  get isOvertime(): boolean {
    if (this.progressTime === undefined) {
      return false
    }

    return this.progressTime.includes('+')
  }

  get progressBarStyle() {
    if (!this.progress) {
      return
    }

    const progress = Math.min(this.progress, 100)

    return {
      transform: `translateX(${-100 + progress}%)`,
    }
  }

  get finishedSince(): string | undefined {
    if (this.finishedAt === undefined) {
      return
    }

    return formatDistance(this.finishedAt, this.now) + ' ago'
  }

  get author(): TileAuthor | undefined {
    if (this.build === undefined) {
      return
    }

    return this.build.author
  }

  get showAuthor(): boolean {
    return this.author !== undefined && this.status === TileStatus.Failed
  }
}
