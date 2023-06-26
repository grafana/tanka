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
    "url": "webpack-runtime-dda88764348a98346931.js"
  },
  {
    "url": "framework-3d5d6f3c5ba6d5f455c5.js"
  },
  {
    "url": "styles.fe4e60cd646311c716e1.css"
  },
  {
    "url": "styles-95d00f891faf7b06d026.js"
  },
  {
    "url": "f0e45107-91cefa4606c78384043e.js"
  },
  {
    "url": "app-8b002200dbbf0650e746.js"
  },
  {
    "url": "offline-plugin-app-shell-fallback/index.html",
    "revision": "981415a80efebc092177550d32421c9d"
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
    "revision": "a562923da56233a8bdebe3702d40eb16"
  },
  {
    "url": "page-data/app-data.json",
    "revision": "97c00c140d5772cba60eeace4ed329bd"
  },
  {
    "url": "polyfill-8d34c62e07e660a83f29.js"
  },
  {
    "url": "component---docs-introduction-mdx-53424fcd5b068bb753bd.js"
  },
  {
    "url": "page-data/index/page-data.json",
    "revision": "ae1c67dd0d8b8a6bfb1addee109d1e31"
  },
  {
    "url": "component---docs-installation-mdx-c6c0750e330b95ae7923.js"
  },
  {
    "url": "page-data/install/page-data.json",
    "revision": "70e59f6d7f4c06980a518910dbeba161"
  },
  {
    "url": "component---docs-tutorial-overview-mdx-74033b3037057f4a5e63.js"
  },
  {
    "url": "page-data/tutorial/overview/page-data.json",
    "revision": "4188dac6afc4f2e9c6cb4bc583711075"
  },
  {
    "url": "manifest.webmanifest",
    "revision": "f319a0c559f86a7480b781187d3324a2"
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
  pathname = pathname.replace(new RegExp(`^/pr-preview/pr-872`), ``)

  // Check for resources + the app bundle
  // The latter may not exist if the SW is updating to a new version
  const resources = await idbKeyval.get(`resources:${pathname}`)
  if (!resources || !(await caches.match(`/pr-preview/pr-872/app-8b002200dbbf0650e746.js`))) {
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

  const offlineShell = `/pr-preview/pr-872/offline-plugin-app-shell-fallback/index.html`
  const offlineShellWithKey = workbox.precaching.getCacheKeyForURL(offlineShell)
  return await caches.match(offlineShellWithKey)
})

workbox.routing.registerRoute(navigationRoute)

// this route is used when performing a non-navigation request (e.g. fetch)
workbox.routing.registerRoute(/\/.gatsby-plugin-offline:.+/, handleAPIRequest)
