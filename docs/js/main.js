'use strict'

/*
 * Home
 */

// Update stars count
;(function () {
  const starCountElement = document.getElementById('star-count')

  if (starCountElement === null) {
    return
  }

  function setStarsCount (stars) {
    starCountElement.innerText = stars
  }

  const date = new Date()
  const starsKey = 'stars_' + date.getDay() + '_' + date.getMonth() + '_' + date.getFullYear()
  const stars = localStorage.getItem(starsKey)
  if (!stars) {
    fetch('https://api.github.com/repos/monitoror/monitoror').then((response) => {
      response.json().then((body) => {
        const stars = body.stargazers_count
        localStorage.clear()
        localStorage.setItem(starsKey, stars)
        setStarsCount(stars)
      })
    })
  } else {
    setStarsCount(stars)
  }
})()

// Show on scroll
;(function () {
  const showOnScroll = function (entries) {
    entries.forEach(entry => {
      if (entry.isIntersecting && entry.target.dataset.showOnScroll !== 'visible') {
        return entry.target.dataset.showOnScroll = 'visible'
      }

      if (!entry.isIntersecting && entry.target.dataset.showOnScroll === 'visible') {
        entry.target.dataset.showOnScroll = ''
      }
    })
  }
  const observer = new IntersectionObserver(showOnScroll)
  Array.from(document.querySelectorAll('[data-show-on-scroll]')).forEach((element) => {
    observer.observe(element)
  })
})()

/*
 * Docs
 */

// Toggle menu
;(function () {
  const toggleMenuButton = document.getElementById('js-toggle-menu')

  if (!toggleMenuButton) {
    return
  }

  toggleMenuButton.addEventListener('click', (e) => {
    e.stopPropagation()
    document.body.classList.toggle('m-documentation__sidebar-open')
  })

  Array.from(document.querySelectorAll('.m-sidebar a[href^="#"]')).forEach((link) => {
    link.addEventListener('click', function (e) {
      document.body.classList.remove('m-documentation__sidebar-open')

      if (toggleMenuButton.style.display !== 'none') {
        e.preventDefault()
        const targetElementId = this.href.split('#')[1]
        const targetPosition = document.getElementById(targetElementId).getBoundingClientRect()
        window.scroll(window.pageXOffset, targetPosition.top + window.pageYOffset - 50)
      }
    })
  })


  document.querySelector('.m-sidebar').addEventListener('click', (e) => {
    e.stopPropagation()
  })
  document.body.addEventListener('click', () => {
    document.body.classList.remove('m-documentation__sidebar-open')
  })
})()

// Run highlight.js
if (typeof hljs !== 'undefined') {
  hljs.initHighlightingOnLoad()
}

// Input with auto "select on click" behaviour
Array.from(document.querySelectorAll('[data-select-on-click]')).forEach((input) => {
  input.addEventListener('click', function () {
    this.select()
  })
})

// Set demo default state
Array.from(document.querySelectorAll('.m-documentation--demo-switch-label:first-child')).forEach((label) => {
  const tileElement = label.parentNode.parentNode.querySelector('.m-documentation--demo-tile')
  const input = label.querySelector('input')
  tileElement.classList.add('m-documentation--demo-tile__status-' + input.value)
})

// Demo state switch
Array.from(document.querySelectorAll('[data-state-switch]')).forEach((input) => {
  input.addEventListener('change', function () {
    const tileElement = this.parentNode.parentNode.parentNode.querySelector('.m-documentation--demo-tile')
    tileElement.setAttribute('class', tileElement.getAttribute('class').replace(/m-documentation--demo-tile__status-([^\s]+)/, ''))
    tileElement.classList.add('m-documentation--demo-tile__status-' + this.value)
  })
})
