<template>
  <div class="c-monitoror-sub-tile" :class="classes">
    <div class="c-monitoror-sub-tile--content">
      <div class="c-monitoror-sub-tile--label">
        {{ label }}
      </div>

      <div class="c-monitoror-sub-tile--status">
        <monitoror-tile-icon :tile-type="type" class="c-monitoror-sub-tile--icon"></monitoror-tile-icon>
      </div>

      <template v-if="isRunning || isQueued">
        <div class="c-monitoror-sub-tile--progress-time">
          <template v-if="isRunning">
            {{ progressTime }}
          </template>
          <template v-else>
            &bull;&bull;&bull;
          </template>
        </div>
        <div class="c-monitoror-sub-tile--progress">
          <div class="c-monitoror-sub-tile--progress-bar" :style="progressBatStyle"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  import Vue from 'vue'
  import {Component, Prop} from 'vue-property-decorator'

  import MonitororTileIcon from '@/components/TileIcon.vue'
  import {TileConfig, TileState, TileStatus, TileType} from '@/store'

  @Component({
    components: {
      MonitororTileIcon,
    },
  })
  export default class MonitororSubTile extends Vue {
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
        'c-monitoror-sub-tile__status-succeeded': [this.previousStatus, this.status].includes(TileStatus.Success),
        'c-monitoror-sub-tile__status-failed': [this.previousStatus, this.status].includes(TileStatus.Failed),
        'c-monitoror-sub-tile__status-warning': [this.previousStatus, this.status].includes(TileStatus.Warning),
        'c-monitoror-sub-tile__status-running': this.isRunning,
        'c-monitoror-sub-tile__status-queued': this.isQueued,
        'c-monitoror-sub-tile__status-canceled': this.status === TileStatus.Canceled,
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

    get label(): string | undefined {
      if (this.config.label) {
        return this.config.label
      }

      if (this.state) {
        return this.state.label
      }
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
      if (!this.duration || this.estimatedDuration === undefined) {
        return
      }

      const progress = this.duration / this.estimatedDuration * 100

      return progress
    }

    get progressTime(): string | undefined {
      if (!this.progress || this.estimatedDuration === undefined || !this.duration) {
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

      const overtimePrefix = (this.progress > 100 ? '+' : '')

      return overtimePrefix + ((hours) ? `${hours}:` : '') + `${minutesPrefix + minutes}:${secondsPrefix + seconds}`
    }
  }
</script>

<style lang="scss">
  $tile-padding: 15px;
  $border-radius: 4px;

  .c-monitoror-sub-tile {
    --sub-tile-status-color: var(--color-unknown);
    display: inline-block;
    position: relative;
    margin-right: 7px;
    font-size: 22px;
    font-weight: normal;

    &__status-succeeded {
      --sub-tile-status-color: var(--color-succeeded);
    }

    &__status-failed {
      --sub-tile-status-color: var(--color-failed);
    }

    &__status-warning {
      --sub-tile-status-color: var(--color-warning);
    }

    &__status-cancel {
      --sub-tile-status-color: var(--color-warning); // TODO: yellow
    }
  }

  .c-monitoror-sub-tile--content {
    position: relative;
    display: inline-block;
    padding: 2px 15px 2px 35px;
    color: var(--sub-tile-status-color);
    background: var(--color-background);
    border-radius: 15px;
    overflow: hidden;
  }

  .c-monitoror-sub-tile--label {
    display: inline-block;
  }

  .c-monitoror-sub-tile__status-queued .c-monitoror-sub-tile--label,
  .c-monitoror-sub-tile__status-running .c-monitoror-sub-tile--label {
    max-width: calc(100vw / var(--columns) - 153px);
    vertical-align: bottom;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .c-monitoror-sub-tile--progress-time {
    display: inline-block;
    opacity: 0.8;
    text-align: right;
    font-variant-numeric: tabular-nums;
    padding-left: 15px;

    &::first-letter {
      text-transform: uppercase;
    }
  }

  .c-monitoror-sub-tile--progress-time {
    bottom: $tile-padding + 4px;
  }

  .c-monitoror-sub-tile--progress {
    position: absolute;
    right: 5px;
    left: 5px;
    bottom: 0;
    height: 2px;
    overflow: hidden;
    transform: translateZ(0); /* Optimize repaints */
  }

  .c-monitoror-sub-tile--progress-bar {
    width: 100%;
    height: 100%;
    background: #fff;
    transform: translateX(-101%);
    transition: transform 150ms;
    animation: progressBarBlink 3s linear infinite;
  }

  .c-monitoror-sub-tile--progress-bar::after {
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

  .c-monitoror-sub-tile--status {
    position: absolute;
    top: 3px;
    left: 10px;
    width: 20px;
    height: 20px;
    border-radius: 100%;
  }

  .c-monitoror-sub-tile--icon {
    width: 100%;
    height: 100%;
  }
</style>
