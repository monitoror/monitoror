<template>
  <div>
    <h3>Hello {{ name }} from my Vue.js page, built with Webpack 4!</h3>
  </div>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      required: true
    }
  }
}
</script>
