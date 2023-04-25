/**
 * Welcome to your Workbox-powered service worker!
 *
 * You'll need to register this file in your web app and you should
 * disable HTTP caching for this file too.
 * See https://goo.gl/nhQhGp
 *
 * The rest of the code is auto-generated. Please don't update this file
 * directly; instead, make changes to your Workbox build configuration
 * and re-run your build process.
 * See https://goo.gl/2aRDsh
 */

importScripts("workbox-v4.3.1/workbox-sw.js");
workbox.setConfig({modulePathPrefix: "workbox-v4.3.1"});

workbox.core.setCacheNameDetails({prefix: "gatsby-plugin-offline"});

workbox.core.skipWaiting();

workbox.core.clientsClaim();

/**
 * The workboxSW.precacheAndRoute() method efficiently caches and responds to
 * requests for URLs in the manifest.
 * See https://goo.gl/S9QRab
 */
self.__precacheManifest = [
  {
    "url": "webpack-runtime-07166e3f152fc80004ac.js"
  },
  {
    "url": "framework-3d5d6f3c5ba6d5f455c5.js"
  },
  {
    "url": "styles.c913b28f19713112dced.css"
  },
  {
    "url": "styles-95d00f891faf7b06d026.js"
  },
  {
    "url": "f0e45107-de48ac461cc5ac1214b6.js"
  },
  {
    "url": "app-d590807b1d6af19182d6.js"
  },
  {
    "url": "offline-plugin-app-shell-fallback/index.html",
    "revision": "6b983c1bc4ebb21b830a1ba3a9b32f04"
  },
  {
    "url": "component---cache-caches-gatsby-plugin-offline-app-shell-js-099d9276a41f4ba01a3c.js"
  },
  {
    "url": "page-data/offline-plugin-app-shell-fallback/page-data.json",
    "revision": "7b70c9311b1f9e78c50c5991a4513806"
  },
  {
    "url": "page-data/sq/d/1635659820.json",
    "revision": "13b4dfef8a44f8b60e30d9fe782292de"
  },
  {
    "url": "page-data/app-data.json",
    "revision": "4596cd7ce7028f02f7542d630a17f4ec"
  },
  {
    "url": "polyfill-8d34c62e07e660a83f29.js"
  },
  {
    "url": "component---docs-introduction-mdx-2f488764d32620398d47.js"
  },
  {
    "url": "page-data/index/page-data.json",
    "revision": "ae1c67dd0d8b8a6bfb1addee109d1e31"
  },
  {
    "url": "component---docs-installation-mdx-75ff524cf4bd41f720de.js"
  },
  {
    "url": "page-data/install/page-data.json",
    "revision": "5de2f6c35bd5eab44b36f293d9e1acd0"
  },
  {
    "url": "component---docs-tutorial-overview-mdx-583246fcf54f578532fc.js"
  },
  {
    "url": "page-data/tutorial/overview/page-data.json",
    "revision": "4188dac6afc4f2e9c6cb4bc583711075"
  },
  {
    "url": "manifest.webmanifest",
    "revision": "38feed427ac73907546eeaeaf8a8c030"
  }
].concat(self.__precacheManifest || []);
workbox.precaching.precacheAndRoute(self.__precacheManifest, {});

workbox.routing.registerRoute(/(\.js$|\.css$|static\/)/, new workbox.strategies.CacheFirst(), 'GET');
workbox.routing.registerRoute(/^https?:.*\/page-data\/.*\.json/, new workbox.strategies.StaleWhileRevalidate(), 'GET');
workbox.routing.registerRoute(/^https?:.*\.(png|jpg|jpeg|webp|avif|svg|gif|tiff|js|woff|woff2|json|css)$/, new workbox.strategies.StaleWhileRevalidate(), 'GET');
workbox.routing.registerRoute(/^https?:\/\/fonts\.googleapis\.com\/css/, new workbox.strategies.StaleWhileRevalidate(), 'GET');

/* global importScripts, workbox, idbKeyval */
importScripts(`idb-keyval-3.2.0-iife.min.js`)

const { NavigationRoute } = workbox.routing

let lastNavigationRequest = null
let offlineShellEnabled = true

// prefer standard object syntax to support more browsers
const MessageAPI = {
  setPathResources: (event, { path, resources }) => {
    event.waitUntil(idbKeyval.set(`resources:${path}`, resources))
  },

  clearPathResources: event => {
    event.waitUntil(idbKeyval.clear())
  },

  enableOfflineShell: () => {
    offlineShellEnabled = true
  },

  disableOfflineShell: () => {
    offlineShellEnabled = false
  },
}

self.addEventListener(`message`, event => {
  const { gatsbyApi: api } = event.data
  if (api) MessageAPI[api](event, event.data)
})

function handleAPIRequest({ event }) {
  const { pathname } = new URL(event.request.url)

  const params = pathname.match(/:(.+)/)[1]
  const data = {}

  if (params.includes(`=`)) {
    params.split(`&`).forEach(param => {
      const [key, val] = param.split(`=`)
      data[key] = val
    })
  } else {
    data.api = params
  }

  if (MessageAPI[data.api] !== undefined) {
    MessageAPI[data.api]()
  }

  if (!data.redirect) {
    return new Response()
  }

  return new Response(null, {
    status: 302,
    headers: {
      Location: lastNavigationRequest,
    },
  })
}

const navigationRoute = new NavigationRoute(async ({ event }) => {
  // handle API requests separately to normal navigation requests, so do this
  // check first
  if (event.request.url.match(/\/.gatsby-plugin-offline:.+/)) {
    return handleAPIRequest({ event })
  }

  if (!offlineShellEnabled) {
    return await fetch(event.request)
  }

  lastNavigationRequest = event.request.url

  let { pathname } = new URL(event.request.url)
  pathname = pathname.replace(new RegExp(`^/pr-preview/pr-837`), ``)

  // Check for resources + the app bundle
  // The latter may not exist if the SW is updating to a new version
  const resources = await idbKeyval.get(`resources:${pathname}`)
  if (!resources || !(await caches.match(`/pr-preview/pr-837/app-d590807b1d6af19182d6.js`))) {
    return await fetch(event.request)
  }

  for (const resource of resources) {
    // As soon as we detect a failed resource, fetch the entire page from
    // network - that way we won't risk being in an inconsistent state with
    // some parts of the page failing.
    if (!(await caches.match(resource))) {
      return await fetch(event.request)
    }
  }

  const offlineShell = `/pr-preview/pr-837/offline-plugin-app-shell-fallback/index.html`
  const offlineShellWithKey = workbox.precaching.getCacheKeyForURL(offlineShell)
  return await caches.match(offlineShellWithKey)
})

workbox.routing.registerRoute(navigationRoute)

// this route is used when performing a non-navigation request (e.g. fetch)
workbox.routing.registerRoute(/\/.gatsby-plugin-offline:.+/, handleAPIRequest)