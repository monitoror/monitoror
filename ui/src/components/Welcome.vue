<template>
  <div class="c-monitoror-welcome" :class="classes">
    <template v-if="isNewUser">
      <h1 class="c-monitoror-welcome--title">Welcome</h1>
      <h2 class="c-monitoror-welcome--subtitle">You're almost there!</h2>

      <div class="c-monitoror-welcome--content">
        <div class="c-monitoror-welcome--next-step-content">
          <div>
            <div class="c-monitoror-welcome--next-step-sup-title">
              1. Next step
            </div>
            <h3 class="c-monitoror-welcome--next-step-title">
              Create your <br>
              UI configuration
            </h3>
            <p>
              Define what you want to display by creating a JSON file using
              tile types.
            </p>
            <p>
              <a href="https://monitoror.com/documentation/#ui-configuration" class="c-monitoror-welcome--button" target="_blank">
                Go to documentation
                <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z" fill="currentColor"></path>
                  <path d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z" fill="currentColor"></path>
                </svg>
              </a>
            </p>
          </div>
          <div>
            <div class="c-monitoror-welcome--next-step-sup-title">
              2. And finally...
            </div>
            <h3 class="c-monitoror-welcome--next-step-title">
              Set your configuration <br>
              as default
            </h3>
            <p>
              Put the path or URL to your JSON file in
              the <code>MO_CONFIG</code> environment variable.
            </p>
            <p>
              <a href="https://monitoror.com/documentation/#core-configuration" class="c-monitoror-welcome--button" target="_blank">
                Read more about Core Configuration
                <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z" fill="currentColor"></path>
                  <path d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z" fill="currentColor"></path>
                </svg>
              </a>
            </p>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      Named config list
    </template>
  </div>
</template>

<script lang="ts">
  import {Component, Vue} from 'vue-property-decorator'

  @Component({})
  export default class MonitororWelcome extends Vue {
    /*
     * Data
     */
    private shouldHaveScroll: boolean = false

    /*
     * Computed
     */

    get classes() {
      return {
        'c-monitoror-welcome__no-scroll': !this.shouldHaveScroll,
      }
    }

    get isNewUser(): boolean {
      return this.$store.getters.isNewUser
    }

    /*
     * Methods
     */

    private mounted() {
      setTimeout(() => {
        this.shouldHaveScroll = true
      }, 1650)
    }
  }
</script>

<style lang="scss">
  .c-monitoror-welcome {
    position: relative;
    padding-top: 260px;
    font-weight: 400;
    font-size: 24px;
    text-align: center;
    animation: welcomeSlideIn 525ms ease 1.15s both;

    &::before {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 460px;
      background: var(--color-dark-background);
      transform: skewY(-7deg) scaleY(2);
      transform-origin: bottom center;
      z-index: -1;
    }

    &__no-scroll .c-monitoror-welcome--next-step-content {
      display: none;
    }

    .c-monitoror-welcome--title {
      display: inline-block;
      margin: 0 auto;
      font-weight: 700;
      font-size: 72px;
      color: transparent;
      background: linear-gradient(to right, var(--color-succeeded), var(--color-action-required));
      -webkit-background-clip: text;
      background-size: 200% 100%;
      animation: welcomeTitleGradient 10s ease 1.7s infinite;
    }

    .c-monitoror-welcome--subtitle {
      font-weight: 300;
      font-size: 32px;
      margin: 0;
    }
  }

  .c-monitoror-welcome--content {
    width: 870px;
    max-width: calc(100% - 50px);
    margin: 60px auto 0;
    padding: 130px 0 0;
    font-size: 20px;
    line-height: 1.6;
    text-align: left;
  }

  .c-monitoror-welcome--next-step-sup-title {
    text-transform: uppercase;
    font-size: 18px;
    line-height: 1.2;
    letter-spacing: 1.1px;
    color: var(--color-succeeded);
    transform: translateX(2px);
    margin-bottom: 7px;
    font-family: "JetBrains Mono", monospace;
  }

  .c-monitoror-welcome--next-step-title {
    font-size: 38px;
    line-height: 1.3;
    font-weight: 300;
    margin: 0 0 20px;
    color: #ffffff;
  }

  .c-monitoror-welcome--next-step-content {
    p {
      margin: 10px 0 30px;
    }

    & > * {
      animation: welcomeSlideIn 700ms both;
    }

    & > *:nth-child(2) {
      animation-delay: 400ms;
    }

    @media (max-width: 999px) {
      max-width: 400px;
      margin: 0 auto;

      & > * {
        padding-bottom: 70px;
      }
    }

    @media (min-width: 1000px) {
      display: grid;
      grid-auto-flow: column;
      grid-gap: 10%;
    }
  }

  .c-monitoror-welcome--button {
    display: inline-block;
    text-align: center;
    font-size: 16px;
    font-weight: 600;
    line-height: 1.2;
    letter-spacing: 0.3px;
    text-decoration: none;
    padding: 20px 30px;
    border-radius: 4px;
    border: 1px solid hsla(211, 52%, 88%, 0.15);
    color: var(--color-spindle);
    background: rgba(207, 223, 240, 0.03);
    transition: color 200ms, background 200ms, border-color 200ms;

    &::after {
      content: initial;
    }

    &:hover {
      color: var(--color-background);
      background: var(--color-spindle);
      border-color: var(--color-spindle);
    }

    svg {
      display: inline-block;
      width: 13px;
      margin-left: 10px;
      transform: translateY(1px);
    }
  }

  @keyframes welcomeTitleGradient {
    0%,
    100% {
      background-position: 0 50%;
    }
    50% {
      background-position: 100% 50%;
    }
  }

  @keyframes welcomeSlideIn {
    from {
      opacity: 0;
      transform: translateY(10%);
    }
  }
</style>
