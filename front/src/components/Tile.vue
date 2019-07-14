<template>
  <div class="c-monitoror-tile" :class="classes" :style="styles">
    <div class="c-monitoror-tile--content" v-if="!isEmpty">
      <div class="c-monitoror-tile--label">
        {{ label }}
      </div>

      <div class="c-monitoror-tile--message">
        {{ message }}
      </div>

      <div class="c-monitoror-tile--finished-at" v-if="finishedSince">
        {{ finishedSince }}
      </div>

      <template v-if="isRunning || isQueued">
        <div class="c-monitoror-tile--progress-time">
          {{ progressTime || 'Pending...' }}
        </div>
        <div class="c-monitoror-tile--progress">
          <div class="c-monitoror-tile--progress-bar" :style="progressBatStyle"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  import distanceInWordsToNow from 'date-fns/distance_in_words_to_now'
  import Vue from 'vue'
  import {Component, Prop} from 'vue-property-decorator'

  import {TileCategory, TileConfig, TileState, TileStatus, TileType} from '@/store'

  @Component
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
        'c-monitoror-tile__empty': this.isEmpty,
        'c-monitoror-tile__group': this.isGroup,
        'c-monitoror-tile__centered-message': this.category && [TileCategory.Health].includes(this.category),
        'c-monitoror-tile__status-succeeded': [this.previousStatus, this.status].includes(TileStatus.Success),
        'c-monitoror-tile__status-failed': [this.previousStatus, this.status].includes(TileStatus.Failed),
        'c-monitoror-tile__status-warning': [this.previousStatus, this.status].includes(TileStatus.Warning),
        'c-monitoror-tile__status-running': this.isRunning,
        'c-monitoror-tile__status-queued': this.isQueued,
        'c-monitoror-tile__status-canceled': this.status === TileStatus.Canceled,
      }
    }

    get styles() {
      return {
        'grid-columns': `auto / span ${this.columnSpan}`,
        'grid-rows': `auto / span ${this.rowSpan}`,
      }
    }

    get progressBatStyle() {
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

    get label(): string | undefined {
      if (this.config.label) {
        return this.config.label
      }

      if (this.state) {
        return this.state.label
      }
    }

    get columnSpan(): number {
      return this.config.columnSpan || 0
    }
    get rowSpan(): number {
      return this.config.rowSpan || 0
    }

    get url(): string | undefined {
      return this.config.url
    }

    get stateKey(): string {
      return this.config.stateKey
    }

    get state(): TileState | undefined {
      if (!this.$store.state.tilesState.hasOwnProperty(this.stateKey)) {
        return
      }

      return this.$store.state.tilesState[this.stateKey]
    }

    get category(): TileCategory | undefined {
      if (!this.state) {
        return
      }

      return this.state.category
    }

    get status(): string | undefined {
      if (!this.state) {
        return
      }

      return this.state.status
    }

    get isRunning(): boolean {
      return this.status === TileStatus.Running
    }

    get isQueued(): boolean {
      return this.status === TileStatus.Queued
    }

    get previousStatus(): string | undefined {
      if (!this.state) {
        return
      }

      return this.state.previousStatus
    }

    get message(): string | undefined {
      if (!this.state) {
        return
      }

      return this.state.message
    }

    get startedAt(): number | undefined {
      if (!this.state) {
        return
      }

      return this.state.startedAt
    }

    get finishedAt(): number | undefined {
      if (!this.state) {
        return
      }

      return this.state.finishedAt
    }

    get duration(): number | undefined {
      if (!this.state) {
        return
      }

      return this.state.duration
    }

    get estimatedDuration(): number | undefined {
      if (!this.state) {
        return
      }

      return this.state.estimatedDuration
    }

    get progress(): number | undefined {
      if (!this.duration || !this.estimatedDuration) {
        return
      }

      const progress = this.duration / this.estimatedDuration * 100

      return progress
    }

    get progressTime(): string | undefined {
      if (!this.progress || !this.estimatedDuration || !this.duration) {
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

      return distanceInWordsToNow(new Date(this.finishedAt * 1000)) + ' ago'
    }
  }
</script>

<style lang="scss">
  $tile-padding: 15px;
  $border-radius: 4px;

  .c-monitoror-tile {
    position: relative;
    overflow: hidden;
    color: var(--color-text);
    background: var(--color-unknown) linear-gradient(rgba(255, 255, 255, 0.1), transparent);
    border-radius: $border-radius;

    &__empty {
      opacity: 0;
    }

    &__status-succeeded {
      background-color: var(--color-succeeded);
    }

    &__status-failed {
      background-color: var(--color-failed);
    }

    &__status-warning {
      background-color: var(--color-warning);
    }
  }

  .c-monitoror-tile--content {
    height: 100%;
    padding: $tile-padding;
  }

  .c-monitoror-tile--label {
    font-size: 32px;
  }

  .c-monitoror-tile--message {
    padding-top: 5px;
    font-size: 24px;
    opacity: 0.8;
  }

  .c-monitoror-tile__centered-message .c-monitoror-tile--message {
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

  .c-monitoror-tile--progress-time {
    bottom: $tile-padding + 4px;
  }

  .c-monitoror-tile--progress {
    position: absolute;
    right: 0;
    left: 0;
    bottom: 0;
    border-top: 4px solid #2c3e50;
    background: rgba(#2c3e50, 0.5);
    overflow: hidden;
    transform: translateZ(0); /* Optimize repaints */
  }

  .c-monitoror-tile--progress-bar {
    width: 100%;
    height: 10px;
    background: #fff;
    transform: translateX(-101%);
    transition: transform 150ms;
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
</style>
