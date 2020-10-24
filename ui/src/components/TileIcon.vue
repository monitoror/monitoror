<template>
  <div class="c-monitoror-tile-icon">
    <svg viewBox="0 0 64 64">
      <use :xlink:href="'./icons.svg#' + iconId"></use>
    </svg>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue'

import TileIconId from '@/enums/tileIconId'
import TileType from '@/enums/tileType'

export default defineComponent({
  name: 'MonitororTileIcon',
  props: {
    tileType: {
      type: String as () => TileType,
      required: true,
    },
  },
  setup(props) {
    const iconId = computed(() => {
      switch (props.tileType) {
        case TileType.AzureDevOpsBuild:
        case TileType.AzureDevOpsRelease:
          return TileIconId.Azure

        case TileType.GitHubChecks:
        case TileType.GitHubPullRequest:
        case TileType.GitHubCount:
          return TileIconId.GitHub

        case TileType.GitLabIssues:
        case TileType.GitLabMergeRequest:
        case TileType.GitLabPipeline:
          return TileIconId.GitLab

        case TileType.Group:
          return TileIconId.Group

        case TileType.JenkinsBuild:
          return TileIconId.Jenkins

        case TileType.HttpStatus:
        case TileType.HttpRaw:
        case TileType.HttpFormatted:
          return TileIconId.Http

        case TileType.Ping:
          return TileIconId.Ping

        case TileType.PingdomCheck:
        case TileType.PingdomTransactionCheck:
          return TileIconId.Pingdom

        case TileType.Port:
          return TileIconId.Port

        case TileType.TravisCiBuild:
          return TileIconId.TravisCi
      }
    })

    return {
      iconId,
    }
  }
})
</script>

<style lang="scss">
  .c-monitoror-tile-icon svg {
    fill: currentColor;
  }
</style>
