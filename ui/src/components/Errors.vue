<template>
  <div class="c-monitoror-errors">
    <template v-if="!isOnline">
      I'm offline... Gimme my connection back!
    </template>
    <template v-else-if="hasErrors">
      <div class="c-monitoror-errors--config-info">
        <div class="c-monitoror-errors-title">
          We found {{errors.length}} error{{errors.length > 1 ? 's' : ''}} in this configuration:
        </div>
        <template v-if="configUrlOrPath !== 'undefined'">
          <code>{{configUrlOrPath}}</code> <br><br>
        </template>
        Last refresh at {{lastRefreshDate}}
      </div>
      <div class="c-monitoror-errors--error" v-for="error in errors">
        <hr>
        <!-- Blocking single-line errors -->
        <template v-if="error.id === ConfigErrorId.ConfigNotFound">
          <p class="c-monitoror-errors--error-title">
            Your configuration URL or path seems broken, please verify it.
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnableToParseConfig">
          <p class="c-monitoror-errors--error-title">
            Your configuration cannot be parsed, please verify it.
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.MissingPathOrUrl">
          <p class="c-monitoror-errors--error-title">
            Missing <code>configPath</code> or <code>configUrl</code> query param
          </p>
          <p>
            <a href="https://monitoror.com/documentation/#ui-configuration" target="_blank">
              Check UI configuration documentation
            </a>
          </p>
        </template>

        <!-- Config verify errors -->
        <template v-else-if="error.id === ConfigErrorId.MissingRequiredField">
          <p class="c-monitoror-errors--error-title">
            Missing required <code>{{error.data.fieldName}}</code> field.
          </p>
          <pre><code v-html="formatConfigExtract(error)"></code></pre>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnknownTileType">
          <p class="c-monitoror-errors--error-title">
            Unknown
            <code>{{JSON.parse(extractFieldValue(error.data.configExtract, error.data.fieldName))}}</code>
            tile type.
          </p>
          <pre><code v-html="formatConfigExtract(error)"></code></pre>
          <p>
            Expected <code>{{error.data.fieldName}}</code> value to be one of:
            <template v-for="expected in splitList(error.data.expected)">
              <code>{{expected}}</code>,
            </template>
          </p>
          <p>
            <a href="https://monitoror.com/documentation/#tile-definitions">Go to <em>Tile definitions</em> documentation section</a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnauthorizedSubtileType">
          <p class="c-monitoror-errors--error-title">
            Unauthorized
            <code>{{JSON.parse(extractFieldValue(error.data.configExtractHighlight, 'type'))}}</code>
            type as <code>GROUP</code> subtile.
          </p>
          <pre><code v-html="formatConfigExtract(error)"></code></pre>
        </template>
        <template v-else-if="error.id === ConfigErrorId.InvalidFieldValue">
          <p class="c-monitoror-errors--error-title">
            Invalid
            <code>{{error.data.fieldName}}</code>
            value.
          </p>
          <pre v-if="error.data.configExtract"><code v-html="formatConfigExtract(error)"></code></pre>
          <a
            v-if="getTileDocUrl(error) !== undefined"
            :href="getTileDocUrl(error)"
            target="_blank">
            Go to <code>{{JSON.parse(extractFieldValue(error.data.configExtract, 'type'))}}</code> documentation
          </a>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnsupportedVersion">
          <p class="c-monitoror-errors--error-title">
            Invalid
            <code>{{error.data.fieldName}}</code>
            value: <code>{{error.data.value}}</code>
          </p>
          <p>
            Supported config version by Core: <code>{{error.data.expected}}</code>
          </p>
          <p>
            Are you up to date?
            <a href="https://monitoror.com/documentation/#ui-configuration" target="_blank">
              Go check the latest config version on documentation
            </a>
          </p>
        </template>

        <!-- Misc errors -->
        <template v-else-if="error.id === ConfigErrorId.UnexpectedError">
          <p class="c-monitoror-errors--error-title">
            Unexpected error
          </p>
          {{error.message}}
        </template>
        <template v-else>
          <p class="c-monitoror-errors--error-title">
            {{error.message}}
          </p>
          {{error.id}}
        </template>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
  import {format} from 'date-fns'
  import {Component, Vue} from 'vue-property-decorator'

  import ConfigErrorId from '@/enums/configErrorId'
  import ConfigError from '@/interfaces/configError'

  @Component({})
  export default class MonitororErrors extends Vue {
    /**
     * Computed
     */

    get configUrlOrPath(): string {
      return this.$store.getters.configUrl || decodeURIComponent(this.$store.getters.configPath)
    }

    get lastRefreshDate(): string {
      return format(this.$store.state.lastRefreshDate, 'hh:mm:ss a')
    }

    get isOnline(): boolean {
      return this.$store.state.online
    }

    get errors(): ConfigError[] {
      return this.$store.state.errors
    }

    get hasErrors(): boolean {
      return this.errors.length > 0
    }

    get ConfigErrorId() {
      return ConfigErrorId
    }

    /*
     * Methods
     */

    public jsonSyntaxColor(jsonString: string): string {
      const coloredJson = jsonString
        .replace(/:\s+"(.*)"/g, ': <span class="code-string">"$1"</span>')
        .replace(/:\s+([.\d]+)/g, ': <span class="code-number">$1</span>')

      return coloredJson
    }

    public formatConfigExtract(configError: ConfigError): string {
      if (configError.data.configExtract === undefined) {
        return ''
      }

      const formattedConfigExtract = JSON.stringify(JSON.parse(configError.data.configExtract), null, 2)
      let html = formattedConfigExtract

      let configExtractHighlight = configError.data.configExtractHighlight
      let patternPrefix = ''

      if (configExtractHighlight === undefined && configError.data.fieldName !== undefined) {
        patternPrefix = `"${configError.data.fieldName}":\\s+`
        configExtractHighlight = this.extractFieldValue(configError.data.configExtract, configError.data.fieldName)
      }

      if (configExtractHighlight !== undefined) {
        const formattedConfigExtractHighlight = JSON.stringify(JSON.parse(configExtractHighlight), null, 2)
        const isHighlightMultiline = formattedConfigExtractHighlight.includes('\n')
        const multilinePrefix = isHighlightMultiline ? ' *' : ''
        const multilineSuffix = isHighlightMultiline ? ',?' : ''
        const pattern = multilinePrefix + patternPrefix + (formattedConfigExtractHighlight.replace(/\s+/g, '\\s+')) + multilineSuffix
        const find = new RegExp(pattern)
        const matches = formattedConfigExtract.match(find)

        if (matches === null) {
          return html
        }

        const match = matches[0]
        const markClassAttr = isHighlightMultiline ? 'class="multiline-mark"' : ''

        html = formattedConfigExtract.replace(match, `<mark ${markClassAttr}>${match}</mark>`)
      }

      if (html.includes('</mark>')) {
        html = `<span class="has-mark">${html}</span>`
      }

      return this.jsonSyntaxColor(html)
    }

    public splitList(list: string): string[] {
      try {
        return list.split(',').map((item) => item.trim()).sort()
      } catch (e) {
        return [list]
      }
    }

    public extractFieldValue(jsonString: string, fieldName: string): string | undefined {
      try {
        return JSON.stringify(JSON.parse(jsonString)[fieldName])
      } catch (e) {
        return
      }
    }

    public getTileDocUrl(error: ConfigError): string | undefined {
      const tileType = this.extractFieldValue(error.data.configExtract as string, 'type')

      if (tileType === undefined) {
        return
      }

      const url = 'https://monitoror.com/documentation/#tile-' + JSON.parse(tileType).toLowerCase()

      return url
    }
  }
</script>

<style lang="scss">
  .c-monitoror-errors {
    a {
      color: var(--color-succeeded);
      font-size: 20px;

      code {
        color: inherit;
      }
    }

    hr {
      width: 150px;
      border-color: var(--color-logo-background);
      margin: 25px auto;
      transition: border-color 300ms;
    }

    code {
      display: inline-block;
      background: var(--color-code-background);
      margin-top: 5px;
      padding: 3px 7px;
      border-radius: 4px;
    }

    pre {
      position: relative;
      margin-top: 50px;
    }

    pre code {
      color: var(--color-spring-wood);
    }

    pre::before {
      content: "";
      position: absolute;
      top: 0;
      bottom: 0;
      left: -9999px;
      right: 0;
      background: var(--color-code-background);
      box-shadow: 9999px 0 0 var(--color-code-background);
      z-index: -1;
    }

    pre::after {
      content: "Config extract";
      position: absolute;
      top: -1em;
      left: 0;
      font-size: 20px;
      text-transform: uppercase;
      letter-spacing: 1px;
      color: var(--color-unknown);
    }

    .code-string {
      color: var(--color-succeeded);
    }

    .code-number {
      color: var(--color-warning);
    }

    .has-mark {
      color: rgba(255, 255, 255, 0.3);

      .code-string,
      .code-number {
        opacity: 0.3;
      }
    }

    mark {
      display: inline-block;
      color: var(--color-spring-wood);
      background: var(--color-docs-background);
      border-radius: 4px;
      padding: 5px 7px;
      margin: 0 -7px;
      box-shadow: 0 0 5px rgba(0, 0, 0, 0.2), 0 0 25px rgba(0, 0, 0, 0.1);

      &.multiline-mark {
        padding-right: 50px;
      }

      .code-string,
      .code-number {
        opacity: 1 !important;
      }
    }
  }

  .c-monitoror-errors--error {
    margin-bottom: 50px;
    font-size: 18px;
  }

  .c-monitoror-errors--error-title {
    font-size: 22px;
    color: #ffffff;
  }

  .c-monitoror-errors--error code {
    color: var(--color-spring-wood);
  }

  .c-monitoror-errors--error pre code {
    padding: 20px;
  }

  .c-monitoror-errors--config-info {
    text-align: center;

    code {
      color: var(--color-failed);
      padding: 7px 13px;
    }
  }
</style>
