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
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor">
            <use :xlink:href="'./icons.svg#' + (mergeRequest ? 'merge-request' : 'branch')"/>
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

      <div class="c-monitoror-tile--value" v-if="displayedMetric">
        {{ displayedMetric }}
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
import useTileValues from '@/composables/useTileMetrics'
import {defineComponent, computed} from 'vue'
import {useStore} from 'vuex'

import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'
import useTileCommons from '@/composables/useTileCommons'
import TileStatus from '@/enums/tileStatus'
import TileType from '@/enums/tileType'
import TileConfig from '@/types/tileConfig'

import MonitororSubTile from '@/components/SubTile.vue'
import MonitororTileIcon from '@/components/TileIcon.vue'

export default defineComponent({
  name: 'MonitororTile',
  components: {
    MonitororSubTile,
    MonitororTileIcon,
  },
  props: {
    config: {
      type: Object as () => TileConfig,
      required: true
    }
  },
  setup: function (props) {
    const {
      state,
      theme,
      type,

      // status
      status,
      isQueued,
      isRunning,
      isSucceeded,
      isFailed,
      isWarning,

      // content
      label,

      // build
      build,
      branch,
      mergeRequest,
      mergeRequestLabelPrefix,
      progressTime,
      progressBarStyle,
      isOvertime,
      finishedSince,
      author,
      showAuthor,
    } = useTileCommons(props.config)

    const {
      displayedMetric,
    } = useTileValues(state)

    const store = useStore()

    const classes = computed((): Record<string, boolean | string> => {
      return {
        ['c-monitoror-tile__theme-' + theme.value]: true,
        'c-monitoror-tile__empty': isEmpty.value,
        'c-monitoror-tile__group': isGroup.value,
        'c-monitoror-tile__status-succeeded': isSucceeded.value,
        'c-monitoror-tile__status-failed': isFailed.value,
        'c-monitoror-tile__status-warning': isWarning.value,
        'c-monitoror-tile__status-running': isRunning.value,
        'c-monitoror-tile__status-queued': isQueued.value,
        'c-monitoror-tile__status-canceled': status.value === TileStatus.Canceled,
        'c-monitoror-tile__status-action-required': status.value === TileStatus.ActionRequired,
      }
    })

    const styles = computed((): Record<string, string | number> => {
      return {
        'grid-column': `auto / span ${columnSpan.value}`,
        'grid-row': `auto / span ${rowSpan.value}`,
        '--row-span': rowSpan.value,
      }
    })

    const isEmpty = computed((): boolean => {
      return type.value === TileType.Empty
    })

    const isGroup = computed((): boolean => {
      return type.value === TileType.Group
    })

    const columnSpan = computed((): number => {
      return props.config.columnSpan || 1
    })

    const rowSpan = computed((): number => {
      return props.config.rowSpan || 1
    })

    const displayedSubTiles = computed((): TileConfig[] | undefined => {
      if (!props.config.tiles) {
        return
      }

      const displayedSubTiles = props.config.tiles.filter((subTile) => {
        if (!Object.keys(store.state.tilesState).includes(subTile.stateKey)) {
          return false
        }

        return DISPLAYABLE_SUBTILE_STATUS.includes(store.state.tilesState[subTile.stateKey].status)
      })

      return displayedSubTiles
    })

    const message = computed((): string | undefined => {
      if (state.value === undefined) {
        return
      }

      return state.value.message
    })

    const buildId = computed((): string | undefined => {
      if (build.value === undefined) {
        return
      }

      return build.value.id
    })

    return {
      // attributes
      classes,
      styles,

      // type
      type,
      isEmpty,
      isGroup,

      // status
      status,
      isQueued,
      isRunning,
      isSucceeded,
      isFailed,
      isWarning,

      // layout
      columnSpan,
      rowSpan,

      // content
      label,
      message,
      displayedSubTiles,

      // build
      build,
      buildId,
      branch,
      mergeRequest,
      mergeRequestLabelPrefix,
      progressTime,
      progressBarStyle,
      isOvertime,
      finishedSince,
      author,
      showAuthor,

      // metrics
      displayedMetric,
    }
  },
})
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
    overflow-wrap: break-word;
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
    line-break: anywhere;

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
