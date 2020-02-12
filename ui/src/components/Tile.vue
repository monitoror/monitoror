import TileValueUnit from '@/enums/tileValueUnit'
<template>
  <div class="c-monitoror-tile" :class="classes" :style="styles">
    <div class="c-monitoror-tile--content" v-if="!isEmpty">
      <div class="c-monitoror-tile--label">
        {{ label }}
      </div>

      <div class="c-monitoror-tile--message" v-if="message">
        {{ message }}
      </div>

      <div class="c-monitoror-tile--value" v-if="value">
        {{ value }}
      </div>

      <div class="c-monitoror-tile--sub-tiles" v-if="isGroup">
        <monitoror-sub-tile v-for="subTile in displayedSubTiles" :key="subTile.stateKey" :config="subTile"></monitoror-sub-tile>
      </div>

      <monitoror-tile-icon :tile-type="type" class="c-monitoror-tile--icon"></monitoror-tile-icon>

      <div class="c-monitoror-tile--author" v-if="showAuthor">
        <img :src="author.avatarUrl" alt="" class="c-monitoror-tile--author-avatar">
        {{ author.name }}
      </div>

      <div class="c-monitoror-tile--finished-at" v-if="finishedSince">
        {{ finishedSince }}
      </div>

      <template v-if="isRunning || isQueued">
        <div class="c-monitoror-tile--progress-time">
          <template v-if="isRunning">
            {{ progressTime }}
          </template>
          <template v-else>
            Pending...
          </template>
        </div>
        <div class="c-monitoror-tile--progress">
          <div class="c-monitoror-tile--progress-bar" :style="progressBarStyle"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  import {formatDistanceToNow} from 'date-fns'
  import Vue from 'vue'
  import {Component, Prop} from 'vue-property-decorator'

  import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'
  import TileStatus from '@/enums/tileStatus'
  import TileType from '@/enums/tileType'
  import TileValueUnit from '@/enums/tileValueUnit'
  import TileAuthor from '@/interfaces/tileAuthor'
  import TileConfig from '@/interfaces/tileConfig'
  import TileState from '@/interfaces/tileState'

  import MonitororSubTile from '@/components/SubTile.vue'
  import MonitororTileIcon from '@/components/TileIcon.vue'

  @Component({
    components: {
      MonitororSubTile,
      MonitororTileIcon,
    },
  })
  export default class MonitororTile extends Vue {
    /*
     * Props
     */

    @Prop()
    private config!: TileConfig

    /*
     * Computed
     */

    get classes() {
      return {
        ['c-monitoror-tile__theme-' + this.theme]: true,
        'c-monitoror-tile__empty': this.isEmpty,
        'c-monitoror-tile__group': this.isGroup,
        'c-monitoror-tile__centered-message': this.mustCenterMessage,
        'c-monitoror-tile__status-succeeded': this.isSucceeded,
        'c-monitoror-tile__status-failed': this.isFailed,
        'c-monitoror-tile__status-warning': this.isWarning,
        'c-monitoror-tile__status-running': this.isRunning,
        'c-monitoror-tile__status-queued': this.isQueued,
        'c-monitoror-tile__status-canceled': this.status === TileStatus.Canceled,
        'c-monitoror-tile__status-action-required': this.status === TileStatus.ActionRequired,
      }
    }

    get styles() {
      const styles = {
        'grid-column': `auto / span ${this.columnSpan}`,
        'grid-row': `auto / span ${this.rowSpan}`,
        '--row-span': this.rowSpan,
      }

      return styles
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

    get type(): TileType {
      return this.config.type
    }

    get isEmpty(): boolean {
      return this.config.type === TileType.Empty
    }

    get isGroup(): boolean {
      return this.config.type === TileType.Group
    }

    get mustCenterMessage(): boolean {
      return this.values !== undefined
    }

    get label(): string | undefined {
      if (this.config.label) {
        return this.config.label
      }

      if (this.state) {
        return this.state.label
      }
    }

    get columnSpan(): number {
      return this.config.columnSpan || 1
    }
    get rowSpan(): number {
      return this.config.rowSpan || 1
    }

    get url(): string | undefined {
      return this.config.url
    }

    get stateKey(): string {
      return this.config.stateKey
    }

    get theme(): string {
      return this.$store.getters.theme.toString().toLowerCase()
    }

    get displayedSubTiles(): TileConfig[] | undefined {
      if (!this.config.tiles) {
        return
      }

      const displayedSubTiles = this.config.tiles.filter((subTile) => {
        if (!this.$store.state.tilesState.hasOwnProperty(subTile.stateKey)) {
          return false
        }

        return DISPLAYABLE_SUBTILE_STATUS.includes(this.$store.state.tilesState[subTile.stateKey].status)
      })

      return displayedSubTiles
    }

    get state(): TileState | undefined {
      if (!this.$store.state.tilesState.hasOwnProperty(this.stateKey)) {
        return
      }

      return this.$store.state.tilesState[this.stateKey]
    }

    get status(): string | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.status
    }

    get previousStatus(): string | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.previousStatus
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

    get message(): string | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.message
    }

    get unit(): TileValueUnit {
      if (this.state === undefined || this.state.unit === undefined) {
        return TileValueUnit.Default
      }

      return this.state.unit as TileValueUnit
    }

    get values(): number[] | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.values
    }

    get value(): string | undefined {
      if (this.values === undefined) {
        return
      }

      const UNIT_DISPLAY = {
        [TileValueUnit.Millisecond]: 'ms',
        [TileValueUnit.Default]: '',
      }

      let value = this.values[this.values.length - 1]
      if (this.unit === TileValueUnit.Millisecond) {
        value = Math.round(value)
      }

      return value + UNIT_DISPLAY[this.unit]
    }

    get finishedAt(): number | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.finishedAt
    }

    get duration(): number | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.duration
    }

    get estimatedDuration(): number | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.estimatedDuration
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

      let totalSeconds = Math.round((this.estimatedDuration - this.duration))
      if (this.progress > 100) {
        totalSeconds = Math.round((this.duration - this.estimatedDuration))
      }
      const hours = Math.floor(totalSeconds / 3600)
      const minutes = Math.floor((totalSeconds - (hours * 3600)) / 60)
      const seconds = totalSeconds - (hours * 3600) - (minutes * 60)
      let minutesPrefix = ''
      let secondsPrefix = ''

      if (hours > 1 && minutes < 10) {
        minutesPrefix = '0'
      }
      if (seconds < 10) {
        secondsPrefix = '0'
      }

      const overtimePrefix = (this.progress > 100 ? 'Overtime: +' : '')

      return overtimePrefix + ((hours) ? `${hours}:` : '') + `${minutesPrefix + minutes}:${secondsPrefix + seconds}`
    }

    get finishedSince(): string | undefined {
      if (!this.finishedAt) {
        return
      }

      return formatDistanceToNow(new Date(this.finishedAt)) + ' ago'
    }

    get author(): TileAuthor | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.author
    }

    get showAuthor(): boolean {
      return this.author !== undefined && this.status === TileStatus.Failed
    }
  }
</script>

<style lang="scss">
  $tile-padding: 15px;
  $tile-author-height: 40px;
  $border-radius: 4px;

  .c-monitoror-tile {
    --tile-background: var(--color-unknown);
    position: relative;
    overflow: hidden;
    color: var(--color-text);
    background: var(--tile-background) linear-gradient(rgba(255, 255, 255, 0.1), transparent);
    border-radius: $border-radius;

    &__theme-dark {
      color: var(--tile-background);
      background: none;
      border: 5px solid var(--tile-background);

      &.c-monitoror-tile__status-queued,
      &.c-monitoror-tile__status-running {
        overflow: initial;
        border-bottom: none;
      }
    }

    &__empty {
      visibility: hidden;
    }

    &__status-succeeded {
      --tile-background: var(--color-succeeded);
    }

    &__status-failed {
      --tile-background: var(--color-failed);
    }

    &__status-warning {
      --tile-background: var(--color-warning);
    }

    &__status-canceled {
      --tile-background: var(--color-canceled);
    }

    &__status-action-required {
      --tile-background: var(--color-action-required);
    }
  }

  .c-monitoror-tile--content {
    height: 100%;
    padding: $tile-padding;
  }

  .c-monitoror-tile--label {
    font-size: 32px;
  }

  .c-monitoror-tile--message,
  .c-monitoror-tile--value {
    padding-top: 5px;
    font-size: 24px;
    opacity: 0.8;
  }

  .c-monitoror-tile__centered-message .c-monitoror-tile--message,
  .c-monitoror-tile--value {
    position: absolute;
    top: 50%;
    left: 50%;
    padding-top: 25px;
    text-align: center;
    font-size: 50px;
    transform: translate(-50%, -50%);
  }

  .c-monitoror-tile__group .c-monitoror-tile--message,
  .c-monitoror-tile--finished-at,
  .c-monitoror-tile--progress-time {
    position: absolute;
    right: $tile-padding;
    bottom: $tile-padding - 6px;
    font-size: 32px;
    opacity: 0.8;
    text-align: right;
    font-variant-numeric: tabular-nums;

    &::first-letter {
      text-transform: uppercase;
    }
  }

  .c-monitoror-tile--author {
    --tile-author-height: 40px;
    position: absolute;
    right: $tile-padding;
    bottom: $tile-padding + 35px;
    display: inline-block;
    padding: 3px 20px 3px 3px;
    max-width: calc(100% - 2 * #{$tile-padding});
    height: $tile-author-height;
    line-height: $tile-author-height - 6px;
    font-size: 20px;
    font-weight: normal;
    color: var(--tile-background);
    background: var(--color-background);
    border-radius: $tile-author-height;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    .c-monitoror-tile__theme-dark & {
      border: 1px solid var(--tile-background);
      line-height: $tile-author-height - 8px;
    }
  }

  .c-monitoror-tile--author-avatar {
    float: left;
    width: $tile-author-height - 6px;
    height: $tile-author-height - 6px;
    margin-right: 10px;
    background: #fff;
    border-radius: 100%;

    .c-monitoror-tile__theme-dark & {
      width: $tile-author-height - 8px;
      height: $tile-author-height - 8px;
    }
  }

  .c-monitoror-tile--progress-time {
    bottom: $tile-padding + 4px;
  }

  .c-monitoror-tile--progress {
    position: absolute;
    right: 0;
    left: 0;
    bottom: 0;
    border-top: 4px solid #2c3e50;
    background: var(--tile-background) linear-gradient(rgba(#2c3e50, 0.5), rgba(#2c3e50, 0.5));
    overflow: hidden;
    transform: translateZ(0); /* Optimize repaints */

    .c-monitoror-tile__theme-dark & {
      left: -5px;
      right: -5px;
      border-radius: 0 0 $border-radius $border-radius;
    }
  }

  .c-monitoror-tile--progress-bar {
    width: 100%;
    height: 10px;
    background: #fff;
    transform: translateX(-101%);
    transition: transform 150ms;

    .c-monitoror-tile__theme-dark & {
      height: 14px;
    }
  }

  .c-monitoror-tile__status-running .c-monitoror-tile--progress-bar {
    animation: progressBarBlink 3s linear infinite;
  }

  .c-monitoror-tile--progress-bar::after {
    content: "";
    display: block;
    width: 100px;
    height: 100%;
    position: absolute;
    top: 0;
    right: -50px;
    background: linear-gradient(to right, transparent, #fff 50%, var(--color-background) 50.01%, transparent);
    opacity: 0.2;
  }

  @keyframes progressBarBlink {
    0%,
    10%,
    90%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }

  .c-monitoror-tile--icon {
    position: absolute;
    bottom: $tile-padding;
    left: $tile-padding;
    opacity: 0.35;
  }

  .c-monitoror-tile__status-queued .c-monitoror-tile--icon,
  .c-monitoror-tile__status-running .c-monitoror-tile--icon {
    bottom: $tile-padding + 10px;
  }

  .c-monitoror-tile--sub-tiles {
    position: relative;
    overflow: hidden;
    height: calc((100vh / var(--rows)) * var(--row-span) - #{2 * $tile-padding + 85});
    margin-top: 7px;
    z-index: 1;
  }
</style>
