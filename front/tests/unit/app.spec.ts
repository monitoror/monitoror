import { shallowMount } from '@vue/test-utils'

import App from '@/App.vue'

describe('App.vue', () => {
  it('renders HelloWorld', () => {
    const wrapper = shallowMount(App)
    expect(wrapper.find('.hello')).toBeTruthy()
  })
})
