<template>
  <div class="c-monitoror-sub-tile" :class="classes">
    <div class="c-monitoror-sub-tile--content">
      <div class="c-monitoror-sub-tile--label">
        <template v-if="mergeRequestLabelPrefix">{{ mergeRequestLabelPrefix }}</template>
        <template v-else-if="branch">{{ branch }}</template>
        <template v-if="(mergeRequestLabelPrefix || branch) && label"> @ </template>
        {{ label }}
      </div>

      <monitoror-tile-icon :tile-type="type" class="c-monitoror-sub-tile--icon"></monitoror-tile-icon>

      <template v-if="isRunning || isQueued">
        <div class="c-monitoror-sub-tile--progress-time" :class="{'c-monitoror-sub-tile--progress-time__overtime': isOvertime}">
          <template v-if="isRunning">
            {{ progressTime }}
          </template>
          <template v-else>
            &bull;&bull;&bull;
          </template>
        </div>
        <div class="c-monitoror-sub-tile--progress">
          <div class="c-monitoror-sub-tile--progress-bar" :style="progressBarStyle"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  import {Component} from 'vue-property-decorator'

  import TileStatus from '@/enums/tileStatus'

  import AbstractMonitororTile from '@/classes/monitororTile'
  import MonitororTileIcon from '@/components/TileIcon.vue'

  @Component({
    components: {
      MonitororTileIcon,
    },
  })
  export default class MonitororSubTile extends AbstractMonitororTile {
    /*
     * Computed
     */

    get classes() {
      return {
        ['c-monitoror-sub-tile__theme-' + this.theme]: true,
        'c-monitoror-sub-tile__status-succeeded': this.isSucceeded,
        'c-monitoror-sub-tile__status-failed': this.isFailed,
        'c-monitoror-sub-tile__status-warning': this.isWarning,
        'c-monitoror-sub-tile__status-running': this.isRunning,
        'c-monitoror-sub-tile__status-queued': this.isQueued,
        'c-monitoror-sub-tile__status-canceled': this.status === TileStatus.Canceled,
        'c-monitoror-sub-tile__status-action-required': this.status === TileStatus.ActionRequired,
      }
    }

    get progressTime(): string | undefined {
      if (super.progressTime === undefined) {
        return
      }

      return super.progressTime.replace('Overtime: ', '')
    }
  }
</script>

<style lang="scss">
  $tile-padding: 15px;
  $border-radius: 17px;

  .c-monitoror-sub-tile {
    --sub-tile-status-color: var(--color-unknown);
    display: inline-block;
    position: relative;
    margin-right: 7px;
    font-size: 24px;
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

    &__status-canceled {
      --sub-tile-status-color: var(--color-canceled);
    }

    &__status-action-required {
      --sub-tile-status-color: var(--color-action-required);
    }
  }

  .c-monitoror-sub-tile--content {
    position: relative;
    display: inline-block;
    padding: 2px 15px 2px 40px;
    color: var(--sub-tile-status-color);
    background: var(--color-background);
    border: 1px solid var(--color-background);
    border-radius: $border-radius;
    overflow: hidden;

    .c-monitoror-sub-tile__theme-dark & {
      border-color: var(--sub-tile-status-color);
    }
  }

  .c-monitoror-sub-tile--label {
    display: inline-block;
    font-weight: 600;
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
    bottom: $tile-padding + 4px;

    &::first-letter {
      text-transform: uppercase;
    }
  }

  .c-monitoror-sub-tile--progress-time__overtime {
    font-weight: 600;
    opacity: 1;
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
      opacity: 0.6;
    }
  }

  .c-monitoror-sub-tile--icon {
    position: absolute;
    top: 3px;
    left: 10px;
    width: 23px;
    height: 23px;
    border-radius: 100%;
  }
</style>
