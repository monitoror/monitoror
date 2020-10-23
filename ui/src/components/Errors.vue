<template>
  <div class="c-monitoror-errors" :class="classes">
    <template v-if="!isOnline">
      I'm offline... Gimme my connection back!
    </template>
    <template v-else-if="hasErrors">
      <div class="c-monitoror-errors--config-info">
        <div class="c-monitoror-errors-title" v-if="hasConfigVerifyErrors">
          We found {{errors.length}} error{{errors.length > 1 ? 's' : ''}}
          <template v-if="configParam !== undefined">
            in this configuration:
          </template>
        </div>
        <template v-if="configParam !== undefined && !isUsingDefaultConfig">
          <code>{{configParam}}</code> <br><br>
        </template>
        <template v-if="!isCoreDown">
          Last refresh at {{lastRefreshDate}}

          <hr v-if="!hasConfigVerifyErrors">
        </template>
      </div>
      <div class="c-monitoror-errors--error" v-for="error in errors">
        <!-- Blocking single-line errors -->
        <template v-if="isCoreDown">
          <p class="c-monitoror-errors--error-title">
            Monitoror Core seems down
          </p>
          <p>
            Is there a Monitoror Core running at <code>{{apiBaseUrl}}</code>?
          </p>
          <hr>
          <p class="c-monitoror-errors--last-refresh">
            Last refresh at {{lastRefreshDate}}
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.ConfigNotFound">
          <p class="c-monitoror-errors--error-title">
            Your configuration URL or path seems broken, please verify it
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnknownNamedConfig">
          <p class="c-monitoror-errors--error-title">
            <code>{{error.data.value}}</code> named configuration cannot be found
          </p>
          <p v-if="error.data.expected">
            Must be one of:
            <template v-for="(namedConfig, index) in splitList(error.data.expected)">
              <code>{{namedConfig}}</code><template v-if="index !== splitList(error.data.expected).length - 1">, </template>
            </template>
          </p>
          <p v-else>
            You must give a UI configuration via <code>config</code> query parameter in order to use Monitoror
          </p>
        </template>

        <!-- Config verify errors -->
        <template v-else-if="error.id === ConfigErrorId.FieldTypeMismatch">
          <p class="c-monitoror-errors--error-title">
            Type mismatch in
            <code>{{error.data.fieldName}}</code>
            field
          </p>
          <p>
            Expected type: <code>{{error.data.expected}}</code>
          </p>
          <pre v-if="error.data.configExtract" :data-config-param="configParam"><code
            v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/" target="_blank">
              Go to documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.InvalidEscapedCharacter">
          <p class="c-monitoror-errors--error-title">
            Invalid escaped character
            <code>{{error.data.configExtractHighlight}}</code>
          </p>
          <p>
            Try to escape the backslash: <code>\{{error.data.configExtractHighlight}}</code>
          </p>
          <pre v-if="error.data.configExtract" :data-config-param="configParam"><code
            v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/" target="_blank">
              Go to documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.InvalidFieldValue">
          <p class="c-monitoror-errors--error-title">
            Invalid
            <code>{{error.data.fieldName}}</code>
            value
          </p>
          <pre v-if="error.data.configExtract" :data-config-param="configParam"><code
            v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation" v-if="getTileDocUrl(error) !== undefined">
            <a :href="getTileDocUrl(error)" target="_blank">
              Go to <strong>{{parsedExtractFieldValue(error.data.configExtract, 'type')}}</strong> documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.MissingRequiredField">
          <p class="c-monitoror-errors--error-title">
            Missing required <code>{{error.data.fieldName}}</code> field
          </p>
          <pre :data-config-param="configParam"><code v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/" target="_blank">
              Go to documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnableToParseConfig">
          <p class="c-monitoror-errors--error-title">
            Your configuration cannot be parsed, please verify it
          </p>
          <p>
            {{error.message}}
          </p>
          <pre :data-config-param="configParam"><code v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/" target="_blank">
              Go to documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnknownField">
          <p class="c-monitoror-errors--error-title">
            Unknown
            <code>{{error.data.fieldName}}</code>
            field name
          </p>
          <p
            v-if="error.data.expected && guessExpectedFieldName(error.data.fieldName, splitList(error.data.expected)) !== undefined">
            Did you mean
            <code>{{guessExpectedFieldName(error.data.fieldName, splitList(error.data.expected))}}</code>?
          </p>
          <pre :data-config-param="configParam"><code v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/" target="_blank">
              Go to documentation
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnknownTileType">
          <p class="c-monitoror-errors--error-title">
            Unknown
            <code>{{parsedExtractFieldValue(error.data.configExtract, error.data.fieldName)}}</code>
            tile type
          </p>
          <p
            v-if="guessExpectedValue(error.data.configExtract, error.data.fieldName, splitList(error.data.expected)) !== undefined">
            Did you mean
            <code>{{guessExpectedValue(error.data.configExtract, error.data.fieldName,
              splitList(error.data.expected))}}</code>?
          </p>
          <pre :data-config-param="configParam"><code v-html="formatConfigExtract(error)"></code></pre>
          <p class="go-to-documentation">
            <a href="https://monitoror.com/documentation/#tile-definitions" target="_blank">
              Go to <strong>Tile definitions</strong> documentation section
            </a>
          </p>
        </template>
        <template v-else-if="error.id === ConfigErrorId.UnauthorizedSubtileType">
          <p class="c-monitoror-errors--error-title">
            Unauthorized
            <code>{{parsedExtractFieldValue(error.data.configExtractHighlight, 'type')}}</code>
            type as <code>GROUP</code> subtile
          </p>
          <pre :data-config-param="configParam"><code
            v-html="ellipsisUnnecessaryParams(formatConfigExtract(error))"></code></pre>
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
          <p class="go-to-documentation">
            Are you up to date? <br>
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
          <p>
            {{error.message}}
          </p>
        </template>
        <template v-else>
          <p class="c-monitoror-errors--error-title">
            {{error.message}}
          </p>
          <p>
            {{error.id}}
          </p>
          <pre :data-config-param="configParam"><code
            v-html="formatConfigExtract(error)"></code></pre>
        </template>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
  import {format} from 'date-fns'
  import {nextTick} from 'vue'
  import {Options, Vue} from 'vue-class-component'

  import ConfigErrorId from '@/enums/configErrorId'
  import ellipsisUnnecessaryParams from '@/helpers/ellipsisUnnecessaryParams'
  import formatConfigExtract from '@/helpers/formatConfigExtract'
  import getTileDocUrl from '@/helpers/getTileDoc'
  import guessExpectedFieldName from '@/helpers/guessExpectedFieldName'
  import guessExpectedValue from '@/helpers/guessExpectedValue'
  import hasConfigVerifyErrors from '@/helpers/hasConfigVerifyErrors'
  import parsedExtractFieldValue from '@/helpers/parsedExpectedValue'
  import splitList from '@/helpers/splitList'
  import ConfigError from '@/interfaces/configError'
  import {DEFAULT_CONFIG_NAME} from '@/store'

  @Options<MonitororErrors>({
    watch: {
      errors: async () => {
        await nextTick()
        Array.from(document.querySelectorAll('.has-mark')).forEach((errorConfigExtract: Element) => {
          const pre = errorConfigExtract.parentNode?.parentNode as HTMLElement
          const mark = pre.querySelector('mark')

          if (mark === null) {
            return
          }

          pre.scrollTop = mark.offsetTop - pre.offsetTop
        })
      },
    },
  })
  export default class MonitororErrors extends Vue {
    /**
     * Computed
     */

    get classes() {
      return {
        'c-monitoror-errors__config-verify-errors': this.hasConfigVerifyErrors,
      }
    }

    get configParam(): string {
      return this.$store.getters.configParam
    }

    get isUsingDefaultConfig(): boolean {
      return this.configParam === DEFAULT_CONFIG_NAME
    }

    get isCoreDown(): boolean {
      if (!this.hasErrors) {
        return false
      }

      return this.errors[0].id === ConfigErrorId.CannotBeFetched
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

    get hasConfigVerifyErrors(): boolean {
      return hasConfigVerifyErrors(this.errors)
    }

    get apiBaseUrl(): string {
      return this.$store.getters.apiBaseUrl
    }

    get ConfigErrorId() {
      return ConfigErrorId
    }

    /*
     * Methods
     */

    public ellipsisUnnecessaryParams = ellipsisUnnecessaryParams
    public formatConfigExtract = formatConfigExtract
    public getTileDocUrl = getTileDocUrl
    public guessExpectedValue = guessExpectedValue
    public guessExpectedFieldName = guessExpectedFieldName
    public parsedExtractFieldValue = parsedExtractFieldValue
    public splitList = splitList
  }
</script>

<style lang="scss">
  $error-padding: 40px;

  .c-monitoror-errors {
    font-weight: 400;
    font-size: 24px;

    hr {
      width: 250px;
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

    :not(pre) > code {
      word-break: break-all;
    }

    pre {
      display: block;
      position: relative;
      margin-left: -$error-padding;
      margin-right: -$error-padding;
      background: var(--color-code-background);
      max-height: 500px;
      overflow: auto;
    }

    pre code {
      position: relative;
      display: block;
      color: var(--color-spring-wood);
    }

    pre[data-config-param] {
      &::after {
        content: "// " attr(data-config-param);
        position: absolute;
        top: 30px;
        left: $error-padding;
        color: var(--color-unknown);
        font-style: italic;
        opacity: 0.5;
      }

      code {
        padding: 60px $error-padding 30px $error-padding;
      }
    }

    .code-string {
      color: var(--color-succeeded);
    }

    .code-number {
      color: var(--color-warning);
    }

    .has-mark {
      color: rgba(255, 255, 255, 0.3);

      .code-string:not(.contains-mark),
      .code-number:not(.contains-mark) {
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

    .contains-mark mark {
      margin: 0;
    }

    .go-to-documentation {
      $go-to-doc-margin: 20px;
      margin: #{$go-to-doc-margin + 20px} #{-$error-padding + $go-to-doc-margin} $go-to-doc-margin !important;
      padding: #{$error-padding - $go-to-doc-margin};
      line-height: 1.4;
      background: linear-gradient(to right, #293536, var(--color-succeeded-dark));
      box-shadow: 3px 3px 15px rgba(23, 27, 32, .3);
      border-radius: 4px;
    }

    &:not(.c-monitoror-errors__config-verify-errors) p {
      margin: 0.2rem;
    }
  }

  .c-monitoror-errors--last-refresh {
    font-size: 22px;
  }

  .c-monitoror-errors--error {
    margin-bottom: 50px;
    font-size: 18px;
  }

  .c-monitoror-errors__config-verify-errors .c-monitoror-errors--error {
    border-radius: 4px;
    padding: 30px $error-padding 0;
    margin-top: 50px;
    background: var(--color-docs-background);
    border: 1px solid var(--color-cello);
    box-shadow: 0 0 15px var(--color-docs-background);

    & > :last-child {
      margin-bottom: 0;
    }

    .c-monitoror-errors--error-title {
      margin-top: -6px;
      margin-bottom: 1em;

      code {
        color: var(--color-failed);
        font-weight: normal;
      }
    }
  }

  .c-monitoror-errors--error-title {
    position: relative;
    font-size: 24px;
    color: #ffffff;
    margin-top: -2px !important;
    font-weight: bold;
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
