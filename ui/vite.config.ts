import template from 'lodash.template'
import * as path from 'path'
import {UserConfig} from 'vite'
import tsResolver from 'vite-tsconfig-paths'

const BASE_URL = './'

const config: UserConfig = {
  optimizeDeps: {
    include: [
      'ts-md5/dist/md5',
    ],
    exclude: [
      'ts-md5',
    ],
  },
  assetsDir: '.',
  base: BASE_URL,
  port: 8000,
  proxy: {
    '/api': 'http://localhost:8080',
  },
  resolvers: [
    tsResolver,
  ],
  indexHtmlTransforms: [
    {
      apply: 'pre',
      transform({code}) {
        const compiled = template(code)
        return compiled({
          VITE_APP_TITLE: process.env.VITE_APP_TITLE,
          VITE_APP_CANONICAL_URL: process.env.VITE_APP_CANONICAL_URL,
          BASE_URL,
        })
      },
    }
  ],
  configureServer: ({root, app, watcher}) => {
    watcher.add(path.resolve(root, './public/**/*'))
    const publicPath = path.resolve(root, './public')
    watcher.on('change', function (path) {
      console.log(path)
      if (path.startsWith(publicPath)) {
        watcher.send({
          type: 'full-reload',
          path,
        })
      }
    })
  },
}

export default config
