import * as path from 'path'

import {UserConfig} from 'vite'
// @ts-ignore
import tsResolver from 'vite-tsconfig-paths'

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
  base: './',
  port: 8000,
  proxy: {
    '/api': 'http://localhost:8080',
  },
  resolvers: [
    tsResolver,
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
