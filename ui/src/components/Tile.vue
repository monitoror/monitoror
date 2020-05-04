<template>
  <div class="c-monitoror-tile" :class="classes" :style="styles">
    <div class="c-monitoror-tile--content" v-if="!isEmpty">
      <div class="c-monitoror-tile--label">
        <template v-if="mergeRequestLabelPrefix">{{ mergeRequestLabelPrefix }}</template>
        <template v-if="mergeRequestLabelPrefix && label"> @</template>
        {{ label }}
      </div>

      <div class="c-monitoror-tile--build-info" v-if="branch || buildId">
        <template v-if="branch">
          <svg v-if="mergeRequest" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor">
            <!-- Merge request icon -->
            <path fill-rule="evenodd" d="M7.177 3.073L9.573.677A.25.25 0 0110 .854v4.792a.25.25 0 01-.427.177L7.177 3.427a.25.25 0 010-.354zM3.75 2.5a.75.75 0 100 1.5.75.75 0 000-1.5zm-2.25.75a2.25 2.25 0 113 2.122v5.256a2.251 2.251 0 11-1.5 0V5.372A2.25 2.25 0 011.5 3.25zM11 2.5h-1V4h1a1 1 0 011 1v5.628a2.251 2.251 0 101.5 0V5A2.5 2.5 0 0011 2.5zm1 10.25a.75.75 0 111.5 0 .75.75 0 01-1.5 0zM3.75 12a.75.75 0 100 1.5.75.75 0 000-1.5z"/>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor">
            <!-- Branch icon -->
            <path fill-rule="evenodd" d="M11.75 2.5a.75.75 0 100 1.5.75.75 0 000-1.5zm-2.25.75a2.25 2.25 0 113 2.122V6A2.5 2.5 0 0110 8.5H6a1 1 0 00-1 1v1.128a2.251 2.251 0 11-1.5 0V5.372a2.25 2.25 0 111.5 0v1.836A2.492 2.492 0 016 7h4a1 1 0 001-1v-.628A2.25 2.25 0 019.5 3.25zM4.25 12a.75.75 0 100 1.5.75.75 0 000-1.5zM3.5 3.25a.75.75 0 111.5 0 .75.75 0 01-1.5 0z"/>
          </svg>
          {{ branch }}
        </template>
        <template v-if="branch && buildId"> —</template>
        <template v-if="buildId">
          #{{ buildId }}
        </template>
      </div>

      <div class="c-monitoror-tile--message" v-if="message">
        {{ message }}
      </div>

      <div class="c-monitoror-tile--value" v-if="displayedValue">
        {{ displayedValue }}
      </div>

      <div class="c-monitoror-tile--sub-tiles" v-if="isGroup">
        <monitoror-sub-tile v-for="subTile in displayedSubTiles" :key="subTile.stateKey"
                            :config="subTile"></monitoror-sub-tile>
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
  import {Component} from 'vue-property-decorator'

  import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'
  import TileStatus from '@/enums/tileStatus'
  import TileType from '@/enums/tileType'
  import TileValueUnit from '@/enums/tileValueUnit'
  import TileConfig from '@/interfaces/tileConfig'

  import AbstractMonitororTile from '@/classes/monitororTile'
  import MonitororSubTile from '@/components/SubTile.vue'
  import MonitororTileIcon from '@/components/TileIcon.vue'
  import TileValue from '@/interfaces/tileValue'

  @Component({
    components: {
      MonitororSubTile,
      MonitororTileIcon,
    },
  })
  export default class MonitororTile extends AbstractMonitororTile {

    /*
     * Computed
     */

    get classes() {
      return {
        ['c-monitoror-tile__theme-' + this.theme]: true,
        'c-monitoror-tile__empty': this.isEmpty,
        'c-monitoror-tile__group': this.isGroup,
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

    get isEmpty(): boolean {
      return this.type === TileType.Empty
    }

    get isGroup(): boolean {
      return this.type === TileType.Group
    }

    get columnSpan(): number {
      return this.config.columnSpan || 1
    }

    get rowSpan(): number {
      return this.config.rowSpan || 1
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

    get value(): TileValue | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.value
    }

    get message(): string | undefined {
      if (this.state === undefined) {
        return
      }

      return this.state.message
    }

    get buildId(): string | undefined {
      if (this.build === undefined) {
        return
      }

      return this.build.id
    }

    get unit(): TileValueUnit {
      if (this.value === undefined) {
        return TileValueUnit.Raw
      }

      return this.value.unit as TileValueUnit
    }

    get values(): string[] | undefined {
      if (this.value === undefined) {
        return
      }

      return this.value.values
    }

    get displayedValue(): string | undefined {
      if (this.values === undefined) {
        return
      }

      const UNIT_DISPLAY = {
        [TileValueUnit.Millisecond]: 'ms',
        [TileValueUnit.Ratio]: '%',
        [TileValueUnit.Number]: '',
        [TileValueUnit.Raw]: '',
      }

      let value = this.values[this.values.length - 1]
      if (this.unit === TileValueUnit.Millisecond) {
        value = Math.round(parseFloat(value)).toString()
      } else if (this.unit === TileValueUnit.Ratio) {
        value = (parseFloat(value) * 100).toFixed(2).toString()
      }

      return value + UNIT_DISPLAY[this.unit]
    }
  }
</script>

<style lang="scss">
  $tile-author-height: 40px;
  $border-radius: 4px;

  .c-monitoror-tile {
    --tile-background: var(--color-unknown);
    --tile-padding: 15px;

    position: relative;
    overflow: hidden;
    color: var(--color-text);
    background: var(--tile-background) linear-gradient(rgba(255, 255, 255, 0.1), transparent);
    border-radius: $border-radius;

    &__theme-dark {
      --tile-padding: 18px;

      color: var(--tile-background);
      background: none;
      box-shadow:
        inset 5px 0 0 var(--tile-background),
        inset -5px 0 0 var(--tile-background),
        inset 0 5px 0 var(--tile-background),
        inset 0 -5px 0 var(--tile-background);
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
    padding: var(--tile-padding);
    zoom: var(--zoom);

    @media screen and (max-width: 750px) {
      min-height: 160px;
    }
  }

  .c-monitoror-tile--label {
    font-size: 32px;
    line-height: 1.2;
    font-weight: bold;
    margin-bottom: 3px;
  }

  .c-monitoror-tile--message,
  .c-monitoror-tile--value {
    padding-top: 5px;
    font-size: 24px;
    opacity: 0.8;
  }

  .c-monitoror-tile--build-info {
    font-size: 24px;
    font-family: 'JetBrains Mono', monospace;
    opacity: 0.8;

    .c-monitoror-tile__theme-dark & {
      opacity: 0.7;
    }
  }

  .c-monitoror-tile--build-info svg {
    display: inline-block;
    width: 18px;
    vertical-align: middle;
    transform: translate(1px, 0px);
    margin-right: -7px;
  }

  .c-monitoror-tile--value {
    position: absolute;
    top: 50%;
    left: 50%;
    padding-top: 25px;
    text-align: center;
    font-size: 50px;
    font-weight: bold;
    transform: translate(-50%, -50%);
  }

  .c-monitoror-tile__group .c-monitoror-tile--message,
  .c-monitoror-tile--finished-at,
  .c-monitoror-tile--progress-time {
    position: absolute;
    right: var(--tile-padding);
    bottom: calc(var(--tile-padding) - 6px);
    font-size: 30px;
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
    right: var(--tile-padding);
    bottom: calc(var(--tile-padding) + 35px);
    display: inline-block;
    padding: 3px 20px 3px 3px;
    max-width: calc(100% - 2 * var(--tile-padding));
    height: $tile-author-height;
    line-height: $tile-author-height - 6px;
    font-size: 20px;
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
    bottom: calc(var(--tile-padding) + 4px);
  }

  .c-monitoror-tile--progress {
    position: absolute;
    right: 0;
    left: 0;
    bottom: 0;
    border-top: 4px solid var(--color-background);
    background: var(--tile-background) linear-gradient(rgba(#2c3e50, 0.5), rgba(#2c3e50, 0.5));
    overflow: hidden;
    transform: translateZ(0); /* Optimize repaints */
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
    bottom: var(--tile-padding);
    left: var(--tile-padding);
    opacity: 0.35;
    width: 40px;
    height: 40px;
  }

  .c-monitoror-tile--icon svg {
    position: absolute;
    bottom: 0;
    left: 0;
  }

  .c-monitoror-tile__status-queued .c-monitoror-tile--icon,
  .c-monitoror-tile__status-running .c-monitoror-tile--icon {
    bottom: calc(var(--tile-padding) + 10px);
  }

  .c-monitoror-tile--sub-tiles {
    position: relative;
    margin-top: 10px;
    z-index: 1;
  }
</style>
