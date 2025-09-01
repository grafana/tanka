# Changelog

For releases from v0.24.0 to v0.28.3, you can find the changelog in the GitHub releases: https://github.com/grafana/tanka/releases

## [0.34.1](https://github.com/grafana/tanka/compare/v0.34.0...v0.34.1) (2025-09-01)


### 🐛 Bug Fixes

* **deps:** update dependency astro to v5.13.5 ([#1600](https://github.com/grafana/tanka/issues/1600)) ([d3eb55e](https://github.com/grafana/tanka/commit/d3eb55e3a647e31f76a3db188edf201b2b696970))
* **deps:** update module github.com/spf13/pflag to v1.0.8 ([#1603](https://github.com/grafana/tanka/issues/1603)) ([4a2314f](https://github.com/grafana/tanka/commit/4a2314ffdfe7a791e7f0c57c3c1ed15ff7d23fe9))
* remove debug print of carrier in otel.go ([#1606](https://github.com/grafana/tanka/issues/1606)) ([4a6d929](https://github.com/grafana/tanka/commit/4a6d9295d229c6fdd79022102e6c0e6c9b2b2395))


### 🔧 Miscellaneous Chores

* **deps:** update dagger/dagger-for-github action to v8.1.0 ([#1601](https://github.com/grafana/tanka/issues/1601)) ([ad9a3d1](https://github.com/grafana/tanka/commit/ad9a3d1d075d0ff7e0b2c1dc512b493ec055f651))

## [0.34.0](https://github.com/grafana/tanka/compare/v0.33.0...v0.34.0) (2025-08-29)


### 🎉 Features

* add OpenTelemetry tracing support ([#1598](https://github.com/grafana/tanka/issues/1598)) ([9707927](https://github.com/grafana/tanka/commit/97079276a717a3f7beb1602197e212d26a6751d1))
* support `TANKA_DANGEROUS_ALLOW_REDIRECT` env variable ([#1582](https://github.com/grafana/tanka/issues/1582)) ([42750e9](https://github.com/grafana/tanka/commit/42750e9ea8f270fdece38876b90251a185532591))


### 🐛 Bug Fixes

* **deps:** update dependency astro to v5.13.2 [security] ([#1568](https://github.com/grafana/tanka/issues/1568)) ([3b381c2](https://github.com/grafana/tanka/commit/3b381c2512bc1067a20ef4bea83693c9fd90383a))
* **deps:** update dependency astro to v5.13.4 ([#1594](https://github.com/grafana/tanka/issues/1594)) ([be783c7](https://github.com/grafana/tanka/commit/be783c7bc72e3fabc9a725ee15571d9272320f75))
* **deps:** update dependency typescript to ^5.9.2 ([#1543](https://github.com/grafana/tanka/issues/1543)) ([b008238](https://github.com/grafana/tanka/commit/b00823891dc8cc4214faf32f3eebe402da791001))
* **deps:** update k8s.io/utils digest to 0af2bda ([#1572](https://github.com/grafana/tanka/issues/1572)) ([d12a2be](https://github.com/grafana/tanka/commit/d12a2be848a248e939fd1c125585f82899845260))
* **deps:** update kubernetes packages to v0.33.4 ([#1563](https://github.com/grafana/tanka/issues/1563)) ([40e2724](https://github.com/grafana/tanka/commit/40e272470fafb69bea77721577df0ba857b6c6bb))
* **deps:** update module github.com/stretchr/testify to v1.11.0 ([#1585](https://github.com/grafana/tanka/issues/1585)) ([9df78f6](https://github.com/grafana/tanka/commit/9df78f60a12f7d20cd65331972010bfacf2b8347))
* **deps:** update module github.com/stretchr/testify to v1.11.1 ([#1595](https://github.com/grafana/tanka/issues/1595)) ([e4f5077](https://github.com/grafana/tanka/commit/e4f5077ea150a2cc5b2223ef5016f4ae6dffaa1c))
* **deps:** update module go.opentelemetry.io/proto/otlp to v1.7.1 ([#1541](https://github.com/grafana/tanka/issues/1541)) ([9192382](https://github.com/grafana/tanka/commit/91923828d23c9ad274305ef005faae7f1a5f37cf))
* **deps:** update module google.golang.org/grpc to v1.75.0 ([#1581](https://github.com/grafana/tanka/issues/1581)) ([b966c2c](https://github.com/grafana/tanka/commit/b966c2c0ece1dd522c14260d7e3899085460c255))
* **deps:** update tailwindcss monorepo to v4.1.12 ([#1564](https://github.com/grafana/tanka/issues/1564)) ([0da537b](https://github.com/grafana/tanka/commit/0da537bae0fac1bd57466ad6e8573c7fe2d60f3d))
* **setup-goversion:** fix string manipulation to exclude unwanted strings and chars ([#1561](https://github.com/grafana/tanka/issues/1561)) ([c72740b](https://github.com/grafana/tanka/commit/c72740ba9dc10fa08273aa0799977136556ea078))


### 📝 Documentation

* fix k-lib usage ([#1524](https://github.com/grafana/tanka/issues/1524)) ([997ca97](https://github.com/grafana/tanka/commit/997ca97b42237d5d5bdf7dc3e85f1fc5832e435f))


### 🏗️ Build System

* **deps:** bump @astrojs/starlight ([ca9346b](https://github.com/grafana/tanka/commit/ca9346b1f8faed649a9150e49e54fba788060613))
* **deps:** bump @astrojs/starlight from 0.34.4 to 0.34.5 in /docs ([#1511](https://github.com/grafana/tanka/issues/1511)) ([ca9346b](https://github.com/grafana/tanka/commit/ca9346b1f8faed649a9150e49e54fba788060613))
* **deps:** bump actions/create-github-app-token from 2.0.6 to 2.1.0 ([#1558](https://github.com/grafana/tanka/issues/1558)) ([a78b0df](https://github.com/grafana/tanka/commit/a78b0df80cd3489ea653b0ac3789e6f2cb083339))
* **deps:** bump astro in /docs in the docs-dependencies group ([750b1c9](https://github.com/grafana/tanka/commit/750b1c9d457a39321872177bdda91ad9bad4ea5e))
* **deps:** bump astro to 5.13.3 in /docs ([#1592](https://github.com/grafana/tanka/issues/1592)) ([750b1c9](https://github.com/grafana/tanka/commit/750b1c9d457a39321872177bdda91ad9bad4ea5e))
* **deps:** bump github.com/spf13/pflag from 1.0.6 to 1.0.7 ([#1520](https://github.com/grafana/tanka/issues/1520)) ([989d8f8](https://github.com/grafana/tanka/commit/989d8f81eb63d04884ae9d5a6d85982c1aee791e))
* **deps:** bump github.com/stretchr/testify ([70846ab](https://github.com/grafana/tanka/commit/70846abad8624575059293b1130bf14366c80bd8))
* **deps:** bump github.com/stretchr/testify from 1.10.0 to 1.11.0 ([#1589](https://github.com/grafana/tanka/issues/1589)) ([6c7ecec](https://github.com/grafana/tanka/commit/6c7ececc2494b7d31c126b1732b3d32d3f1082ac))
* **deps:** bump github.com/stretchr/testify to 1.11.0 in /acceptance-tests ([#1590](https://github.com/grafana/tanka/issues/1590)) ([70846ab](https://github.com/grafana/tanka/commit/70846abad8624575059293b1130bf14366c80bd8))
* **deps:** bump golang from 1.24.5 to 1.24.6 ([4d25565](https://github.com/grafana/tanka/commit/4d255658edee6e0956f820280264c1468a401f33))
* **deps:** bump golang from 1.24.5 to 1.25.0 ([#1552](https://github.com/grafana/tanka/issues/1552)) ([4d25565](https://github.com/grafana/tanka/commit/4d255658edee6e0956f820280264c1468a401f33))
* **deps:** bump golang.org/x/sync ([8c70a87](https://github.com/grafana/tanka/commit/8c70a872e6bb05ee0ca2bce82b2a095714303a3f))
* **deps:** bump golang.org/x/sync from 0.15.0 to 0.16.0 in /dagger ([#1514](https://github.com/grafana/tanka/issues/1514)) ([8c70a87](https://github.com/grafana/tanka/commit/8c70a872e6bb05ee0ca2bce82b2a095714303a3f))
* **deps:** bump golang.org/x/term from 0.32.0 to 0.33.0 ([#1512](https://github.com/grafana/tanka/issues/1512)) ([8ede85c](https://github.com/grafana/tanka/commit/8ede85c8b169e95b144b8ee6d4b2cc75c9bbc130))
* **deps:** bump golang.org/x/term from 0.33.0 to 0.34.0 ([#1554](https://github.com/grafana/tanka/issues/1554)) ([333ffe5](https://github.com/grafana/tanka/commit/333ffe5fc40c02cb8b3be1641e623c87bd726174))
* **deps:** bump golang.org/x/text from 0.26.0 to 0.27.0 ([#1515](https://github.com/grafana/tanka/issues/1515)) ([704a294](https://github.com/grafana/tanka/commit/704a294eaafb5f67acb3bb161067a897e53e4e53))
* **deps:** bump golang.org/x/text from 0.27.0 to 0.28.0 ([#1555](https://github.com/grafana/tanka/issues/1555)) ([89c6bc3](https://github.com/grafana/tanka/commit/89c6bc326c83219e36e0b9732c0aedb3ed40b0bc))
* **deps:** bump k8s.io/apimachinery from 0.33.2 to 0.33.3 ([#1519](https://github.com/grafana/tanka/issues/1519)) ([2dd6ac2](https://github.com/grafana/tanka/commit/2dd6ac23d69bed79ab965b7721d562246c286bd4))
* **deps:** bump ncipollo/release-action from 1.16.0 to 1.18.0 ([#1496](https://github.com/grafana/tanka/issues/1496)) ([f7a0748](https://github.com/grafana/tanka/commit/f7a07480fcc76bad8171917f1483d5e0b1c0755e))
* **deps:** bump renovatebot/github-action from 43.0.1 to 43.0.2 ([#1497](https://github.com/grafana/tanka/issues/1497)) ([b054a60](https://github.com/grafana/tanka/commit/b054a6098470e4670507bfd764e845e3b4b0e042))
* **deps:** bump renovatebot/github-action from 43.0.2 to 43.0.3 ([#1513](https://github.com/grafana/tanka/issues/1513)) ([c553a8a](https://github.com/grafana/tanka/commit/c553a8a0c542b880354132b9e181dd7b374a15a9))
* **deps:** bump renovatebot/github-action from 43.0.3 to 43.0.4 ([#1522](https://github.com/grafana/tanka/issues/1522)) ([d029c30](https://github.com/grafana/tanka/commit/d029c305624f8b9e5e496f4edc13258686324da3))
* **deps:** bump renovatebot/github-action from 43.0.4 to 43.0.5 ([#1525](https://github.com/grafana/tanka/issues/1525)) ([bc96c1d](https://github.com/grafana/tanka/commit/bc96c1d2978edfbf433573ce9528ba05836fcd4c))
* **deps:** bump renovatebot/github-action from 43.0.5 to 43.0.7 ([#1557](https://github.com/grafana/tanka/issues/1557)) ([69b9ee8](https://github.com/grafana/tanka/commit/69b9ee8045689998ae2c2d40879eecbdf83464d0))
* **deps:** bump renovatebot/github-action from 43.0.8 to 43.0.9 ([#1587](https://github.com/grafana/tanka/issues/1587)) ([40dc31f](https://github.com/grafana/tanka/commit/40dc31ffb08d001d8abeee4cdce4af17a789f0e9))
* **deps:** bump rossjrw/pr-preview-action from 1.6.1 to 1.6.2 ([#1502](https://github.com/grafana/tanka/issues/1502)) ([dfdec7f](https://github.com/grafana/tanka/commit/dfdec7f0ec15e4c47b1e898db2f1a6cc860de083))
* **deps:** bump sigs.k8s.io/yaml ([1c81a59](https://github.com/grafana/tanka/commit/1c81a597b033e620bb2b41d4ce29dfe13e692ec2))
* **deps:** bump sigs.k8s.io/yaml ([d59463e](https://github.com/grafana/tanka/commit/d59463ee3a812f8a68404ebaeee32f2bb14e74d9))
* **deps:** bump sigs.k8s.io/yaml from 1.4.0 to 1.5.0 ([#1500](https://github.com/grafana/tanka/issues/1500)) ([e6cf4e2](https://github.com/grafana/tanka/commit/e6cf4e20599f9007574d186016b72f6acb969b25))
* **deps:** bump sigs.k8s.io/yaml from 1.4.0 to 1.5.0 in /acceptance-tests ([#1499](https://github.com/grafana/tanka/issues/1499)) ([d59463e](https://github.com/grafana/tanka/commit/d59463ee3a812f8a68404ebaeee32f2bb14e74d9))
* **deps:** bump sigs.k8s.io/yaml from 1.5.0 to 1.6.0 ([#1528](https://github.com/grafana/tanka/issues/1528)) ([df800a9](https://github.com/grafana/tanka/commit/df800a9e7f764128ff8b5735e563196eeb2c7f47))
* **deps:** bump sigs.k8s.io/yaml to1.6.0 in /acceptance-tests ([#1527](https://github.com/grafana/tanka/issues/1527)) ([1c81a59](https://github.com/grafana/tanka/commit/1c81a597b033e620bb2b41d4ce29dfe13e692ec2))
* **deps:** bump the acceptance-tests-dependencies group ([b2bbb72](https://github.com/grafana/tanka/commit/b2bbb7289d7d4c78a4d35b6230176e037e493384))
* **deps:** bump the acceptance-tests-dependencies group in acceptance-tests with 2 updates ([#1518](https://github.com/grafana/tanka/issues/1518)) ([b2bbb72](https://github.com/grafana/tanka/commit/b2bbb7289d7d4c78a4d35b6230176e037e493384))
* **deps:** bump the dagger-dependencies group ([367b677](https://github.com/grafana/tanka/commit/367b6772efe714bccf4d45c9d7f404dece9d191d))
* **deps:** bump the dagger-dependencies group in /dagger with 11 updates ([#1498](https://github.com/grafana/tanka/issues/1498)) ([d09427e](https://github.com/grafana/tanka/commit/d09427e8468da85eb2338fb6021b31e9f0d61fbc))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1529](https://github.com/grafana/tanka/issues/1529)) ([367b677](https://github.com/grafana/tanka/commit/367b6772efe714bccf4d45c9d7f404dece9d191d))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1501](https://github.com/grafana/tanka/issues/1501)) ([88010cf](https://github.com/grafana/tanka/commit/88010cfe8b256ac30b6cad32fea16f1cffd85a6e))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1526](https://github.com/grafana/tanka/issues/1526)) ([17cf293](https://github.com/grafana/tanka/commit/17cf29324b9a24c7da95668bb9981f3f802fca7d))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1546](https://github.com/grafana/tanka/issues/1546)) ([045ae32](https://github.com/grafana/tanka/commit/045ae320ea7b6f6dfcb60e829134438d7d960cc4))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1556](https://github.com/grafana/tanka/issues/1556)) ([dd97c55](https://github.com/grafana/tanka/commit/dd97c55105d02f541df7d190e1920cc89a2dcd3a))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1521](https://github.com/grafana/tanka/issues/1521)) ([f8fb6d3](https://github.com/grafana/tanka/commit/f8fb6d359cebce4eb329f3b032a4c762d595363f))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1495](https://github.com/grafana/tanka/issues/1495)) ([09946fd](https://github.com/grafana/tanka/commit/09946fd5c3dc4491d63530579801280a7c579220))


### 🔧 Miscellaneous Chores

* **deps:** pin dependencies ([#1532](https://github.com/grafana/tanka/issues/1532)) ([47082b8](https://github.com/grafana/tanka/commit/47082b84f0fda9756a8de519960887f0461936e4))
* **deps:** update actions/cache action to v4.2.4 ([#1550](https://github.com/grafana/tanka/issues/1550)) ([e06912b](https://github.com/grafana/tanka/commit/e06912b27e51bf1551cc9162f202f0f0f2482d7c))
* **deps:** update actions/checkout action to v4.3.0 ([#1565](https://github.com/grafana/tanka/issues/1565)) ([8144f86](https://github.com/grafana/tanka/commit/8144f867a1c122190d37c497730d1d40f6dae05a))
* **deps:** update actions/checkout action to v5 ([#1586](https://github.com/grafana/tanka/issues/1586)) ([d8757cd](https://github.com/grafana/tanka/commit/d8757cd174aa335eab5f48cecc11b5fe3b267237))
* **deps:** update actions/create-github-app-token action to v2.1.1 ([#1573](https://github.com/grafana/tanka/issues/1573)) ([cbd727a](https://github.com/grafana/tanka/commit/cbd727a49ee6674336852f589a9ba6de8693dc98))
* **deps:** update actions/download-artifact action to v5 ([#1547](https://github.com/grafana/tanka/issues/1547)) ([7096802](https://github.com/grafana/tanka/commit/70968029c5c6202603245ade7f17841c228eb519))
* **deps:** update actions/setup-go action to v5.5.0 ([#1534](https://github.com/grafana/tanka/issues/1534)) ([beeb921](https://github.com/grafana/tanka/commit/beeb9217ffc5c28c3d5cbd0cfc1af1cb8f828dd2))
* **deps:** update azure/setup-helm action to v4.3.1 ([#1574](https://github.com/grafana/tanka/issues/1574)) ([f75b9ad](https://github.com/grafana/tanka/commit/f75b9ad9bca0e9fd0de1aeb8184c8252cec3a181))
* **deps:** update dependency @types/node to v24.0.12 ([#1506](https://github.com/grafana/tanka/issues/1506)) ([bb989fc](https://github.com/grafana/tanka/commit/bb989fce59c295d884d3a0fd3b277d81a06a1f68))
* **deps:** update dependency @types/node to v24.0.13 ([#1510](https://github.com/grafana/tanka/issues/1510)) ([dddb355](https://github.com/grafana/tanka/commit/dddb355f0cbca166dd68735a193f5d8181e0f995))
* **deps:** update dependency @types/node to v24.3.0 ([#1576](https://github.com/grafana/tanka/issues/1576)) ([8d27649](https://github.com/grafana/tanka/commit/8d2764947248581cb0d0371c9fa55a6bddaf9857))
* **deps:** update dependency go to v1.24.4 ([#1504](https://github.com/grafana/tanka/issues/1504)) ([63f4bd7](https://github.com/grafana/tanka/commit/63f4bd7dcab7157d1a145aca76cf956b398e9eb8))
* **deps:** update dependency go to v1.24.5 ([#1507](https://github.com/grafana/tanka/issues/1507)) ([1260749](https://github.com/grafana/tanka/commit/1260749040499a4f5aca507d6e328f56ffc38ca7))
* **deps:** update dependency go to v1.24.6 ([#1548](https://github.com/grafana/tanka/issues/1548)) ([5c70662](https://github.com/grafana/tanka/commit/5c70662f3009e7cac0846cac38c2b6406fe94af0))
* **deps:** update dependency go to v1.25.0 ([#1577](https://github.com/grafana/tanka/issues/1577)) ([e270383](https://github.com/grafana/tanka/commit/e270383f5132c1f941e847b311f4ecb985ba4780))
* **deps:** update dependency helm to v3.18.4 ([#1505](https://github.com/grafana/tanka/issues/1505)) ([6f689e3](https://github.com/grafana/tanka/commit/6f689e3a72efcf636a6e4279cb00c042c02176ed))
* **deps:** update dependency helm to v3.18.5 ([#1562](https://github.com/grafana/tanka/issues/1562)) ([8182ce0](https://github.com/grafana/tanka/commit/8182ce0022524ac467753264ba4293bb35f30cfa))
* **deps:** update dependency helm to v3.18.6 ([#1567](https://github.com/grafana/tanka/issues/1567)) ([7bd5bbb](https://github.com/grafana/tanka/commit/7bd5bbb2252c1b2ddb84c557f43242309f8584cc))
* **deps:** update dependency kubectl to v1.33.3 ([#1516](https://github.com/grafana/tanka/issues/1516)) ([b5b6d72](https://github.com/grafana/tanka/commit/b5b6d723ef09003ed5807688ef65223d8d8d2ec2))
* **deps:** update dependency kubectl to v1.33.4 ([#1560](https://github.com/grafana/tanka/issues/1560)) ([bf90224](https://github.com/grafana/tanka/commit/bf902249d27e99ab266cf530fe30cefff90356cf))
* **deps:** update dependency kubectl to v1.34.0 ([#1596](https://github.com/grafana/tanka/issues/1596)) ([e08ead5](https://github.com/grafana/tanka/commit/e08ead538ecfa416fc089a1ac2aa3abf294b762e))
* **deps:** update dependency kustomize to v5.7.0 ([#1493](https://github.com/grafana/tanka/issues/1493)) ([d1af250](https://github.com/grafana/tanka/commit/d1af25043d8b38d3c9cd002b6555d64ae889a371))
* **deps:** update dependency kustomize to v5.7.1 ([#1523](https://github.com/grafana/tanka/issues/1523)) ([b75931e](https://github.com/grafana/tanka/commit/b75931e3a71cc46c5e311336e19d520549f4c017))
* **deps:** update dependency sharp to v0.34.3 ([#1508](https://github.com/grafana/tanka/issues/1508)) ([8e79042](https://github.com/grafana/tanka/commit/8e79042c089408309e548b13a7af1d271392d822))
* **deps:** update docker/metadata-action action to v5.8.0 ([#1544](https://github.com/grafana/tanka/issues/1544)) ([8c9967b](https://github.com/grafana/tanka/commit/8c9967b73681df43fc830e988b85c6adfc9bb9eb))
* **deps:** update golang docker tag to v1.24.5 ([#1509](https://github.com/grafana/tanka/issues/1509)) ([6a5c559](https://github.com/grafana/tanka/commit/6a5c55958370b243d984156bfddf81cf35d95c93))
* **deps:** update golang docker tag to v1.25.0 ([#1549](https://github.com/grafana/tanka/issues/1549)) ([7adb868](https://github.com/grafana/tanka/commit/7adb868585979803b1ec6557201c59d6d7b801b5))
* **deps:** update golang:1.25.0 docker digest to 5502b0e ([#1583](https://github.com/grafana/tanka/issues/1583)) ([e55d6ff](https://github.com/grafana/tanka/commit/e55d6ff5f838d0f7cdf4e824805296388afadd9c))
* **deps:** update golang:1.25.0 docker digest to 91e2cd4 ([#1569](https://github.com/grafana/tanka/issues/1569)) ([e873c0f](https://github.com/grafana/tanka/commit/e873c0fa67452f94074d5438e80bb4a0de97aa75))
* **deps:** update golang:1.25.0-alpine docker digest to f18a072 ([#1570](https://github.com/grafana/tanka/issues/1570)) ([3a5c518](https://github.com/grafana/tanka/commit/3a5c51840804cd674b6afd7af519829ac9f45a45))
* **deps:** update googleapis/release-please-action action to v4.3.0 ([#1578](https://github.com/grafana/tanka/issues/1578)) ([56578af](https://github.com/grafana/tanka/commit/56578af179648e1d1681d7aa66d77ec04e6b046a))
* **deps:** update grafana/shared-workflows/dockerhub-login action to v1.0.2 ([#1533](https://github.com/grafana/tanka/issues/1533)) ([4eeb9ec](https://github.com/grafana/tanka/commit/4eeb9ec55c8ebca7a4ccd87f39e0e9196886ca63))
* **deps:** update grafana/shared-workflows/get-vault-secrets action to v1.2.1 ([#1535](https://github.com/grafana/tanka/issues/1535)) ([4caf6a5](https://github.com/grafana/tanka/commit/4caf6a5f21b524fbc6b581ad7b3aefcc0e1e2041))
* **deps:** update grafana/shared-workflows/get-vault-secrets action to v1.3.0 ([#1579](https://github.com/grafana/tanka/issues/1579)) ([139b4ca](https://github.com/grafana/tanka/commit/139b4caaf320928bf8b975a91435548a342cb3c3))
* **deps:** update grafana/shared-workflows/lint-pr-title action to v1.2.0 ([#1537](https://github.com/grafana/tanka/issues/1537)) ([0c77a14](https://github.com/grafana/tanka/commit/0c77a146dbe6c299281f0fc0eafebbd4a9d08f73))
* **deps:** update k8s.io/utils digest to 4c0f3b2 ([#1503](https://github.com/grafana/tanka/issues/1503)) ([c083023](https://github.com/grafana/tanka/commit/c0830235fcba635d43a3f21891fcadae68c4fefd))
* **deps:** update pnpm to v10 [security] ([#1531](https://github.com/grafana/tanka/issues/1531)) ([dbdf9d3](https://github.com/grafana/tanka/commit/dbdf9d3318c1d9cf2e5ee2ca0b6161d28b8a8b48))
* **deps:** update pnpm to v10.14.0 ([#1542](https://github.com/grafana/tanka/issues/1542)) ([7ce35fe](https://github.com/grafana/tanka/commit/7ce35fe5afa951cf5356441df73e6fdc2613927e))
* **deps:** update pnpm to v10.15.0 ([#1580](https://github.com/grafana/tanka/issues/1580)) ([6cec1e9](https://github.com/grafana/tanka/commit/6cec1e953a00d715c4f5a82c9081df2dfa2143a9))
* **deps:** update renovatebot/github-action action to v43.0.7 ([#1551](https://github.com/grafana/tanka/issues/1551)) ([ae73eca](https://github.com/grafana/tanka/commit/ae73eca6cd7df03ef2fe183dc001e0cd6cf27176))
* **deps:** update renovatebot/github-action action to v43.0.8 ([#1575](https://github.com/grafana/tanka/issues/1575)) ([b426691](https://github.com/grafana/tanka/commit/b4266918657ba0dd6c7cb0a78e6bfc532f65c3d1))
* **deps:** update renovatebot/github-action action to v43.0.9 ([#1584](https://github.com/grafana/tanka/issues/1584)) ([cfff6ee](https://github.com/grafana/tanka/commit/cfff6ee1b5bec14c58e577e135906ca0f99bbac7))

## [0.33.0](https://github.com/grafana/tanka/compare/v0.32.0...v0.33.0) (2025-06-25)


### 🎉 Features

* **jsonnet lib:** add function to find transitive importers ([#1464](https://github.com/grafana/tanka/issues/1464)) ([fa219d3](https://github.com/grafana/tanka/commit/fa219d35d24f14acdb86fc9954b26d17f0865ac7))


### 🏗️ Build System

* **deps:** bump actions/create-github-app-token from 2.0.2 to 2.0.6 ([#1446](https://github.com/grafana/tanka/issues/1446)) ([47892f0](https://github.com/grafana/tanka/commit/47892f02dda7aeffdff3e7088bb545cd6b791f12))
* **deps:** bump actions/download-artifact from 4.2.1 to 4.3.0 ([#1438](https://github.com/grafana/tanka/issues/1438)) ([cd1123d](https://github.com/grafana/tanka/commit/cd1123d164289be5666155372252bff2ca2015fb))
* **deps:** bump alpine from 3.21 to 3.22 ([#1469](https://github.com/grafana/tanka/issues/1469)) ([7a918ba](https://github.com/grafana/tanka/commit/7a918baf52c4a4a7bff77f1172a278cbd4b1f6e0))
* **deps:** bump docker/build-push-action from 6.15.0 to 6.16.0 ([#1437](https://github.com/grafana/tanka/issues/1437)) ([a08234f](https://github.com/grafana/tanka/commit/a08234feafd4da3c2fde4dd77baaa8aa85707817))
* **deps:** bump docker/build-push-action from 6.16.0 to 6.17.0 ([#1462](https://github.com/grafana/tanka/issues/1462)) ([b176e7e](https://github.com/grafana/tanka/commit/b176e7ef511ae88d5216c0e38a345837834b9e72))
* **deps:** bump docker/build-push-action from 6.17.0 to 6.18.0 ([#1472](https://github.com/grafana/tanka/issues/1472)) ([664d712](https://github.com/grafana/tanka/commit/664d712315f2a3651daf78eea839b9f55094a264))
* **deps:** bump docker/setup-buildx-action from 3.10.0 to 3.11.1 ([#1485](https://github.com/grafana/tanka/issues/1485)) ([af014ab](https://github.com/grafana/tanka/commit/af014abe95d1dcb88a72ff06ba0af11bb0a101b0))
* **deps:** bump github.com/99designs/gqlgen ([e06239b](https://github.com/grafana/tanka/commit/e06239bd160f6250665a33e92957a53c0673e077))
* **deps:** bump github.com/99designs/gqlgen from 0.17.74 to 0.17.75 in /dagger ([#1488](https://github.com/grafana/tanka/issues/1488)) ([e06239b](https://github.com/grafana/tanka/commit/e06239bd160f6250665a33e92957a53c0673e077))
* **deps:** bump github.com/google/go-jsonnet from 0.20.0 to 0.21.0 ([#1449](https://github.com/grafana/tanka/issues/1449)) ([ff44275](https://github.com/grafana/tanka/commit/ff442751bc61ca7126c2b74e1749d3e1e258a71e))
* **deps:** bump github.com/vektah/gqlparser/v2 ([866326a](https://github.com/grafana/tanka/commit/866326abd04570416cbe13bfd29506c28fffcd94))
* **deps:** bump github.com/vektah/gqlparser/v2 from 2.5.27 to 2.5.28 ([#1480](https://github.com/grafana/tanka/issues/1480)) ([866326a](https://github.com/grafana/tanka/commit/866326abd04570416cbe13bfd29506c28fffcd94))
* **deps:** bump golang from 1.24.2 to 1.24.3 ([#1454](https://github.com/grafana/tanka/issues/1454)) ([f14153a](https://github.com/grafana/tanka/commit/f14153aca5762c98178edf5b31247e8f61c19011))
* **deps:** bump golang from 1.24.3 to 1.24.4 ([#1476](https://github.com/grafana/tanka/issues/1476)) ([c043bf0](https://github.com/grafana/tanka/commit/c043bf026196c814f5f01e2328f94bf1c11718e5))
* **deps:** bump golang.org/x/term from 0.31.0 to 0.32.0 ([#1450](https://github.com/grafana/tanka/issues/1450)) ([5a70656](https://github.com/grafana/tanka/commit/5a70656863ad6bc863bd488debebf394b4e415ba))
* **deps:** bump golang.org/x/text from 0.24.0 to 0.25.0 ([#1451](https://github.com/grafana/tanka/issues/1451)) ([8136dc5](https://github.com/grafana/tanka/commit/8136dc5e00f3bd2c9923648bea9b88a6e2b145a2))
* **deps:** bump golang.org/x/text from 0.25.0 to 0.26.0 ([#1475](https://github.com/grafana/tanka/issues/1475)) ([940c4d3](https://github.com/grafana/tanka/commit/940c4d37b9d888c7868e6d07120a9aa8bf9a4f83))
* **deps:** bump k8s.io/apimachinery from 0.32.3 to 0.33.0 ([#1440](https://github.com/grafana/tanka/issues/1440)) ([0a97be9](https://github.com/grafana/tanka/commit/0a97be99be9d62ab43932fe28c41b398b9107088))
* **deps:** bump k8s.io/apimachinery from 0.33.0 to 0.33.1 ([#1457](https://github.com/grafana/tanka/issues/1457)) ([2f632a7](https://github.com/grafana/tanka/commit/2f632a7d0d08ae39ef4b1f398a23c4896ce59dc9))
* **deps:** bump k8s.io/apimachinery from 0.33.1 to 0.33.2 ([#1487](https://github.com/grafana/tanka/issues/1487)) ([a0df773](https://github.com/grafana/tanka/commit/a0df7737b66a01edb57e3a64bc0e145977115f82))
* **deps:** bump renovatebot/github-action from 41.0.21 to 41.0.22 ([#1436](https://github.com/grafana/tanka/issues/1436)) ([7b1367c](https://github.com/grafana/tanka/commit/7b1367ce08ecb7cac062289c82fcecd691d63469))
* **deps:** bump renovatebot/github-action from 41.0.22 to 42.0.1 ([#1447](https://github.com/grafana/tanka/issues/1447)) ([b95f821](https://github.com/grafana/tanka/commit/b95f8211847161ec0ca6516a47eab482fb106c2e))
* **deps:** bump renovatebot/github-action from 42.0.1 to 42.0.2 ([#1453](https://github.com/grafana/tanka/issues/1453)) ([19fb364](https://github.com/grafana/tanka/commit/19fb364b8d93303b11da577e2e7d552276f3be32))
* **deps:** bump renovatebot/github-action from 42.0.2 to 42.0.3 ([#1461](https://github.com/grafana/tanka/issues/1461)) ([9a18527](https://github.com/grafana/tanka/commit/9a18527ef1e4a5e06f11ed2d157d98bc1622d202))
* **deps:** bump renovatebot/github-action from 42.0.3 to 42.0.4 ([#1470](https://github.com/grafana/tanka/issues/1470)) ([acf65e8](https://github.com/grafana/tanka/commit/acf65e8ae44f478abd09de215646518b5210d81e))
* **deps:** bump renovatebot/github-action from 42.0.4 to 42.0.5 ([#1477](https://github.com/grafana/tanka/issues/1477)) ([7c417b9](https://github.com/grafana/tanka/commit/7c417b990ada7ccd68c4bd0a252a99b9613c24ee))
* **deps:** bump renovatebot/github-action from 42.0.5 to 42.0.6 ([#1482](https://github.com/grafana/tanka/issues/1482)) ([b6228b4](https://github.com/grafana/tanka/commit/b6228b4f8b95250cb7927b8d6eedeca66e459efb))
* **deps:** bump renovatebot/github-action from 42.0.6 to 43.0.1 ([#1486](https://github.com/grafana/tanka/issues/1486)) ([162ca11](https://github.com/grafana/tanka/commit/162ca1115a80a877cb355407535ed054e71d389c))
* **deps:** bump the acceptance-tests-dependencies group ([d6fe7e9](https://github.com/grafana/tanka/commit/d6fe7e96c91332efe554ed16195802bc7eef2c7f))
* **deps:** bump the acceptance-tests-dependencies group ([#1490](https://github.com/grafana/tanka/issues/1490)) ([dacefb8](https://github.com/grafana/tanka/commit/dacefb8a59a429bde240a0029d7c7bd5d3594319))
* **deps:** bump the acceptance-tests-dependencies group in /acceptance-tests ([#1460](https://github.com/grafana/tanka/issues/1460)) ([d6fe7e9](https://github.com/grafana/tanka/commit/d6fe7e96c91332efe554ed16195802bc7eef2c7f))
* **deps:** bump the acceptance-tests-dependencies group with 2 updates ([#1435](https://github.com/grafana/tanka/issues/1435)) ([3de332b](https://github.com/grafana/tanka/commit/3de332b3558ae692ebe799e49454bca1916bdd18))
* **deps:** bump the dagger-dependencies group ([57ad0f6](https://github.com/grafana/tanka/commit/57ad0f67b739c765e1153ba731e09435bf66fda5))
* **deps:** bump the dagger-dependencies group ([2045ea2](https://github.com/grafana/tanka/commit/2045ea20e1546b452617b97b2cc2b02fc4eb95a7))
* **deps:** bump the dagger-dependencies group ([e1ebc9d](https://github.com/grafana/tanka/commit/e1ebc9d04a492d2513855404413ab12ae584943d))
* **deps:** bump the dagger-dependencies group ([23016f6](https://github.com/grafana/tanka/commit/23016f6ac6f9103abfefd35892a62a650fd6cf14))
* **deps:** bump the dagger-dependencies group ([cddd02a](https://github.com/grafana/tanka/commit/cddd02af3ca0f5e7d1ca83fc30a85df01a7d94b1))
* **deps:** bump the dagger-dependencies group in /dagger with 11 updates ([#1471](https://github.com/grafana/tanka/issues/1471)) ([2045ea2](https://github.com/grafana/tanka/commit/2045ea20e1546b452617b97b2cc2b02fc4eb95a7))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1439](https://github.com/grafana/tanka/issues/1439)) ([f505b1b](https://github.com/grafana/tanka/commit/f505b1bbf654b268e5bc16b29bf051556d4b0e47))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1445](https://github.com/grafana/tanka/issues/1445)) ([cddd02a](https://github.com/grafana/tanka/commit/cddd02af3ca0f5e7d1ca83fc30a85df01a7d94b1))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1455](https://github.com/grafana/tanka/issues/1455)) ([23016f6](https://github.com/grafana/tanka/commit/23016f6ac6f9103abfefd35892a62a650fd6cf14))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1459](https://github.com/grafana/tanka/issues/1459)) ([e1ebc9d](https://github.com/grafana/tanka/commit/e1ebc9d04a492d2513855404413ab12ae584943d))
* **deps:** bump the dagger-dependencies group in /dagger with 3 updates ([#1478](https://github.com/grafana/tanka/issues/1478)) ([57ad0f6](https://github.com/grafana/tanka/commit/57ad0f67b739c765e1153ba731e09435bf66fda5))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1434](https://github.com/grafana/tanka/issues/1434)) ([979d796](https://github.com/grafana/tanka/commit/979d796805adea4b7b4525aac96ecf9ce897bcc9))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1474](https://github.com/grafana/tanka/issues/1474)) ([a718ef3](https://github.com/grafana/tanka/commit/a718ef3d2d87de1e7ae39260345a06b2d8252670))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1489](https://github.com/grafana/tanka/issues/1489)) ([543670a](https://github.com/grafana/tanka/commit/543670a8ad496273f2d0a9b70180d0274418a938))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1444](https://github.com/grafana/tanka/issues/1444)) ([d021aca](https://github.com/grafana/tanka/commit/d021acafd4fd6648862a2b6c669a6d6efdfb7f51))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1458](https://github.com/grafana/tanka/issues/1458)) ([3766943](https://github.com/grafana/tanka/commit/3766943f4fa12468282934676b3970df16751cdd))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1452](https://github.com/grafana/tanka/issues/1452)) ([45c21f2](https://github.com/grafana/tanka/commit/45c21f2c72bf5d95ad714b1cac5f092ee2cb81fb))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1468](https://github.com/grafana/tanka/issues/1468)) ([ba6bd90](https://github.com/grafana/tanka/commit/ba6bd90f46ad1c6eb40abe27275cc418976d9bc2))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1481](https://github.com/grafana/tanka/issues/1481)) ([d34903e](https://github.com/grafana/tanka/commit/d34903e477db5ebd69328c29c309e380df8f8b89))


### 🤖 Continuous Integration

* apply zizmor findings ([#1441](https://github.com/grafana/tanka/issues/1441)) ([0cbb9b6](https://github.com/grafana/tanka/commit/0cbb9b62a8456d47ff61ab0903352d66c772e7c3))
* escape docker output ([#1443](https://github.com/grafana/tanka/issues/1443)) ([c7f0cd3](https://github.com/grafana/tanka/commit/c7f0cd3e339833c55e4bb71adb24c2bac331300d))


### 🔧 Miscellaneous Chores

* **deps:** update dependency helm to v3.18.0 ([#1463](https://github.com/grafana/tanka/issues/1463)) ([d1e6347](https://github.com/grafana/tanka/commit/d1e63477d97dbfeef4ba2b17d8504557799ca48e))
* **deps:** update dependency helm to v3.18.1 ([#1466](https://github.com/grafana/tanka/issues/1466)) ([aa93a16](https://github.com/grafana/tanka/commit/aa93a160b3f11320a929da796ccaa22ccdc6529b))
* **deps:** update dependency helm to v3.18.2 ([#1473](https://github.com/grafana/tanka/issues/1473)) ([7d6d2a9](https://github.com/grafana/tanka/commit/7d6d2a9193c81c2622919d7c63cb11ee5e4da718))
* **deps:** update dependency helm to v3.18.3 ([#1483](https://github.com/grafana/tanka/issues/1483)) ([f8a0f93](https://github.com/grafana/tanka/commit/f8a0f9301ab2328d10170d20bbdee263797cd321))
* **deps:** update dependency kubectl to v1.33.0 ([#1432](https://github.com/grafana/tanka/issues/1432)) ([ee74690](https://github.com/grafana/tanka/commit/ee7469029d12b3cf9714c669ad1240a4e0b19891))
* **deps:** update dependency kubectl to v1.33.1 ([#1456](https://github.com/grafana/tanka/issues/1456)) ([26b308c](https://github.com/grafana/tanka/commit/26b308cba85aa15f721c9bcdffc45ae37c079775))
* **deps:** update dependency kubectl to v1.33.2 ([#1484](https://github.com/grafana/tanka/issues/1484)) ([7eb5730](https://github.com/grafana/tanka/commit/7eb573048cec3122a8e058a55563e5f3eb6d5572))
* **tests:** drop use of testify/mock.Mock.TestData ([#1467](https://github.com/grafana/tanka/issues/1467)) ([433e534](https://github.com/grafana/tanka/commit/433e53471564ccb102378e0a78f28d88b81d0176))

## [0.32.0](https://github.com/grafana/tanka/compare/v0.31.3...v0.32.0) (2025-04-23)


### 🎉 Features

* allow special char in repo name ([ea63f8d](https://github.com/grafana/tanka/commit/ea63f8d443ff2ca38b60065583f3bb2111ad30cf))


### 🐛 Bug Fixes

* **helm:** allow special char in chart repo name ([#1366](https://github.com/grafana/tanka/issues/1366)) ([ea63f8d](https://github.com/grafana/tanka/commit/ea63f8d443ff2ca38b60065583f3bb2111ad30cf))
* **jsonnet/implementations:** capture stderr separately ([#1423](https://github.com/grafana/tanka/issues/1423)) ([1b26f20](https://github.com/grafana/tanka/commit/1b26f2012f1bc508063fe835c1c30b4243320cbf))
* **tanka/inline:** ensure Peek only grabs metadata ([#1425](https://github.com/grafana/tanka/issues/1425)) ([6408b5f](https://github.com/grafana/tanka/commit/6408b5f84e411d378813c4dca91b0fc4f4f24122))


### 📝 Documentation

* update docs to reduce confusion in tutorial ([#1391](https://github.com/grafana/tanka/issues/1391)) ([5a3aa1e](https://github.com/grafana/tanka/commit/5a3aa1e780aff643e94d298c502886688d49f3ed))


### 🏗️ Build System

* **deps:** bump actions/cache from 4.2.2 to 4.2.3 ([#1399](https://github.com/grafana/tanka/issues/1399)) ([3dea644](https://github.com/grafana/tanka/commit/3dea6442b0062e2dc298b80aa331a3fcc71de561))
* **deps:** bump actions/create-github-app-token from 1.11.6 to 1.11.7 ([#1400](https://github.com/grafana/tanka/issues/1400)) ([df27333](https://github.com/grafana/tanka/commit/df273331ca76c85d005ca90d7edebb88157e69d3))
* **deps:** bump actions/create-github-app-token from 1.11.7 to 1.12.0 ([#1409](https://github.com/grafana/tanka/issues/1409)) ([5fc266e](https://github.com/grafana/tanka/commit/5fc266eed807749d1af146d838c7b1e84af9e5cb))
* **deps:** bump actions/create-github-app-token from 1.12.0 to 2.0.2 ([#1414](https://github.com/grafana/tanka/issues/1414)) ([eda32fb](https://github.com/grafana/tanka/commit/eda32fbb8fd99ed93c2048c261a56d3c6f9bf2c8))
* **deps:** bump actions/download-artifact from 4.1.9 to 4.2.1 ([#1408](https://github.com/grafana/tanka/issues/1408)) ([7f1af96](https://github.com/grafana/tanka/commit/7f1af96ba07215238302bd079b1168000aaecdb6))
* **deps:** bump actions/setup-node from 4.2.0 to 4.3.0 ([#1398](https://github.com/grafana/tanka/issues/1398)) ([c439965](https://github.com/grafana/tanka/commit/c4399652919e5b5be5b6e9adcaab5d4438b1c1ba))
* **deps:** bump actions/setup-node from 4.3.0 to 4.4.0 ([#1419](https://github.com/grafana/tanka/issues/1419)) ([083219b](https://github.com/grafana/tanka/commit/083219b056e760e355186d22f0b5a87c2cd4de87))
* **deps:** bump actions/upload-artifact from 4.6.1 to 4.6.2 ([#1402](https://github.com/grafana/tanka/issues/1402)) ([40df914](https://github.com/grafana/tanka/commit/40df91418a55470428f9467ee7d481881e8e13ce))
* **deps:** bump docker/setup-buildx-action from 3.9.0 to 3.10.0 ([#1397](https://github.com/grafana/tanka/issues/1397)) ([b9651e1](https://github.com/grafana/tanka/commit/b9651e12bf2a0fc711f869235ed3d3455c8dfff5))
* **deps:** bump github.com/99designs/gqlgen ([4e04a2a](https://github.com/grafana/tanka/commit/4e04a2a6f401cdc9217d3abb5dfd782c1a57ca63))
* **deps:** bump github.com/99designs/gqlgen ([c9fb06d](https://github.com/grafana/tanka/commit/c9fb06da2fa2ae0d4061337a26615b293ebd41b0))
* **deps:** bump github.com/99designs/gqlgen from 0.17.68 to 0.17.70 ([#1406](https://github.com/grafana/tanka/issues/1406)) ([4e04a2a](https://github.com/grafana/tanka/commit/4e04a2a6f401cdc9217d3abb5dfd782c1a57ca63))
* **deps:** bump github.com/99designs/gqlgen to 0.17.68 in /dagger ([#1392](https://github.com/grafana/tanka/issues/1392)) ([c9fb06d](https://github.com/grafana/tanka/commit/c9fb06da2fa2ae0d4061337a26615b293ebd41b0))
* **deps:** bump github.com/rs/zerolog from 1.33.0 to 1.34.0 ([#1404](https://github.com/grafana/tanka/issues/1404)) ([39b5ed4](https://github.com/grafana/tanka/commit/39b5ed48e8c805b97ce7e1a3135a728aa69ac288))
* **deps:** bump github.com/vektah/gqlparser/v2 ([1aa3b1e](https://github.com/grafana/tanka/commit/1aa3b1e1512b70b92929eabd7d8722160871b104))
* **deps:** bump github.com/vektah/gqlparser/v2 from 2.5.23 to 2.5.24 ([#1417](https://github.com/grafana/tanka/issues/1417)) ([1aa3b1e](https://github.com/grafana/tanka/commit/1aa3b1e1512b70b92929eabd7d8722160871b104))
* **deps:** bump golang from 1.24.1 to 1.24.2 ([#1411](https://github.com/grafana/tanka/issues/1411)) ([3b6b0ed](https://github.com/grafana/tanka/commit/3b6b0ed93c46aa71de5e9623ac38f29ef45bb471))
* **deps:** bump golang.org/x/crypto from 0.31.0 to 0.35.0 ([#1424](https://github.com/grafana/tanka/issues/1424)) ([14242d4](https://github.com/grafana/tanka/commit/14242d40bb4162de568d67be9ed143ab3f306df6))
* **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 in /acceptance-tests ([#1388](https://github.com/grafana/tanka/issues/1388)) ([b036484](https://github.com/grafana/tanka/commit/b0364844f66c920fe6e8920ca1dd502655efd64b))
* **deps:** bump golang.org/x/net from 0.35.0 to 0.36.0 in /dagger ([#1389](https://github.com/grafana/tanka/issues/1389)) ([c97ec16](https://github.com/grafana/tanka/commit/c97ec1643fc04eda1ae6ac98838d3f927b60d257))
* **deps:** bump golang.org/x/net from 0.36.0 to 0.38.0 in /acceptance-tests ([#1426](https://github.com/grafana/tanka/issues/1426)) ([7c95553](https://github.com/grafana/tanka/commit/7c95553839c41fe45d3ec2eab2c8f61d097220a0))
* **deps:** bump golang.org/x/net from 0.37.0 to 0.38.0 in /dagger ([#1427](https://github.com/grafana/tanka/issues/1427)) ([298940f](https://github.com/grafana/tanka/commit/298940fe618c9cb16f79d248379dea05ffc8f434))
* **deps:** bump golang.org/x/net in /acceptance-tests ([7c95553](https://github.com/grafana/tanka/commit/7c95553839c41fe45d3ec2eab2c8f61d097220a0))
* **deps:** bump golang.org/x/net in /acceptance-tests ([b036484](https://github.com/grafana/tanka/commit/b0364844f66c920fe6e8920ca1dd502655efd64b))
* **deps:** bump golang.org/x/term from 0.30.0 to 0.31.0 ([#1413](https://github.com/grafana/tanka/issues/1413)) ([a193121](https://github.com/grafana/tanka/commit/a19312141a36b3ec480c3e3e24de5cf430bd2544))
* **deps:** bump golang.org/x/text from 0.23.0 to 0.24.0 ([#1412](https://github.com/grafana/tanka/issues/1412)) ([89d0c8d](https://github.com/grafana/tanka/commit/89d0c8df7943b8a89c58fd5f0b15067fee285e8d))
* **deps:** bump k8s.io/apimachinery from 0.32.2 to 0.32.3 ([#1395](https://github.com/grafana/tanka/issues/1395)) ([f229afe](https://github.com/grafana/tanka/commit/f229afe07738784f12afbd27ce056adc10a15c02))
* **deps:** bump renovatebot/github-action from 41.0.14 to 41.0.16 ([#1396](https://github.com/grafana/tanka/issues/1396)) ([e496dab](https://github.com/grafana/tanka/commit/e496dab7146f3b8678b4848e4592dcbd1f05b0b1))
* **deps:** bump renovatebot/github-action from 41.0.16 to 41.0.17 ([#1401](https://github.com/grafana/tanka/issues/1401)) ([0a913d4](https://github.com/grafana/tanka/commit/0a913d45576a3a54bf12bf6da92fba97d3cd57b7))
* **deps:** bump renovatebot/github-action from 41.0.17 to 41.0.18 ([#1407](https://github.com/grafana/tanka/issues/1407)) ([52b1a0f](https://github.com/grafana/tanka/commit/52b1a0fb16a7ad08a378dea0c80cd129a13921e6))
* **deps:** bump renovatebot/github-action from 41.0.18 to 41.0.20 ([#1420](https://github.com/grafana/tanka/issues/1420)) ([c6fbf8a](https://github.com/grafana/tanka/commit/c6fbf8ad0609ef18f0b05d5201c5376b88a78217))
* **deps:** bump renovatebot/github-action from 41.0.20 to 41.0.21 ([#1430](https://github.com/grafana/tanka/issues/1430)) ([b698943](https://github.com/grafana/tanka/commit/b6989430b29286a1703293b7b68233f4ce346f19))
* **deps:** bump rossjrw/pr-preview-action from 1.6.0 to 1.6.1 ([#1421](https://github.com/grafana/tanka/issues/1421)) ([b19ca92](https://github.com/grafana/tanka/commit/b19ca9251bf9f91d41c157ef24b4c8c956a13f6a))
* **deps:** bump the acceptance-tests-dependencies group ([6f395af](https://github.com/grafana/tanka/commit/6f395affdc383fa340cf695f079de33737dafa8f))
* **deps:** bump the acceptance-tests-dependencies group with 2 updates ([#1393](https://github.com/grafana/tanka/issues/1393)) ([6f395af](https://github.com/grafana/tanka/commit/6f395affdc383fa340cf695f079de33737dafa8f))
* **deps:** bump the dagger-dependencies group ([ff10be3](https://github.com/grafana/tanka/commit/ff10be30cbf23be3426b95af59e25cb40d1d3ad2))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1415](https://github.com/grafana/tanka/issues/1415)) ([ff10be3](https://github.com/grafana/tanka/commit/ff10be30cbf23be3426b95af59e25cb40d1d3ad2))
* **deps:** bump the dagger-dependencies group in /dagger with 2 updates ([#1429](https://github.com/grafana/tanka/issues/1429)) ([73e201e](https://github.com/grafana/tanka/commit/73e201e3ce1102d9e9a57052f355184c8e56e00e))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1394](https://github.com/grafana/tanka/issues/1394)) ([aa57c47](https://github.com/grafana/tanka/commit/aa57c47fae97e3a83c111622fe2d1e8845e99d7b))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1418](https://github.com/grafana/tanka/issues/1418)) ([924e875](https://github.com/grafana/tanka/commit/924e87522105fe00c2df8e8e41974bf5a26cdf50))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1405](https://github.com/grafana/tanka/issues/1405)) ([02c0aa4](https://github.com/grafana/tanka/commit/02c0aa47486e8df84452457d51e22c17537b6a9f))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1403](https://github.com/grafana/tanka/issues/1403)) ([0cd2506](https://github.com/grafana/tanka/commit/0cd25068522f77e2387293329ec6b7d02715336d))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1428](https://github.com/grafana/tanka/issues/1428)) ([366f67e](https://github.com/grafana/tanka/commit/366f67e4206c2f953b8021f494ed856b20fd711a))
* **deps:** bump the docs-dependencies group in /docs with 6 updates ([#1410](https://github.com/grafana/tanka/issues/1410)) ([d2f6a87](https://github.com/grafana/tanka/commit/d2f6a87be6995d8dd759fb7989431c2e25232495))


### 🔧 Miscellaneous Chores

* **deps:** update dependency helm to v3.17.2 ([#1390](https://github.com/grafana/tanka/issues/1390)) ([9ae944f](https://github.com/grafana/tanka/commit/9ae944f2a1c3f82fe1d129d558c36bf14810104b))
* **deps:** update dependency helm to v3.17.3 ([#1416](https://github.com/grafana/tanka/issues/1416)) ([3fe52eb](https://github.com/grafana/tanka/commit/3fe52eb3d7f449dd6d608fdca972d7aebd4509c1))
* **deps:** update dependency kubectl to v1.32.3 ([#1387](https://github.com/grafana/tanka/issues/1387)) ([90158a3](https://github.com/grafana/tanka/commit/90158a35b560520ededea18a70bb9e27603053eb))
* **deps:** update dependency kubectl to v1.32.4 ([#1431](https://github.com/grafana/tanka/issues/1431)) ([bfc9e15](https://github.com/grafana/tanka/commit/bfc9e15b51bd6348584942f6f5cffeffc81b9771))

## [0.31.3](https://github.com/grafana/tanka/compare/v0.31.2...v0.31.3) (2025-03-10)


### 🐛 Bug Fixes

* mimic --tla-code behavior from jsonnet for functions in main.jsonnet in env discovery ([#1251](https://github.com/grafana/tanka/issues/1251)) ([3065778](https://github.com/grafana/tanka/commit/3065778d03f9cecb4896fc0275af510c836e991d))


### 📝 Documentation

* add note about Tanka deploying test resources ([#1329](https://github.com/grafana/tanka/issues/1329)) ([d54c97f](https://github.com/grafana/tanka/commit/d54c97faf168dab7c51f08597e1a1505ade3761e))


### 🏗️ Build System

* **deps:** bump actions/cache from 4.2.0 to 4.2.1 ([#1362](https://github.com/grafana/tanka/issues/1362)) ([210ee51](https://github.com/grafana/tanka/commit/210ee5194b69e4373fa99893422a7b7029aec871))
* **deps:** bump actions/cache from 4.2.1 to 4.2.2 ([#1371](https://github.com/grafana/tanka/issues/1371)) ([fbfa04f](https://github.com/grafana/tanka/commit/fbfa04fd4ce43c9526a120d82fc7f14ab21fcec3))
* **deps:** bump actions/create-github-app-token from 1.11.0 to 1.11.1 ([#1323](https://github.com/grafana/tanka/issues/1323)) ([be354cc](https://github.com/grafana/tanka/commit/be354cca593e07d17f32262ae7b81ffb6690a5c4))
* **deps:** bump actions/create-github-app-token from 1.11.1 to 1.11.2 ([#1330](https://github.com/grafana/tanka/issues/1330)) ([0e930c6](https://github.com/grafana/tanka/commit/0e930c6432ba97687d7d5e6e36825d5d76239acf))
* **deps:** bump actions/create-github-app-token from 1.11.2 to 1.11.3 ([#1341](https://github.com/grafana/tanka/issues/1341)) ([aa5d4d1](https://github.com/grafana/tanka/commit/aa5d4d180d704fc008c0f37fefdedc446118c697))
* **deps:** bump actions/create-github-app-token from 1.11.3 to 1.11.5 ([#1355](https://github.com/grafana/tanka/issues/1355)) ([06b3df4](https://github.com/grafana/tanka/commit/06b3df46e6ec4cc36138e527acb2b03498fd06e3))
* **deps:** bump actions/create-github-app-token from 1.11.5 to 1.11.6 ([#1370](https://github.com/grafana/tanka/issues/1370)) ([e3e6046](https://github.com/grafana/tanka/commit/e3e60468e9eed82f52e265f948a5a0bc06b49c8b))
* **deps:** bump actions/download-artifact from 4.1.8 to 4.1.9 ([#1385](https://github.com/grafana/tanka/issues/1385)) ([e90062a](https://github.com/grafana/tanka/commit/e90062ae01d7ad963c8c7c37e9040f95e61f9d98))
* **deps:** bump actions/setup-node from 4.1.0 to 4.2.0 ([#1321](https://github.com/grafana/tanka/issues/1321)) ([1067607](https://github.com/grafana/tanka/commit/106760704673bf8e8a61e64e3b08120677251d65))
* **deps:** bump actions/upload-artifact from 4.6.0 to 4.6.1 ([#1360](https://github.com/grafana/tanka/issues/1360)) ([cc3c4a3](https://github.com/grafana/tanka/commit/cc3c4a3b66f79f1e33f32c8252ddfb359dd5431f))
* **deps:** bump astro ([e75bbf8](https://github.com/grafana/tanka/commit/e75bbf888c4fcfe1fb933fc60e947c35b3dca603))
* **deps:** bump astro from 5.1.7 to 5.1.10 in /docs ([#1328](https://github.com/grafana/tanka/issues/1328)) ([e75bbf8](https://github.com/grafana/tanka/commit/e75bbf888c4fcfe1fb933fc60e947c35b3dca603))
* **deps:** bump azure/setup-helm from 4.2.0 to 4.3.0 ([#1361](https://github.com/grafana/tanka/issues/1361)) ([2f870c2](https://github.com/grafana/tanka/commit/2f870c20494727f13b79baa496bda6d9479087fd))
* **deps:** bump dagger/dagger-for-github from 7.0.4 to 7.0.5 ([#1340](https://github.com/grafana/tanka/issues/1340)) ([5645dfd](https://github.com/grafana/tanka/commit/5645dfdf23d929fbc73b62205dcbbf832832dadb))
* **deps:** bump dagger/dagger-for-github from 7.0.5 to 7.0.6 ([#1356](https://github.com/grafana/tanka/issues/1356)) ([816fc4b](https://github.com/grafana/tanka/commit/816fc4b1e4f2ab77f83bdffe920e16c98e7e8b59))
* **deps:** bump dagger/dagger-for-github from 7.0.6 to 8.0.0 ([#1383](https://github.com/grafana/tanka/issues/1383)) ([65d55c9](https://github.com/grafana/tanka/commit/65d55c9c42c4c8e5ee8cf7b119a5e450033719a2))
* **deps:** bump docker/build-push-action from 6.10.0 to 6.13.0 ([#1320](https://github.com/grafana/tanka/issues/1320)) ([fb3d545](https://github.com/grafana/tanka/commit/fb3d545762daa11bd0fdd8874a8b410eda674f79))
* **deps:** bump docker/build-push-action from 6.13.0 to 6.14.0 ([#1363](https://github.com/grafana/tanka/issues/1363)) ([a71fd1b](https://github.com/grafana/tanka/commit/a71fd1b71ef882d6e2769886360bbcb73fb6baf0))
* **deps:** bump docker/build-push-action from 6.14.0 to 6.15.0 ([#1381](https://github.com/grafana/tanka/issues/1381)) ([63aee73](https://github.com/grafana/tanka/commit/63aee73667250e1207ff715e772465d4c364ff13))
* **deps:** bump docker/metadata-action from 5.6.1 to 5.7.0 ([#1369](https://github.com/grafana/tanka/issues/1369)) ([8fdda37](https://github.com/grafana/tanka/commit/8fdda3780cb37768f3a357fa9d07c004e2a655fb))
* **deps:** bump docker/setup-buildx-action from 3.7.1 to 3.8.0 ([#1332](https://github.com/grafana/tanka/issues/1332)) ([b6ffaa9](https://github.com/grafana/tanka/commit/b6ffaa9f151e211ec8051a906bc51998a1d111a6))
* **deps:** bump docker/setup-buildx-action from 3.8.0 to 3.9.0 ([#1343](https://github.com/grafana/tanka/issues/1343)) ([65c0b67](https://github.com/grafana/tanka/commit/65c0b67c1361715088ca4055d19e0bae56a074e0))
* **deps:** bump github.com/99designs/gqlgen ([e0ad1a5](https://github.com/grafana/tanka/commit/e0ad1a5fa404789ed7c6ba23527becf56e05fcab))
* **deps:** bump github.com/99designs/gqlgen from 0.17.64 to 0.17.66 ([#1353](https://github.com/grafana/tanka/issues/1353)) ([e0ad1a5](https://github.com/grafana/tanka/commit/e0ad1a5fa404789ed7c6ba23527becf56e05fcab))
* **deps:** bump github.com/google/go-cmp from 0.6.0 to 0.7.0 ([#1364](https://github.com/grafana/tanka/issues/1364)) ([5daa8a8](https://github.com/grafana/tanka/commit/5daa8a841bb6f87980a0fbdfb624e11c3a1fc5dd))
* **deps:** bump github.com/spf13/pflag from 1.0.5 to 1.0.6 ([#1333](https://github.com/grafana/tanka/issues/1333)) ([64c51f9](https://github.com/grafana/tanka/commit/64c51f996ec8ac41ce0a627b337eebc98dc9f654))
* **deps:** bump github.com/vektah/gqlparser/v2 ([bc3e6d1](https://github.com/grafana/tanka/commit/bc3e6d183a0b35c878def9b89ab812de173ca9f9))
* **deps:** bump github.com/vektah/gqlparser/v2 from 2.5.22 to 2.5.23 ([#1372](https://github.com/grafana/tanka/issues/1372)) ([bc3e6d1](https://github.com/grafana/tanka/commit/bc3e6d183a0b35c878def9b89ab812de173ca9f9))
* **deps:** bump golang from 1.23.5 to 1.23.6 ([#1344](https://github.com/grafana/tanka/issues/1344)) ([04acba0](https://github.com/grafana/tanka/commit/04acba05c14923044fb335bbb82ac09cccd6efb1))
* **deps:** bump golang from 1.23.6 to 1.24.0 ([#1350](https://github.com/grafana/tanka/issues/1350)) ([65c8b12](https://github.com/grafana/tanka/commit/65c8b1277b6ef9fb87a006c2478051e33121cc41))
* **deps:** bump golang from 1.24.0 to 1.24.1 ([#1380](https://github.com/grafana/tanka/issues/1380)) ([426b010](https://github.com/grafana/tanka/commit/426b0106672aeb4ecc92415e5d7d680fcc71ec1b))
* **deps:** bump golang.org/x/sync ([b0938b8](https://github.com/grafana/tanka/commit/b0938b82f38d639110708b05dc820df9525191db))
* **deps:** bump golang.org/x/sync from 0.10.0 to 0.11.0 in the dagger-dependencies group ([#1346](https://github.com/grafana/tanka/issues/1346)) ([b0938b8](https://github.com/grafana/tanka/commit/b0938b82f38d639110708b05dc820df9525191db))
* **deps:** bump golang.org/x/term from 0.28.0 to 0.29.0 ([#1337](https://github.com/grafana/tanka/issues/1337)) ([7357261](https://github.com/grafana/tanka/commit/73572618605a9b4fc976524c5fab3be9f288fb15))
* **deps:** bump golang.org/x/term from 0.29.0 to 0.30.0 ([#1376](https://github.com/grafana/tanka/issues/1376)) ([e149906](https://github.com/grafana/tanka/commit/e149906a74e42157c41aedc697b9b7beb9de980e))
* **deps:** bump golang.org/x/text from 0.21.0 to 0.22.0 ([#1338](https://github.com/grafana/tanka/issues/1338)) ([bf90b60](https://github.com/grafana/tanka/commit/bf90b601b41474317d42bb916f2fe5f46116c734))
* **deps:** bump golang.org/x/text from 0.22.0 to 0.23.0 ([#1377](https://github.com/grafana/tanka/issues/1377)) ([444d4ee](https://github.com/grafana/tanka/commit/444d4ee07a50e815175390b8c57d807c763edc4f))
* **deps:** bump google.golang.org/grpc ([d349a7a](https://github.com/grafana/tanka/commit/d349a7a8437db9fa696127a290b70d4430d2cbed))
* **deps:** bump google.golang.org/grpc from 1.69.4 to 1.70.0 ([#1326](https://github.com/grafana/tanka/issues/1326)) ([d349a7a](https://github.com/grafana/tanka/commit/d349a7a8437db9fa696127a290b70d4430d2cbed))
* **deps:** bump googleapis/release-please-action from 4.1.3 to 4.2.0 ([#1384](https://github.com/grafana/tanka/issues/1384)) ([ead7027](https://github.com/grafana/tanka/commit/ead70279180afdc4dc821262a4210b7ea914c7bc))
* **deps:** bump JamesIves/github-pages-deploy-action ([8fc36b6](https://github.com/grafana/tanka/commit/8fc36b6d2f678b04dc68d8be8c3b32286733bd34))
* **deps:** bump JamesIves/github-pages-deploy-action from 4.7.2 to 4.7.3 ([#1367](https://github.com/grafana/tanka/issues/1367)) ([8fc36b6](https://github.com/grafana/tanka/commit/8fc36b6d2f678b04dc68d8be8c3b32286733bd34))
* **deps:** bump k8s.io/apimachinery from 0.32.1 to 0.32.2 ([#1354](https://github.com/grafana/tanka/issues/1354)) ([b5bb9ce](https://github.com/grafana/tanka/commit/b5bb9cec513e53811f34b6470ed28cf1754ad300))
* **deps:** bump ncipollo/release-action from 1.15.0 to 1.16.0 ([#1368](https://github.com/grafana/tanka/issues/1368)) ([9f92ae2](https://github.com/grafana/tanka/commit/9f92ae2add9ea4b48766551fd1ec648756ed0d41))
* **deps:** bump pnpm/action-setup from 4.0.0 to 4.1.0 ([#1339](https://github.com/grafana/tanka/issues/1339)) ([7ca7ab1](https://github.com/grafana/tanka/commit/7ca7ab1810c3fc8109a66c8985bfdbc3a3b1848f))
* **deps:** bump renovatebot/github-action from 41.0.11 to 41.0.12 ([#1331](https://github.com/grafana/tanka/issues/1331)) ([98f7b1d](https://github.com/grafana/tanka/commit/98f7b1dadfe5b6601162370f00fef04e244aa5d2))
* **deps:** bump renovatebot/github-action from 41.0.12 to 41.0.13 ([#1342](https://github.com/grafana/tanka/issues/1342)) ([18c6f45](https://github.com/grafana/tanka/commit/18c6f45b82bc42f80289404623812ec9203df675))
* **deps:** bump renovatebot/github-action from 41.0.13 to 41.0.14 ([#1359](https://github.com/grafana/tanka/issues/1359)) ([d2fa103](https://github.com/grafana/tanka/commit/d2fa103f162e9310c307f7c92433bef2253487e4))
* **deps:** bump renovatebot/github-action from 41.0.6 to 41.0.11 ([#1322](https://github.com/grafana/tanka/issues/1322)) ([d46c9bb](https://github.com/grafana/tanka/commit/d46c9bb26d187c4185df2549608b403360171334))
* **deps:** bump the acceptance-tests-dependencies group ([#1352](https://github.com/grafana/tanka/issues/1352)) ([8102deb](https://github.com/grafana/tanka/commit/8102debf3cee8f7e4e1ceb02e887e3dc8625758b))
* **deps:** bump the dagger-dependencies group ([05dfa38](https://github.com/grafana/tanka/commit/05dfa38aff1cd0cb56071696f0013c953fe6aca1))
* **deps:** bump the dagger-dependencies group ([3d19732](https://github.com/grafana/tanka/commit/3d197322d1760859709b2de2ffef571ffcfb810e))
* **deps:** bump the dagger-dependencies group in /dagger with 11 updates ([#1379](https://github.com/grafana/tanka/issues/1379)) ([05dfa38](https://github.com/grafana/tanka/commit/05dfa38aff1cd0cb56071696f0013c953fe6aca1))
* **deps:** bump the dagger-dependencies group in /dagger with 3 updates ([#1335](https://github.com/grafana/tanka/issues/1335)) ([3d19732](https://github.com/grafana/tanka/commit/3d197322d1760859709b2de2ffef571ffcfb810e))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1334](https://github.com/grafana/tanka/issues/1334)) ([50b0b27](https://github.com/grafana/tanka/commit/50b0b278e1c0948030f62b1ba5dd330230a0cc1f))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1378](https://github.com/grafana/tanka/issues/1378)) ([ebe3764](https://github.com/grafana/tanka/commit/ebe37644de10ed81701beecf0eee4a15e6c230b2))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1345](https://github.com/grafana/tanka/issues/1345)) ([cd614ed](https://github.com/grafana/tanka/commit/cd614ed72c026813115315479c788f7aa6f1fc17))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1358](https://github.com/grafana/tanka/issues/1358)) ([28068dc](https://github.com/grafana/tanka/commit/28068dcf8ee6451389d83d853bc338155521f506))
* **deps:** bump the docs-dependencies group in /docs with 5 updates ([#1351](https://github.com/grafana/tanka/issues/1351)) ([fa1c248](https://github.com/grafana/tanka/commit/fa1c248d49abc7caab8dd218c377a2507480e492))
* **deps:** bump the docs-dependencies group in /docs with 7 updates ([#1373](https://github.com/grafana/tanka/issues/1373)) ([1b1f030](https://github.com/grafana/tanka/commit/1b1f030f56b5bcbf73a4779c5f339d00d7331e35))


### 🤖 Continuous Integration

* use ubuntu-24.04 and ubuntu-24.04-arm runners ([#1357](https://github.com/grafana/tanka/issues/1357)) ([656b9c0](https://github.com/grafana/tanka/commit/656b9c0435252d155338b1824fbb99ceb9d758a2))


### 🔧 Miscellaneous Chores

* **deps:** update dependency helm to v3.17.1 ([#1348](https://github.com/grafana/tanka/issues/1348)) ([d0df7f1](https://github.com/grafana/tanka/commit/d0df7f1f8bfdc5400b8a500fb38eb177f7803f90))
* **deps:** update dependency kubectl to v1.32.2 ([#1349](https://github.com/grafana/tanka/issues/1349)) ([a71826d](https://github.com/grafana/tanka/commit/a71826de254863af2af27f71468c2cdfa9822ebf))
* **docs:** bump to tailwind 4.0 ([#1327](https://github.com/grafana/tanka/issues/1327)) ([6c5d340](https://github.com/grafana/tanka/commit/6c5d340ef9cdd8992f9e96f421e4f5ae7ee81ef1))
* pin login-to-dockerhub action ([#1365](https://github.com/grafana/tanka/issues/1365)) ([3d87962](https://github.com/grafana/tanka/commit/3d879625b5a61ecd7fbe34df55b61fc35b26c7ba))


### ✅ Tests

* add acceptance test for tk export -l ([#1347](https://github.com/grafana/tanka/issues/1347)) ([3109dc9](https://github.com/grafana/tanka/commit/3109dc9d63da175d95646c955c34ff9934ea97da))

## [0.31.2](https://github.com/grafana/tanka/compare/v0.31.1...v0.31.2) (2025-01-20)


### 🐛 Bug Fixes

* `Unexpected input(s) 'github-token'` in `release-please` job ([#1301](https://github.com/grafana/tanka/issues/1301)) ([e465a31](https://github.com/grafana/tanka/commit/e465a3123ea9b18f6b3afe4414f7bb80b57efa32))


### 🏗️ Build System

* **deps:** bump actions/upload-artifact from 4.4.3 to 4.6.0 ([#1315](https://github.com/grafana/tanka/issues/1315)) ([bb482d9](https://github.com/grafana/tanka/commit/bb482d927062d2fe0f5a53b482dd36353c50520b))
* **deps:** bump astro from 5.1.1 to 5.1.2 in /docs in the docs-dependencies group ([#1303](https://github.com/grafana/tanka/issues/1303)) ([27991f4](https://github.com/grafana/tanka/commit/27991f4c9cfab28bf366e6c3991bbfa12dac374e))
* **deps:** bump dagger/dagger-for-github from 7.0.3 to 7.0.4 ([#1314](https://github.com/grafana/tanka/issues/1314)) ([e7acea4](https://github.com/grafana/tanka/commit/e7acea427c7fb4441e5796c0530a337c6f6e7b6e))
* **deps:** bump github.com/99designs/gqlgen from 0.17.61 to 0.17.62 in /dagger ([#1302](https://github.com/grafana/tanka/issues/1302)) ([59b78ed](https://github.com/grafana/tanka/commit/59b78ed994966256025b69b381c3376e5acc3ac6))
* **deps:** bump golang from 1.23.4 to 1.23.5 ([#1312](https://github.com/grafana/tanka/issues/1312)) ([d6e7417](https://github.com/grafana/tanka/commit/d6e74172984dfb268e887ccb94d8f68c0c50dc49))
* **deps:** bump golang.org/x/term from 0.27.0 to 0.28.0 ([#1304](https://github.com/grafana/tanka/issues/1304)) ([c071693](https://github.com/grafana/tanka/commit/c071693977b19aa35d05179f33313426f4c88ce2))
* **deps:** bump k8s.io/apimachinery from 0.32.0 to 0.32.1 ([#1318](https://github.com/grafana/tanka/issues/1318)) ([d8ba4dd](https://github.com/grafana/tanka/commit/d8ba4dd02e3e620e29dc60bd2e5dbbfa250e83f0))
* **deps:** bump ncipollo/release-action from 1.14.0 to 1.15.0 ([#1313](https://github.com/grafana/tanka/issues/1313)) ([8504473](https://github.com/grafana/tanka/commit/85044735ab2932f46095fe392411ce155f2fa88d))
* **deps:** bump rossjrw/pr-preview-action from 1.4.8 to 1.6.0 ([#1316](https://github.com/grafana/tanka/issues/1316)) ([40ac63d](https://github.com/grafana/tanka/commit/40ac63dda4481cf6080d0b83b09f5209a6abbed2))
* **deps:** bump the acceptance-tests-dependencies group with 2 updates ([#1317](https://github.com/grafana/tanka/issues/1317)) ([8821f83](https://github.com/grafana/tanka/commit/8821f832624d394ec17dc1ec5f02c744b033c593))
* **deps:** bump the dagger-dependencies group in /dagger with 3 updates ([#1306](https://github.com/grafana/tanka/issues/1306)) ([8d77eb3](https://github.com/grafana/tanka/commit/8d77eb3646193b93c1445272e1a7e5a873640c21))
* **deps:** bump the dagger-dependencies group in /dagger with 9 updates ([#1319](https://github.com/grafana/tanka/issues/1319)) ([41228bb](https://github.com/grafana/tanka/commit/41228bb023f02d54fe432a6ba6b3ebafad01d5e8))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1311](https://github.com/grafana/tanka/issues/1311)) ([dcc213d](https://github.com/grafana/tanka/commit/dcc213ddde50e06c1efe9bdfddb3eb45b50455d0))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1307](https://github.com/grafana/tanka/issues/1307)) ([7c745e1](https://github.com/grafana/tanka/commit/7c745e1fd1f5357637aaa18f2b3afef1992fac1b))


### 🔧 Miscellaneous Chores

* **deps:** update dependency helm to v3.17.0 ([#1310](https://github.com/grafana/tanka/issues/1310)) ([9811551](https://github.com/grafana/tanka/commit/9811551e0d030e07900935d5a429e56fa72cfeee))
* **deps:** update dependency kubectl to v1.32.1 ([#1309](https://github.com/grafana/tanka/issues/1309)) ([2a6ae21](https://github.com/grafana/tanka/commit/2a6ae21326cc3550ac33276f214520d18abd76ab))
* **deps:** update dependency kustomize to v5.6.0 ([#1308](https://github.com/grafana/tanka/issues/1308)) ([43069eb](https://github.com/grafana/tanka/commit/43069ebb99fa51aab9eeb825a0db1987617ddc30))

## [0.31.1](https://github.com/grafana/tanka/compare/v0.31.0...v0.31.1) (2025-01-02)


### 🏗️ Build System

* **deps:** bump astro from 5.0.5 to 5.0.8 in /docs ([#1293](https://github.com/grafana/tanka/issues/1293)) ([060d5e8](https://github.com/grafana/tanka/commit/060d5e88b733b3792e23fa98945d3c82451450ac))
* **deps:** bump the dagger-dependencies group across 1 directory with 3 updates ([#1297](https://github.com/grafana/tanka/issues/1297)) ([d7857f0](https://github.com/grafana/tanka/commit/d7857f0ab75384aab5c6399c07f3c1d8e17b4565))
* **deps:** bump the docs-dependencies group across 1 directory with 5 updates ([#1296](https://github.com/grafana/tanka/issues/1296)) ([cb353f8](https://github.com/grafana/tanka/commit/cb353f8c6a21e230e8ca12cc771be93d2b69aa02))


### 🤖 Continuous Integration

* bump x/net to v0.33 in acceptance-tests ([#1298](https://github.com/grafana/tanka/issues/1298)) ([50f7149](https://github.com/grafana/tanka/commit/50f714975b06d8df6362a84703d9fa1187f534f3))
* use ncipollo/release-action for uploading release artifacts ([#1292](https://github.com/grafana/tanka/issues/1292)) ([4df6ff7](https://github.com/grafana/tanka/commit/4df6ff7c347b01a177eb730cb9368d4f1d9ec4e1))


### 🔧 Miscellaneous Chores

* **deps:** update dependency helm to v3.16.4 ([#1290](https://github.com/grafana/tanka/issues/1290)) ([336b926](https://github.com/grafana/tanka/commit/336b926b117bcb2e3c52facc66a5647ab42d6aa1))

## [0.31.0](https://github.com/grafana/tanka/compare/v0.30.2...v0.31.0) (2024-12-16)


### 🎉 Features

* support --{tla,ext}-{code,str}-file flag in "tk eval" ([#1238](https://github.com/grafana/tanka/issues/1238)) ([a93627a](https://github.com/grafana/tanka/commit/a93627ab3abb165f3d2323abb277fda5bda1fb46))


### 🏗️ Build System

* **deps:** bump actions/cache from 4.1.2 to 4.2.0 ([#1267](https://github.com/grafana/tanka/issues/1267)) ([c3f9ceb](https://github.com/grafana/tanka/commit/c3f9ceb35dd22056302a24c0460bd10da4ba932f))
* **deps:** bump alpine from 3.20 to 3.21 ([#1265](https://github.com/grafana/tanka/issues/1265)) ([b9f4911](https://github.com/grafana/tanka/commit/b9f49116b87764636f4fa26aa29f786f9a83bbef))
* **deps:** bump dagger/dagger-for-github from 7.0.1 to 7.0.3 ([#1285](https://github.com/grafana/tanka/issues/1285)) ([a5ec928](https://github.com/grafana/tanka/commit/a5ec928aa7b793b67f7f6599fbb23e977ca90327))
* **deps:** bump github.com/99designs/gqlgen from 0.17.57 to 0.17.60 in /dagger ([#1276](https://github.com/grafana/tanka/issues/1276)) ([71defaa](https://github.com/grafana/tanka/commit/71defaa3eeb0a836b5d6944fd5aaebb081489123))
* **deps:** bump golang from 1.23.3 to 1.23.4 ([#1266](https://github.com/grafana/tanka/issues/1266)) ([7f18b87](https://github.com/grafana/tanka/commit/7f18b87f291dfa6d7e6c8d12f63024f394637e0c))
* **deps:** bump golang.org/x/crypto from 0.26.0 to 0.31.0 ([#1284](https://github.com/grafana/tanka/issues/1284)) ([6885695](https://github.com/grafana/tanka/commit/68856959a6abc32d5d0a28f9f9879fb32632d553))
* **deps:** bump golang.org/x/term from 0.26.0 to 0.27.0 ([#1264](https://github.com/grafana/tanka/issues/1264)) ([dc946ad](https://github.com/grafana/tanka/commit/dc946ad3228c31d33efa16e7f3aaa80de0569557))
* **deps:** bump golang.org/x/text from 0.20.0 to 0.21.0 ([#1263](https://github.com/grafana/tanka/issues/1263)) ([95258f7](https://github.com/grafana/tanka/commit/95258f75ee28c4defcc8c759af2231db807fdac5))
* **deps:** bump JamesIves/github-pages-deploy-action from 4.7.1 to 4.7.2 ([#1269](https://github.com/grafana/tanka/issues/1269)) ([5b59c97](https://github.com/grafana/tanka/commit/5b59c9758bab3a26cdf21c6afd0a3b3bfd81a8a9))
* **deps:** bump k8s.io/apimachinery from 0.31.3 to 0.31.4 ([#1275](https://github.com/grafana/tanka/issues/1275)) ([333fc0d](https://github.com/grafana/tanka/commit/333fc0d18f52d6839ff53ee19ce34053689f7461))
* **deps:** bump k8s.io/apimachinery from 0.31.4 to 0.32.0 ([#1283](https://github.com/grafana/tanka/issues/1283)) ([b475bca](https://github.com/grafana/tanka/commit/b475bca24e40c830316ffbcbbf1fdbfb9383aa18))
* **deps:** bump renovatebot/github-action from 41.0.5 to 41.0.6 ([#1268](https://github.com/grafana/tanka/issues/1268)) ([bf679bf](https://github.com/grafana/tanka/commit/bf679bf09e4aab27c5ce7fc08fb965ce397bee54))
* **deps:** bump the acceptance-tests-dependencies group with 2 updates ([#1277](https://github.com/grafana/tanka/issues/1277)) ([7b5140c](https://github.com/grafana/tanka/commit/7b5140c1e421f524131238dc17ae9ccc8fd89e4c))
* **deps:** bump the acceptance-tests-dependencies group with 2 updates ([#1288](https://github.com/grafana/tanka/issues/1288)) ([9b33ff5](https://github.com/grafana/tanka/commit/9b33ff5a5450e0321932f4bf2acdc9350217de77))
* **deps:** bump the dagger-dependencies group in /dagger with 10 updates ([#1287](https://github.com/grafana/tanka/issues/1287)) ([5952a38](https://github.com/grafana/tanka/commit/5952a38a3cd15c81fc18c8a115913e0fa3d3865d))
* **deps:** bump the dagger-dependencies group in /dagger with 3 updates ([#1271](https://github.com/grafana/tanka/issues/1271)) ([ddb7d4e](https://github.com/grafana/tanka/commit/ddb7d4e59e6db2e806ed8097cc2069d642b69731))
* **deps:** bump the docs-dependencies group across 1 directory with 3 updates ([#1289](https://github.com/grafana/tanka/issues/1289)) ([7df2d2b](https://github.com/grafana/tanka/commit/7df2d2bf1b0d93c86f4183b3dbf2591c951dc3cc))
* **deps:** bump the docs-dependencies group in /docs with 4 updates ([#1278](https://github.com/grafana/tanka/issues/1278)) ([7aba4bd](https://github.com/grafana/tanka/commit/7aba4bd84d97682d2a86489ac24b1dcff5b88563))


### 🤖 Continuous Integration

* add renovate ([#1262](https://github.com/grafana/tanka/issues/1262)) ([3c9a48d](https://github.com/grafana/tanka/commit/3c9a48d04ee48ef6c1f6dbd7b11c28c71c8ed5a2))
* ignore Astro 5 for now as Starlight does not support it yet ([#1274](https://github.com/grafana/tanka/issues/1274)) ([30d907e](https://github.com/grafana/tanka/commit/30d907e08374f196a64ad241fecce8618c5ea6eb))


### 🔧 Miscellaneous Chores

* **deps:** update dependency kubectl to v1.31.4 ([#1273](https://github.com/grafana/tanka/issues/1273)) ([1fbf2a2](https://github.com/grafana/tanka/commit/1fbf2a2d78f3431cda5a0e8f7e1f743c894fb851))
* **deps:** update dependency kubectl to v1.32.0 ([#1280](https://github.com/grafana/tanka/issues/1280)) ([feac755](https://github.com/grafana/tanka/commit/feac755c7e353136f92bd345b8bfa4ba24205128))


### ♻️ Code Refactoring

* define jsonnet-implementation flag in a single place ([#1260](https://github.com/grafana/tanka/issues/1260)) ([d30882b](https://github.com/grafana/tanka/commit/d30882bcd28e77976f7028365851f4911eeb5a1e))

## [0.30.2](https://github.com/grafana/tanka/compare/v0.30.1...v0.30.2) (2024-12-02)


### 🏗️ Build System

* **deps:** bump docker/build-push-action from 6.9.0 to 6.10.0 ([#1255](https://github.com/grafana/tanka/issues/1255)) ([2118b15](https://github.com/grafana/tanka/commit/2118b153ed24ae5440c89fa552a5430d2f6741df))
* **deps:** bump github.com/vektah/gqlparser/v2 from 2.5.19 to 2.5.20 ([#1253](https://github.com/grafana/tanka/issues/1253)) ([9b531c9](https://github.com/grafana/tanka/commit/9b531c96c095038593765af187bc8230f6d45a70))
* **deps:** bump JamesIves/github-pages-deploy-action from 4.6.9 to 4.7.1 ([#1256](https://github.com/grafana/tanka/issues/1256)) ([b3dc764](https://github.com/grafana/tanka/commit/b3dc7645e9ac136a4dc0e270ecb2f2af9978fb99))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1254](https://github.com/grafana/tanka/issues/1254)) ([b9db0ce](https://github.com/grafana/tanka/commit/b9db0ce6d1a6d480534c5e7a4ee2509d29ecb0c9))


### 🔧 Miscellaneous Chores

* update location from where kubectl is downloaded from ([#1257](https://github.com/grafana/tanka/issues/1257)) ([c18e134](https://github.com/grafana/tanka/commit/c18e134590d09e331f2cba3467d6399a005f1bc9))

## [0.30.1](https://github.com/grafana/tanka/compare/v0.30.0...v0.30.1) (2024-11-26)


### 🐛 Bug Fixes

* handle quotes in --tla-str and --ext-str in "tk eval" ([#1237](https://github.com/grafana/tanka/issues/1237)) ([7cba21d](https://github.com/grafana/tanka/commit/7cba21d3ea83b20f359516b7dc2e91424c8f48da))


### 🏗️ Build System

* **deps:** bump docker/metadata-action from 5.5.1 to 5.6.1 ([#1245](https://github.com/grafana/tanka/issues/1245)) ([e16af88](https://github.com/grafana/tanka/commit/e16af885811ec4302fbce34223529f4907357dd0))
* **deps:** bump github.com/99designs/gqlgen from 0.17.56 to 0.17.57 in /dagger ([#1244](https://github.com/grafana/tanka/issues/1244)) ([c03cb00](https://github.com/grafana/tanka/commit/c03cb0098c1079f72ae56d5a69cf3160c5bdef48))
* **deps:** bump github.com/stretchr/testify from 1.9.0 to 1.10.0 ([#1243](https://github.com/grafana/tanka/issues/1243)) ([ec8ec69](https://github.com/grafana/tanka/commit/ec8ec690057666ee6aeef8d67236d0e9e450d44f))
* **deps:** bump k8s.io/apimachinery from 0.31.2 to 0.31.3 ([#1242](https://github.com/grafana/tanka/issues/1242)) ([42663ac](https://github.com/grafana/tanka/commit/42663acd25ddc43a9b1e9ad6028ed9318663f86a))
* **deps:** bump the acceptance-tests-dependencies group with 3 updates ([#1241](https://github.com/grafana/tanka/issues/1241)) ([7fdc5e1](https://github.com/grafana/tanka/commit/7fdc5e16123ff89871e076a0dffd34f815af7c73))
* **deps:** bump the docs-dependencies group in /docs with 3 updates ([#1240](https://github.com/grafana/tanka/issues/1240)) ([48e2c12](https://github.com/grafana/tanka/commit/48e2c121b2f7b6a720a1e9a9246b6762749b7ec2))


### 🤖 Continuous Integration

* create release docker image through workflow-call ([#1246](https://github.com/grafana/tanka/issues/1246)) ([fb6380f](https://github.com/grafana/tanka/commit/fb6380fb0e46bef6ab1657e9f3bdeeea3997aa35))
* inject dockerfile dependency versions from workflow ([#1247](https://github.com/grafana/tanka/issues/1247)) ([eb9aac0](https://github.com/grafana/tanka/commit/eb9aac0f9fa393d9b010d68289cc3981cc7c5a1f))
* relevant workflows should react also to ready_for_review ([#1248](https://github.com/grafana/tanka/issues/1248)) ([3183efa](https://github.com/grafana/tanka/commit/3183efa8cb38abc13db1d8a8e4d585f3d8d6c1fb))
* run lint-pr-title workflow on ready_for_review ([#1249](https://github.com/grafana/tanka/issues/1249)) ([45c822e](https://github.com/grafana/tanka/commit/45c822ed39bef4a50dfc024d6507b5096dd42f71))
* run lint-pr-title workflow on ready-for-review ([45c822e](https://github.com/grafana/tanka/commit/45c822ed39bef4a50dfc024d6507b5096dd42f71))

## [0.30.0](https://github.com/grafana/tanka/compare/v0.29.0...v0.30.0) (2024-11-22)


### 🎉 Features

* new `tk tool importers-count` ([#1232](https://github.com/grafana/tanka/issues/1232)) ([5dcb6c5](https://github.com/grafana/tanka/commit/5dcb6c5bb56a704ed806c74d37beede93352415a))


### 🐛 Bug Fixes

* delete command supports kinds that do not match singular/plural names ([#1236](https://github.com/grafana/tanka/issues/1236)) ([bf702ef](https://github.com/grafana/tanka/commit/bf702ef1fb562be8f1e998e0468dd7e178f20cac))


### 🏗️ Build System

* **deps:** bump actions/checkout from 4.2.0 to 4.2.2 ([#1233](https://github.com/grafana/tanka/issues/1233)) ([cb5d0c8](https://github.com/grafana/tanka/commit/cb5d0c8cafc0d7b855e4210f0e64a205e3cec163))
* **deps:** bump dagger/dagger-for-github from 6.14.0 to 7.0.1 ([#1224](https://github.com/grafana/tanka/issues/1224)) ([932e7ce](https://github.com/grafana/tanka/commit/932e7ceecf85aa15b9f855eec47d62c4e8cd254f))
* **deps:** bump golang from 1.23.2 to 1.23.3 ([#1227](https://github.com/grafana/tanka/issues/1227)) ([576bfe5](https://github.com/grafana/tanka/commit/576bfe561c39c783cfe6e1b92918d5483954a060))
* **deps:** bump golang.org/x/term from 0.25.0 to 0.26.0 ([#1223](https://github.com/grafana/tanka/issues/1223)) ([70d96f8](https://github.com/grafana/tanka/commit/70d96f8c1590563aa33a9ce054d484f22e150cc5))
* **deps:** bump golang.org/x/text from 0.19.0 to 0.20.0 ([#1222](https://github.com/grafana/tanka/issues/1222)) ([b8119ad](https://github.com/grafana/tanka/commit/b8119adc6a5a893d755e01653682dc7bec754f32))
* **deps:** bump JamesIves/github-pages-deploy-action from 4.6.8 to 4.6.9 ([#1225](https://github.com/grafana/tanka/issues/1225)) ([e754340](https://github.com/grafana/tanka/commit/e754340c4fc6a8ed7b0f12830b1b573bf591f9b4))
* **deps:** bump the dagger-dependencies group in /dagger with 3 updates ([#1235](https://github.com/grafana/tanka/issues/1235)) ([7023072](https://github.com/grafana/tanka/commit/70230726ab7e42f8716597174f1b8d2ef9eb5bce))
* **deps:** bump the dagger-dependencies group in /dagger with 9 updates ([#1221](https://github.com/grafana/tanka/issues/1221)) ([ef52a66](https://github.com/grafana/tanka/commit/ef52a662b6c7dfb4b833115ad082323566e1eae7))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1226](https://github.com/grafana/tanka/issues/1226)) ([24eeca3](https://github.com/grafana/tanka/commit/24eeca32610d1cf668930558ac970161956d2393))
* **deps:** bump the docs-dependencies group in /docs with 2 updates ([#1234](https://github.com/grafana/tanka/issues/1234)) ([fe13459](https://github.com/grafana/tanka/commit/fe134590c8a0efd25ad2fe916dd7917b5be2d5a5))


### 🤖 Continuous Integration

* prevent breaking change from creating major release &lt; 1.0.0 ([#1231](https://github.com/grafana/tanka/issues/1231)) ([1f34a6e](https://github.com/grafana/tanka/commit/1f34a6e5693af10bf0a0d243b51e87be4d017969))
* add release-please for release-automation ([#1195](https://github.com/grafana/tanka/issues/1195)) ([6918cec](https://github.com/grafana/tanka/commit/6918ceccee590e225deb1466fa202211a8a554a6))


## 0.23.1 (2022-09-28)

### Bug Fixes/Enhancements

- **export**: Fix `getSnippetHash` not considering all files
  **[#765](https://github.com/grafana/tanka/pull/765)**

## 0.23.0 (2022-09-26)

### Features

- **cli/tanka**: Add new `--auto-approve=(always|never|if-no-changes)` option to the `apply` command
  **[#754](https://github.com/grafana/tanka/pull/754)**
- **cli/export**: Expand merging capabilities with new `--merge-strategy` flag
  **[#760](https://github.com/grafana/tanka/pull/760)**

### Bug Fixes/Enhancements

- **cli/tanka**: Use exact match to find context from API server
  **[#750](https://github.com/grafana/tanka/pull/750)**
- **helm**: Handle dirs missing the `Chart.yaml` file
  **[#752](https://github.com/grafana/tanka/pull/752)**
- **export**: `getSnippetHash`: Use regexp instead of parsing whole AST for performance
  **[#758](https://github.com/grafana/tanka/pull/758)**

## 0.22.1 (2022-06-15)

### Bug Fixes/Enhancements

- **helm**: Fix `vendor --prune` deleting charts with a custom directory
  **[#717](https://github.com/grafana/tanka/pull/717)**
- **helm**: Add validation at vendoring time for invalid chart names
  **[#718](https://github.com/grafana/tanka/pull/718)**
- **helm**: Fix cross-device link error when tmp is mounted on a different device
  **[#720](https://github.com/grafana/tanka/pull/720)**

## 0.22.0 (2022-06-03)

### Features

- **cli**: Add lint command
  **[#592](https://github.com/grafana/tanka/pull/592)**
- **cli**: Support a diff-strategy of "none" for "tk apply" to suppress diffing
  **[#700](https://github.com/grafana/tanka/pull/700)** (**jphx**)
- **cli**: Add a fallback to inline environment when path doesn't exist
  **[#637](https://github.com/grafana/tanka/pull/637)** (**josephglanville**)
- **kubectl**: Support interactive diff utilities
  **[#690](https://github.com/grafana/tanka/pull/690)** (**partcyborg**)
- **helm**: Allow defining a custom dir for each chart
  **[#706](https://github.com/grafana/tanka/pull/706)**
- **helm**: Add `--prune` option to the vendor command
  **[#707](https://github.com/grafana/tanka/pull/707)**
- **helm**: Check for output dir conflicts
  **[#710](https://github.com/grafana/tanka/pull/710)**
- **cli/tanka**: Adds support for contextNames in tk env subcommands
  **[#704](https://github.com/grafana/tanka/pull/704)** (**Nashluffy**)

### Bug Fixes/Enhancements

- **helm**: Compare semvers when checking if existing chart is up-to-date
  **[#702](https://github.com/grafana/tanka/pull/702)** (**kklimonda-fn**)
- **tanka**: Omit empty `apiServer` or `contextNames` when listing environments
  **[#709](https://github.com/grafana/tanka/pull/709)**
- **export**: Fix caching in case of missing import
  **[#712](https://github.com/grafana/tanka/pull/712)**
- **helm**: Tighten validations for `add` and `add-repo` commands
  **[#713](https://github.com/grafana/tanka/pull/713)**
- **cli/export**: Un-hide the memory ballast setting
  **[#714](https://github.com/grafana/tanka/pull/714)**

## 0.21.0 (2022-04-28)

### Features

- **cli**: Add Apple Silicon binary
  **[#685](https://github.com/grafana/tanka/pull/685)** (**BeyondEvil**)
- **tanka/cli**: Add server-side apply mode
  **[#651](https://github.com/grafana/tanka/pull/651)** (**smuth4**)
- **tanka**: Adds support for specifying valid context names for an environment
  **[#674](https://github.com/grafana/tanka/pull/674)** (**Nashluffy**)

### Bug Fixes/Enhancements

- **cli**: Remove backticks from -inject-labels flag desc
  **[#688](https://github.com/grafana/tanka/pull/688)** (**colega**)
- **tanka**: Fix target must be a non-nil pointer
  **[#684](https://github.com/grafana/tanka/pull/684)** (**maoueh**)
- **tanka**: Upgrade to Go 1.18 + Upgrade dependencies
  **[#697](https://github.com/grafana/tanka/pull/697)**

## 0.20.0 (2022-02-01)

### Features

- **jsonnet**: Update `go-jsonnet` to version 0.18.0
  **[#660](https://github.com/grafana/tanka/pull/660)**
- **cli**: Add `--dry-run` kubectl option
  **[#667](https://github.com/grafana/tanka/pull/667)**
- **helm**: Add option to pass `--skip-tests`
  **[#654](https://github.com/grafana/tanka/pull/654)** (**jouve**)
- **export**: Introduce a configurable memory ballast
  **[#669](https://github.com/grafana/tanka/pull/669)**

## 0.19.0 (2021-11-22)

### Notice

The go.yaml library's version lock was removed. Sequence items in YAML generated from the `manifestYamlFromJson` native function will have a different indent level.

If you are exporting manifests from multiple environments with `tk export` and you wish to do it gradually, you can do it using the `--selector` argument. Here's an example where environments have a `cluster` label:

```console
// Export the dev cluster with the new version
tk-new export outputs-dir tanka-dir --merge --selector cluster=dev
// Export other clusters with the old version
tk export outputs-dir tanka-dir --merge --selector cluster!=dev
```

### Bug Fixes

- **helm**: match `Add()` and `AddRepos()` and correct typos
  **[#641](https://github.com/grafana/tanka/pull/641)** (**redradrat**)
- **yaml**: Remove yaml.v3's version lock
  **[#643](https://github.com/grafana/tanka/pull/643)**

## 0.18.2 (2021-10-14)

### Features

- **cli**: Add `--max-stack` jsonnet option
  **[#619](https://github.com/grafana/tanka/pull/619)**

### Bug Fixes/Enhancements

- **cli**: If there's a full match on an inline environment name, use it
  **[#620](https://github.com/grafana/tanka/pull/620)**
- **cli**: Add instructions to use `--name` on multiple envs error
  **[#621](https://github.com/grafana/tanka/pull/621)**
- **export**: Remove unnecessary `os.Stat` in eval cache
  **[#624](https://github.com/grafana/tanka/pull/624)**
- **tanka**: Upgrade to Go 1.17
  **[#625](https://github.com/grafana/tanka/pull/625)**
- **jsonnet**: Fix `std.thisFile`
  **[#626](https://github.com/grafana/tanka/pull/626)**

## 0.18.1 (2021-10-04)

### Bug Fixes

- **kubernetes**: Fix api-resources table parsing
  **[#605](https://github.com/grafana/tanka/pull/605)**
- **yaml**: Revert yaml.v3 bump due to changes to indent
  **[#616](https://github.com/grafana/tanka/pull/616)**

## 0.18.0 (2021-10-01)

### Features

- **export**: Implement environment caching for `tk export`
  **[#603](https://github.com/grafana/tanka/pull/603)**
- **cli**: Allow partial matches in the --name option
  **[#613](https://github.com/grafana/tanka/pull/613)**

### Bug Fixes

- **tanka**: Check executable prefix before calling stat
  **[#601](https://github.com/grafana/tanka/pull/601)**  (**neerolyte**)
- **cli**: Add hint to inline environment error
  **[#606](https://github.com/grafana/tanka/pull/606)**
- **cli**: Bump cli to `0.2.0`
  **[#611](https://github.com/grafana/tanka/pull/611)**
- **cli**: Add check to prevent using `spec.json` and inline envs simultaneously
  **[#614](https://github.com/grafana/tanka/pull/614)**

## 0.17.3 (2021-08-16)

### Features

- **docker**: Add Kustomize binary
  **([#597](https://github.com/grafana/tanka/pull/597))**

## 0.17.2 (2021-08-10)

### Bug Fixes

- **export**: Add more context when extract fails
  **([#583](https://github.com/grafana/tanka/pull/583))**
- **tanka**: Do not remove `data`, hide it
  **([#585](https://github.com/grafana/tanka/pull/585))**

## 0.17.1 (2021-07-08)

### Bug Fixes

- **cli**: Preserve compatibility for `tk init --k8s=false`
  **([#582](https://github.com/grafana/tanka/pull/582))** (**harmjanblok**)

## 0.17.0 (2021-07-02)

:tada: Big shout out to the community in this release, well done!

### Features

- **helm**: Add support to specify `--kube-version`
  **([#578](https://github.com/grafana/tanka/pull/578))** (**@olegmayko**)
- **kubectl**: Add "validate" diff strategy with `kubectl diff --server-side`
  **([#538](https://github.com/grafana/tanka/pull/538))**

### Bug Fixes

- **helm**: Pass multiple `--api-versions` flags
  **([#576](https://github.com/grafana/tanka/pull/576))** (**@jtdoepke**)
- **jsonnet**: Handle TLA code properly
  **([#574](https://github.com/grafana/tanka/pull/574))** (**@mihaitodor**)
- **export**: Make `--format` respect "/" in template actions
  **([#572](https://github.com/grafana/tanka/pull/572))** (**@dewe**)


## 0.16.0 (2021-06-01)

#### :sparkles: Tanka now defaults to [`k8s-alpha`](https://github.com/jsonnet-libs/k8s-alpha)

`tk init` will now install `k8s-alpha` as the default library for `k.libsonnet`. It is currently defaults to Kubernetes v1.20 however you can pick your own version or disable it:

```console
tk init --k8s=1.18
tk init --k8s=false
```

### Features

- **cli** :sparkles:: `tk init` now defaults to [`k8s-alpha`](https://github.com/jsonnet-libs/k8s-alpha)
  **([#563](https://github.com/grafana/tanka/pull/563))**

### Bug Fixes

- **kubernetes**: Remove resources with altered state
  **([#539](https://github.com/grafana/tanka/pull/539))** (**@StevenPG**)


## 0.15.1 (2021-04-27)

### Features

- **helm**: Add support for `--no-hooks` switch in Helm template
  **([#545](https://github.com/grafana/tanka/pull/545))** (**@PatTheSilent**)
- **export**: Only call FindEnvs once
  **([#553](https://github.com/grafana/tanka/pull/553))**

### Bug Fixes

- **kubernetes**: Don't fail on listing namespaces
  **([#549](https://github.com/grafana/tanka/pull/549))**

## 0.15 (2021-03-22)

Half the changes introduced in this version come from the community, great job y'all!

#### :warning: Pruning label changed

With enabling pruning on inline environments
([#511](https://github.com/grafana/tanka/pull/511)) we fixed pruning. Instead of just
setting the prune label `tanka.dev/environment` to the environment name, Tanka now sets it
to a hash of the environment name and path it is on.

This solves 2 problems:

* Ensures pruning works properly on inline environments.
* Label values in Kubernetes have a 63 characters limit, environment names can be longer.

The effect of this is that all environments using `spec.injectLabels` will show a diff
on the `tanka.dev/environment` label. To ensure a proper migration, execute `tk apply` and
`tk prune` with Tanka v0.14 before running v0.15 so all stale objects are pruned before
the label changes.

Thanks **@craigfurman** for pulling this together.

### Features

- **cli**: Add  `tk env add|set --inject-labels` flag
  **([#505](https://github.com/grafana/tanka/pull/505))** (**@curusarn**)
- **cli**: Add `tk diff --exit-zero` flag
  **([#506](https://github.com/grafana/tanka/pull/506))** (**@craigfurman**)
- **cli**: `tk env list` sorts environments by name
  **([#521](https://github.com/grafana/tanka/pull/521))**
- **cli**: Pruning warns before deleting namespaces
  **([#531](https://github.com/grafana/tanka/pull/531))**
- **cli**: Add `tk status --name` flag and sort Spec.data
  **([#533](https://github.com/grafana/tanka/pull/533))**

* **tooling**: `tk tool imports` works on both files and paths
  **([#517](https://github.com/grafana/tanka/pull/517))**
* **kubernetes**: support .metadata.generateName
  **([#529](https://github.com/grafana/tanka/pull/529))** (**@wojciechka**)
* **helm**: Only update helm repositories when necessary
  **([#535](https://github.com/grafana/tanka/pull/535))** (**@craigfurman**)

### Bug Fixes

- **cli** :sparkles:: Enable pruning on inline environments
  **([#511](https://github.com/grafana/tanka/pull/511))** (**@craigfurman**)
- **cli**: Do not silently fail on find/List
  **([#515](https://github.com/grafana/tanka/pull/515))**
- **cli**: Split diff and non-diff output
  **([#537](https://github.com/grafana/tanka/pull/537))** (**@craigfurman**)

* **tooling**: `tk tool imports` shows path info to error message
  **([#518](https://github.com/grafana/tanka/pull/518))**
* **jsonnet**: TLA in export panic
  **([#519](https://github.com/grafana/tanka/pull/519))** (**@morlay**)

## 0.14 (2021-02-03)

#### :building_construction: Multiple inline environments

As a next step in the inline environment area, this release supports multiple inline
environments. In case there are multiple environments in your workflow, you can use
`--name` to specify the environment you want to diff or apply.

```console
tk apply --name us-central1 environments/dev
tk diff --name europe-west2 environments/prod
```

#### :hammer: Export multiple environments

As part of Grafana Labs' continuous delivery setup, we developed a fast and effective way
to export all our environments. In v0.14, we have built this into `tk export`.

> :warning: breaking change: the arguments for `tk export` have switched places!

`path` to an environment can be added multiple times:

```console
tk export <outputDir> <path> [<path>...] [flags]
```

Some examples:

```bash
# Format based on environment {{env.<...>}}
$ tk export exportDir environments/dev/ --format '{{env.metadata.labels.cluster}}/{{env.spec.namespace}}//{{.kind}}-{{.metadata.name}}'
# Export multiple environments
$ tk export exportDir environments/dev/ environments/qa/
# Recursive export
$ tk export exportDir environments/ --recursive
# Recursive export with labelSelector
$ tk export exportDir environments/ -r -l team=infra
```

### Features

- **tanka** :sparkles:: Handle multiple inline environments
  **([#476](https://github.com/grafana/tanka/pull/476))**
- **jsonnet**: Vendor jsonnet v0.17.0
  **([#445](https://github.com/grafana/tanka/pull/445))**

* **cli**: Extend Tanka with scripts through `tk-` prefix on PATH
  **([#412](https://github.com/grafana/tanka/pull/412))**
* **cli** :sparkles:: Export multiple environments with a single `tk export` command
  **([#450](https://github.com/grafana/tanka/pull/450))**
* **cli**: Initialize inline environments
  **([#451](https://github.com/grafana/tanka/pull/451))**
* **cli**: Add Helm Chart repositories with `tk tool charts add-repo`
  **([#455](https://github.com/grafana/tanka/pull/455))**
* **cli**: Add `--with-prune` option for `tk diff`
  **([#469](https://github.com/grafana/tanka/pull/469))** (**@curusarn**)

- **api**: `Loader` interface
  **([#459](https://github.com/grafana/tanka/pull/459),
  [#467](https://github.com/grafana/tanka/pull/467))**
- **api**: Faster environment discovery and faster `tk env list`
  **([#468](https://github.com/grafana/tanka/pull/468))**

### Bug Fixes

- **jpath**: Support nested calling again
  **([#456](https://github.com/grafana/tanka/pull/456))**
- **cli**: Ensure TLACode works with `EvalScript`
  **([#464](https://github.com/grafana/tanka/pull/464))**
- **jsonnet**: Restore tk.env
  **([#482](https://github.com/grafana/tanka/pull/482),
  [#498](https://github.com/grafana/tanka/pull/498))**

### BREAKING

- **cli**: The argument order of `tk export` changed due to 
  **[#450](https://github.com/grafana/tanka/pull/450)**:

```console
# old:
$ tk export <environment> <outputDir>

# new:
$ tk export <outputDir> <environment> [<environment...>]
```

## 0.13 (2020-12-11)

#### :building_construction: Inline environments

One of the most debated features of the past months has landed: defining Tanka
Environments inline. It is now possible to leverage all powerful Jsonnet
concepts to modify your Tanka Environments. See
https://tanka.dev/inline-environments for more details.

#### :wheel_of_dharma: Kustomize support

In 0.12 we brought in [Helm support](#wheel_of_dharma-helm-support), similarly
Tanka now also comes with Kustomize support.

* We have refactored `helm-util` into
  [`tanka-util`](https://github.com/grafana/jsonnet-libs/blob/master/tanka-util)
  to support both `helm.template()` and `kustomize.build()` use cases from a
  common base.
* Jsonnet-native overwriting of Kustomizations

Have a look at https://tanka.dev/kustomize on how to use all this goodness.
Also https://tanka.dev/helm has been updated to match the library changes.

> This feature is currently experimental. We believe it is feature complete, but
> further usage in the field may lead to adjustments

#### :sparkles: Export `JSONNET_PATH`

Tanka :heart: Jsonnet, and so do you. With this release, you can now access the
`JSONNET_PATH` that Tanka uses to find all libraries. Try something this in your
environment:

```console
$ JSONNET_PATH=$(tk tool jpath environments/prometheus) jsonnet-lint environments/prometheus/main.jsonnet
```

#### :speech_balloon: Github Discussions

The Tanka project is trying out GitHub Discussions as the primary support channel:

- :mag: It is searchable, so information never gets lost in Slack again
- :busts_in_silhouette: No longer sign-up for two platforms
- :mega: Reach both, the team and other community members

Head over to https://github.com/grafana/tanka/discussions and start the discussion!

### Features

- **jsonnet**: Allow alternative entrypoints
  **([#389](https://github.com/grafana/tanka/pull/389))**
- **jsonnet** :sparkles:: Support for inline environment
  **([#403](https://github.com/grafana/tanka/pull/403))**
- **jsonnet, cli** :sparkles:: `tk tool jpath` can be used to export `JSONNET_PATH`
  **([#427](https://github.com/grafana/tanka/pull/427))**
- **jsonnet, docker**: Add `openssh-client` to Docker image
  **([#429](https://github.com/grafana/tanka/pull/429))** (**@xvzf**)

* **k8s** :sparkles:: Render Kustomize into Jsonnet
  **([#422](https://github.com/grafana/tanka/pull/422))**
* **k8s**: Add `metadata.Namespace` directive, always the path relative to the project root
  **([#435](https://github.com/grafana/tanka/pull/435))**

- **helm**: Chart tool: Check chart versions and update accordingly
  **([#420](https://github.com/grafana/tanka/pull/420))** (**@craigfurman**)
- **helm,docker**: Add `helm` client to Docker image
  **([#430](https://github.com/grafana/tanka/pull/430))** (**@ducharmemp**)

* **api**: Introduce the concept of Evaluators
  **([#431](https://github.com/grafana/tanka/pull/431))**


### Bug Fixes

- **helm**: Chart tool: Detect already pulled charts
  **([#402](https://github.com/grafana/tanka/pull/402))** (**justinwalz**)
- **helm**: Chart tool: Use new URL for stable helm repo
  **([#425](https://github.com/grafana/tanka/pull/425))**

* **cli**: Environment path as name relative to the project root
  **([#404](https://github.com/grafana/tanka/pull/404))** (**mwasilew2**)
* **cli**: Normalize `tk fmt` paths on Windows
  **([#411](https://github.com/grafana/tanka/pull/411))** (**nlowe**)
* **cli**: Confirmation prompts on Windows
  **([#413](https://github.com/grafana/tanka/pull/413))** (**nlowe**)


## 0.12 (2020-10-05)

Like good wine, some things need time. After 3 months of intense development we
have another Tanka release ready:

#### :wheel_of_dharma: Helm support

This one is huge! Tanka can now **load Helm Charts**:

- [`helm-util`](https://github.com/grafana/jsonnet-libs/tree/master/helm-util)
  provides `helm.template()` to load them from inside Jsonnet
- Declarative vendoring using `tk tool charts`
- Jsonnet-native overwriting of chart contents

Just by upgrading to 0.12, you have access to every single Helm chart on the
planet, right inside of Tanka! Read more on https://tanka.dev/helm

> This feature is currently experimental. We believe it is feature complete, but
> further usage in the field may lead to adjustments

#### :house: Top Level Arguments

Tanka now supports the `--tla-str` and `--tla-code` flags from the `jsonnet` cli
to late-bind data into the evaluation in a well-defined way. See
https://tanka.dev/jsonnet/injecting-values for more details.

#### :sparkles: Inline Eval

Ever wanted to pull another value out of Jsonnet that does not comply to the
Kubernetes object rules Tanka imposes onto everything? Wait no longer and use
`tk eval -e`:

```console
$ tk eval environments/prometheus -e prometheus_rules
```

Above returns `$.prometheus_rules` as JSON. Every Jsonnet selector is supported:

```console
$ tk eval environments/prometheus -e 'foo.bar[0]'
```

### Features

- **k8s, jsonnet** :sparkles:: Support for [Helm](https://helm.sh). In combination with
  [`helm-util`](https://github.com/grafana/jsonnet-libs/tree/master/helm-util),
  Tanka can now load resources from Helm Charts.
  **([#336](https://github.com/grafana/tanka/pull/336))**
- **k8s**: Default metadata from `spec.json`
  **([#366](https://github.com/grafana/tanka/pull/366))**

* **helm**: Charttool: Adds `tk tool charts` for easy management of vendored
  Helm charts **([#367](https://github.com/grafana/tanka/pull/367))**,
  **([#369](https://github.com/grafana/tanka/pull/369))**
* **helm**: Require Helm Charts to be available locally
  **([#370](https://github.com/grafana/tanka/pull/370))**
* **helm**: Configurable name format
  **([#381](https://github.com/grafana/tanka/pull/381))**

- **cli**: Filtering (`-t`) now supports negative expressions (`-t !deployment/.*`) to exclude resources
  **([#339](https://github.com/grafana/tanka/pull/339))**
- **cli** :sparkles:: Inline eval (Use `tk eval -e` to extract nested fields)
  **([#378](https://github.com/grafana/tanka/pull/378))**
- **cli**: Custom paging (`PAGER` env var)
  **([#373](https://github.com/grafana/tanka/pull/373))**
- **cli**: Predict plain directories if outside a project
  **([#357](https://github.com/grafana/tanka/pull/357))**

* **jsonnet** :sparkles:: Top Level Arguments can now be specified using `--tla-str` and
  `--tla-code` **([#340](https://github.com/grafana/tanka/pull/340))**

### Bug Fixes

- **yaml**: Pin yaml library to v2.2.8 to avoid whitespace changes
  **([#386](https://github.com/grafana/tanka/pull/386))**
- **cli**: Actually respect `TANKA_JB_PATH`
  **([#350](https://github.com/grafana/tanka/pull/350))**
- **k8s**: Update `kubectl v1.18.0` warning
  **([#371](https://github.com/grafana/tanka/pull/371))**

* **jsonnet**: Load `main.jsonnet` using full path. This makes `std.thisFile`
  usable **([#370](https://github.com/grafana/tanka/pull/370))**
* **jsonnet**: Import path resolution now works on Windows
  **([#331](https://github.com/grafana/tanka/pull/331))**
* **jsonnet**: Arrays are now supported at the top level
  **([#321](https://github.com/grafana/tanka/pull/321))**

### BREAKING

- **api**: Struct based Go API: Modifies our Go API
  (`github.com/grafana/tanka/pkg/tanka`) to be based on structs instead of
  variadic arguments. This has no impact on daily usage of Tanka.
  **([#376](https://github.com/grafana/tanka/pull/376))**
- **jsonnet**: ExtVar flags are now `--ext-str` and `--ext-code` (were `--extVar` and `--extCode`)
  **([#340](https://github.com/grafana/tanka/pull/340))**

## 0.11.1 (2020-07-17)

This is a minor release with one bugfix and one minor feature.

### Features

- **process**: With 0.11.0, tanka started automatically adding namespaces to _all_ manifests it processed. We updated this to _not_
  add a namespace to cluster-wide object types in order to make handling of these resources more consistent in different workflows. **([#320](https://github.com/grafana/tanka/pull/320))**

### Bug Fixes

- **export**: Fix inverted logic while checking if a file already exists. This broke `tk export` entirely.
  **([#317](https://github.com/grafana/tanka/pull/317))**

## 0.11.0 (2020-07-07)

2 months later and here we are with another release! Packed with many
detail-improvements, this is what we want to highlight:

#### :sparkles: Enhanced Kubernetes resource handling

From now on, Tanka handles the resources it extracts from your Jsonnet output in
an enhanced way:

1. **Lists**: Contents of lists, such as `RoleBindingList` are automatically
   flattened into the resource stream Tanka works with. This makes sure they are
   properly labeled for garbage collection, etc.
2. **Default namespaces**: While you could always define the default namespace
   (the one for resources without an explicit one) in `spec.json`, this
   information is now also persisted into the YAML returned by `tk show` and `tk export`.
   See https://tanka.dev/namespaces for more information.

#### :hammer: More powerful exporting

`tk export` can now do even more than just writing YAML files to disk:

1. `--extension` can be used to control the file-extension (defaults to `.yaml`)
2. When you put a `/` in your `--format` for the filename, Tanka creates a
   directory or you. This allows e.g. sorting by namespace:
   `--format='{{.metadata.namespace}}/{{.kind}}-{{.metadata.name}}'`
3. Using `--merge`, you can export multiple environments into the same directory
   tree, so you get the full YAML picture of your entire cluster!

#### :fax: Easier shell scripting

The `tk env list` command now has a `--names` option making it easy to operate on multiple environments:

```bash
# diff all environments:
for e in $(tk env list --names); do
  tk diff $e;
done
```

Also, to use a more granular subset of your environments, you can now use
`--selector` / `-l` to match against `metadata.labels` of defined in your
`spec.json`:

```bash
$ tk env list -l status=dev
```

### Features

- **cli**: `tk env list` now supports label selectors, similar to `kubectl get -l` **([#295](https://github.com/grafana/tanka/pull/295))**
- **cli**: If `spec.apiServer` of `spec.json` lacks a protocol, it now defaults
  to `https` **([#289](https://github.com/grafana/tanka/pull/289))**
- **cli**: `tk delete` command to teardown environments
  **([#313](https://github.com/grafana/tanka/pull/313))**

* **cli**: Support different file-extensions than `.yaml` for `tk export`
  **([#294](https://github.com/grafana/tanka/pull/394))** (**@marthjod**)
* **cli**: Support creating sub-directories in `tk export`
  **([#300](https://github.com/grafana/tanka/pull/300))** (**@simonfrey**)
* **cli**: Allow writing into existing folders during `tk export`
  **([#314](https://github.com/grafana/tanka/pull/314))**

- **tooling**: `tk tool imports` now follows symbolic links
  **([#302](https://github.com/grafana/tanka/pull/302))**,
  **([#303](https://github.com/grafana/tanka/pull/303))**

* **process**: `List` types are now unwrapped by Tanka itself
  **([#306](https://github.com/grafana/tanka/pull/306))**
* **process**: Automatically set `metadata.namespace` to the value of
  `spec.namespace` if not set from Jsonnet
  **([#312](https://github.com/grafana/tanka/pull/312))**

### Bug Fixes

- **jsonnet**: Using `import "tk"` twice no longer panics
  **([#290](https://github.com/grafana/tanka/pull/290))**
- **tooling**: `tk tool imports` no longer gets stuck when imports are recursive
  **([#298](https://github.com/grafana/tanka/pull/298))**
- **process**: Fully deterministic recursion, so that error messages are
  consistent **([#307](https://github.com/grafana/tanka/pull/307))**

## 0.10.0 (2020-05-07)

New month, new release! And this one ships with a long awaited feature:

#### :sparkles: Garbage collection

Tanka can finally clean up behind itself. By optionally attaching a
`tanka.dev/environment` label to each resource it creates, we can find these
afterwards and purge those removed from the Jsonnet code. No more dangling
resources!

> :warning: Keep in mind this is still experimental!

To get started, enable labeling in your environment's `spec.json`:

```diff
  "spec": {
+   "injectLabels": true,
  }
```

Don't forget to `tk apply` afterwards! From now on, Tanka can clean up using `tk prune`.

Docs: https://tanka.dev/garbage-collection

#### :boat: Logo

Tanka now has it's very own logo, and here it is:

<img src="docs/img/logo.svg" width="400px" />

#### :package: Package managers

Tanka is now present in some package managers, notably `brew` for macOS and the
AUR of ArchLinux! See the updated [install
instructions](https://tanka.dev/install#using-a-package-manager-recommended) to
make sure to use these if possible.

### Features:

- **cli**: `TANKA_JB_PATH` environment variable introduced to set the `jb`
  binary if required **([#272](https://github.com/grafana/tanka/pull/272))**.
  Thanks [@qckzr](https://github.com/qckzr)

* **kubernetes**: Garbage collection
  **([#251](https://github.com/grafana/tanka/pull/251))**

### Bug Fixes

- **kubernetes**: Resource sorting is now deterministic
  **([#259](https://github.com/grafana/tanka/pull/259))**

## 0.9.2 (2020-04-19)

Mini-release to fix an issue with our Makefile (required for packaging). No
changes in functionality.

### Bug Fixes

- **build**: Enable `static` Makefile target on all operating systems
  ([#262](https://github.com/grafana/tanka/pull/262))

## 0.9.1 (2020-04-08)

Small patch release to fix a `panic` issue with `tk apply`.

### Bug Fixes

- **kubernetes**: don't panic on failed diff
  **([#256](https://github.com/grafana/tanka/pull/256))**

## 0.9.0 (2020-04-07)

**This release includes a critical fix, update ASAP**.

Another Tanka release is here, just in time for Easter. Enjoy the built-in
[formatter](#sparkles-highlight-jsonnet-formatter-tk-fmt), much [more
intelligent apply](#rocket-highlight-sorting-during-apply) and several important
bug fixes.

#### :rotating_light: Alert: `kubectl diff` changes resources :rotating_light:

The recently released `kubectl` version `v1.18.0` includes a **critical issue**
that causes `kubectl diff` (and so `tk diff` as well) to **apply** the changes.

This can be very **harmful**, so Tanka decided to require you to **downgrade**
to `v1.17.x`, until the fix in `kubectl` version `v1.18.1` is released.

- Upstream issue: https://github.com/kubernetes/kubernetes/issues/89762)
- Unreleased fix: https://github.com/kubernetes/kubernetes/pull/89795

#### :sparkles: Highlight: Jsonnet formatter (`tk fmt`)

Since `jsonnetfmt` was [rewritten in Go
recently](https://github.com/google/go-jsonnet/pull/388), Tanka now ships it as
`tk fmt`. Just run `tk fmt .` to keep all Jsonnet files recursively formatted.

#### :rocket: Highlight: Sorting during apply

When using `tk apply`, Tanka now automatically **sorts** your objects
based on **dependencies** between them, so that for example
`CustomResourceDefinitions` created before being used, all in the same run. No
more partly failed applies!

### Features

- **kubernetes** :sparkles:: Objects are now sorted by dependency before `apply`
  **([#244](https://github.com/grafana/tanka/pull/244))**
- **cli**: Env var `TANKA_KUBECTL_PATH` can now be used to set a custom
  `kubectl` binary
  **([#221](https://github.com/grafana/tanka/pull/221))**
- **jsonnet** :sparkles: : Bundle `jsonnetfmt` as `tk fmt`
  **([#241](https://github.com/grafana/tanka/pull/241))**

* **docker**: The Docker image now includes GNU `less`, instead of the BusyBox
  one **([#232](https://github.com/grafana/tanka/pull/232))**
* **docker**: Added `kubectl`, `jsonnet-bundler`, `coreutils`, `git` and
  `diffutils` to the Docker image, so Tanka can be fully used in there.
  **([#243](https://github.com/grafana/tanka/pull/243))**

### Bug Fixes

- **cli**: The diff shown on `tk apply` is now colored again
  **([#216](https://github.com/grafana/tanka/pull/216))**

* **client**: The namespace patch file saved to a temporary location is now
  removed after run **([#225](https://github.com/grafana/tanka/pull/225))**
* **client**: Scanning for the correct context won't panic anymore, but print a
  proper error **([#228](https://github.com/grafana/tanka/pull/228))**
* **client**: Use `os.PathListSeparator` during context patching, so that Tanka
  also works on non-UNIX platforms (e.g. Windows)
  **([#242](https://github.com/grafana/tanka/pull/242))**

- **kubernetes** :rotating_light:: Refuse to diff on `kubectl` version `v1.18.0`, because of
  above mentioned unfixed issue
  **([#254](https://github.com/grafana/tanka/pull/254))**
- **kubernetes**: Apply no longer aborts when diff fails
  **([#231](https://github.com/grafana/tanka/pull/231))**
- **kubernetes** :sparkles:: Namespaces that will be created in the same run are now
  properly handled during `diff`
  **([#237](https://github.com/grafana/tanka/pull/237))**

### Other

- **cli**: Migrates from `spf13/cobra` to much smaller `go-clix/cli`. This cuts
  our dependencies to a minimum.
  **([#235](https://github.com/grafana/tanka/pull/235))**

## 0.8.0 (2020-02-13)
 (**@xvzf**)
The next big one is here! Feature packed with environment overriding and `tk export`. Furthermore lots of bugs were fixed, so using Tanka should be much
smoother now!

#### Highlight: Overriding `vendor` per Environment **([#198](https://github.com/grafana/tanka/pull/198))**

It is now possible, to have a `vendor/` directory managed by `jb` on an
environment basis: https://tanka.dev/libraries/overriding. This means you can
test out changes in libraries in single environments (like `dev`), without
affecting others (like `prod`).

#### Notice:

Changes done in the last release (v0.7.1) can cause indentation changes when
using `std.manifestYAMLFromJSON()`, related to bumping `gopkg.in/yaml.v2` to
`gopkg.in/yaml.v3`.  
Please encourage all your teammembers to upgrade to at least v0.7.1 to avoid
whitespace-only diffs on your projects.

### Features

- **cli**: `tk export` can be used to write all generated Kubernetes resources
  to `.yaml` files

### Bug Fixes

- **kubernetes**: Fail on `diff` when `kubectl` had an internal error
  **([#213](https://github.com/grafana/tanka/pull/213))**
- **kubernetes**: Stop injecting namespaces into wrong places:  
  Tanka was injecting the default namespace into resources of all kinds,
  regardless of whether they actually took one. This caused errors, so we
  stopped doing this. From now on, the default namespace will only be injected
  when the resource is actually namespaced.
  **([#208](https://github.com/grafana/tanka/pull/208))**

* **cli**: `tk diff` colors:  
  Before, the coloring was unstable when scrolling up and down. We fixed this by
  pressing CAPS-LOCK.  
  Furthermore, the output of `tk diff` now also works on light color schemes,
  without messing up the background color.
  **([#210](https://github.com/grafana/tanka/pull/210))**
* **cli**: Proper `--version` output:  
  The release binaries now show the real semver on `tk --version`, instead of
  the git commit sha. **([#201](https://github.com/grafana/tanka/pull/201))**
* **cli**: Print diff on apply again:  
  While refactoring, we accidentally forgot to dereference a pointer, so that
  `tk apply` showed a memory address instead of the actual differences, which
  was kinda pointless. **([#200](https://github.com/grafana/tanka/pull/200))**

## 0.7.1 (2020-02-06)

This is a smaller release focused on critical bug fixes and some other minor
enhancements. While features are included, none of them are significant, meaning
they are part of a patch release.

#### Critical: `parseYaml` works now

Before, `std.native('parseYaml')` did not work at all, a line of code got lost
during merge/rebase, resulting in `parseYaml` returning invalid data, that
Jsonnet could not process. This issue has been fixed in
**([#195](https://github.com/grafana/tanka/pull/195))**.

#### Jsonnet update

The built-in Jsonnet compiler has been upgraded to the lastest master
[`07fa4c0`](https://github.com/google/go-jsonnet/commit/07fa4c037b4ff8b5e601546cb5de4abecaf2651d).
In some cases, this should provide up to 50% more speed, especially when
`base64` is involved, which is now natively implemented.
**([#196](https://github.com/grafana/tanka/pull/196))**

### Features

- **cli**: `tk env set|add` has been extended by `--server-from-context`, which
  allows to parse `$KUBECONFIG` to find the apiServer's IP directly from that
  file, instead of having to manually specify it by hand.
  **([#184](https://github.com/grafana/tanka/pull/184))**
- **jsonnet**: `vendor` overrides:  
  It is now possible to have a `vendor/` directory per environment, so that
  updating upstream libraries can be done gradually.
  **([#185](https://github.com/grafana/tanka/pull/185))**
- **kubernetes**: disable `kubectl` validation:
  `tk apply` now takes `--validate=false` to pass that exact flag to `kubectl`
  as well, for disabling the integrated schema validation.
  **([#186](https://github.com/grafana/tanka/pull/186))**

### Bug Fixes

- **jsonnet, cli**: Stable environment name: The value of `(import "tk").env.name`
  does not anymore depend on how Tanka was invoked, but will
  always be the relative path from `<rootDir>` to the environment's directory.
  **([#182](https://github.com/grafana/tanka/pull/182))**
- **jsonnet**: The nativeFunc `parseYaml` has been fixed to actually return a
  valid result **([#195](https://github.com/grafana/tanka/pull/195))**

## 0.7.0 (2020-01-21)

The promised big update is here! In the last couple of weeks a lot has happened.

Grafana Labs [announced Tanka to the
public](https://grafana.com/blog/2020/01/09/introducing-tanka-our-way-of-deploying-to-kubernetes/),
and the project got a lot of positive feedback, shown both on HackerNews and in
a 500+ increase in GitHub stars!

While we do not ship big new features this time, we ironed out many annoyances
and made the overall experience a lot better:

#### Better website + tutorial ([#134](https://github.com/grafana/tanka/pull/134))

Our [new website](https://tanka.dev) is published! It does not only look super
sleek and performs like a supercar, we also revisited (and rewrote) the most of
the content, to provide especially new users a good experience.

This especially includes the **[new
tutorial](https://tanka.dev/tutorial/overview)**, which gives new and probably
even more experienced users a good insight into how Tanka is meant to be used.

#### :rotating_light::rotating_light: Disabling `import ".yaml"` ([#176](https://github.com/grafana/tanka/pull/176)) :rotating_light::rotating_light:

Unfortunately, we **had to disable the feature** that allowed to directly import
YAML files using the familiar `import` syntax, introduced in v0.6.0, because it
caused serious issues with `importstr`, which became unusable.

While our extensions to the Jsonnet language are cool, it is a no-brainer that
compatibility with upstream Jsonnet is more important. We will work with the
maintainers of Jsonnet to find a solution to enable both, `importstr` and
`import ".yaml"`

**Workaround:**

```diff
- import "foo.yaml"
+ std.parseYaml(importstr "foo.yaml")
```

#### `k.libsonnet` installation ([#140](https://github.com/grafana/tanka/pull/140))

Previously, installing `k.libsonnet` was no fun. While the library is required
for nearly every Tanka project, it was not possible to install it properly using
`jb`, manual work was required.

From now on, **Tanka automatically takes care of this**. A regular `tk init`
installs everything you need. In case you prefer another solution, disable this
new thing using `tk init --k8s=false`.

### Features

- **cli**, **kubernetes**: `k.libsonnet` is now automatically installed on `tk init` **([#140](https://github.com/grafana/tanka/pull/140))**:  
  Before, installing `k.libsonnet` was a time consuming manual task. Tanka now
  takes care of this, as long as `jb` is present on the `$PATH`. See
  https://tanka.dev/tutorial/k-lib#klibsonnet for more details.
- **cli**: `tk env --server-from-context`:  
  This new flag allows to infer the cluster IP from an already set up `kubectl`
  context. No need to remember IP's anymore – and they are even autocompleted on
  the shell. **([#145](https://github.com/grafana/tanka/pull/145))**
- **cli**, **jsonnet**: extCode, extVar:  
  `-e` / `--extCode` and `--extVar` allow using `std.extVar()` in Tanka as well.
  In general, `-e` is the flag to use, because it correctly handles all Jsonnet
  types (string, int, bool). Strings need quoting!
  **([#178](https://github.com/grafana/tanka/pull/178))**

* **jsonnet**: The contents of `spec.json` are now accessible from Jsonnet using
  `(import "tk").env`. **([#163](https://github.com/grafana/tanka/pull/163))**
* **jsonnet**: Lists (`[ ]`) are now fully supported, at an arbitrary level of
  nesting! **([#166](https://github.com/grafana/tanka/pull/166))**

### Bug Fixes

- **jsonnet**: `nil` values are ignored from the output. This allows to disable
  objects using the `if ... then {}` pattern, which returns nil if `false`
  **([#162](https://github.com/grafana/tanka/pull/162))**.
- **cli**: `-t` / `--target` is now case-insensitive
  **([#130](https://github.com/grafana/tanka/pull/130))**

---

## 0.6.1 (2020-01-06)

First release of the new year! This one is a quick patch that lived on master
for some time, fixing an issue with the recent "missing namespaces" enhancement
leading to `apply` being impossible when no namespace is included in Jsonnet.

More to come soon :D

---

## 0.6.0 (2019-11-27)

It has been quite some time since the last release during which Tanka has become
much more mature, especially regarding the code quality and structure.

Furthermore, Tanka has just hit the 100 Stars :tada:

Notable changes include:

#### API ([#97](https://github.com/grafana/tanka/commit/c5edb8b0153ef991765f2f555c839b0f9a487e75))

The most notable change is probably the **Go API**, available at
`https://godoc.org/github.com/grafana/tanka/pkg/tanka`, which allows to use all
features of Tanka directly from any other Golang application, without needing to
exec the binary. The API is inspired by the command line parameters and should
feel very similar.

#### Importing YAML ([#106](https://github.com/grafana/tanka/commit/8029efa44461b5f7ba83a218ccc45bd758c8a322))

It is now possible to import `.yaml` documents directly from Jsonnet. Just use
the familiar syntax `import "foo.yaml"` like you would with JSON.

#### Missing Namespaces ([#120](https://github.com/grafana/tanka/commit/3b9fac1563a75a571b512887602eb53f82e565bf))

Tanka now handles namespaces that are not yet created, in a more user friendly
way than `kubectl\*\* does natively.  
During diff, all objects of an in-existent namespace are shown as new and when
applying, namespaces are applied first to allow applying in a single step.

### Features

- **tool/imports**: import analysis using upstream jsonnet: Due to recent
  changes to google/jsonnet, we can now use the upstream compiler for static
  import analysis
  ([#84](https://github.com/grafana/tanka/commit/394cb12b28beb0ea05d065594b6cf5c3f92de5e4))
- **Array output**: The output of Jsonnet may now be an array of Manifests.
  Nested arrays are not supported yet.
  ([#112](https://github.com/grafana/tanka/commit/eb647793ff5515bc828e4f91186655c143bb6a04))

### Bug Fixes

- **Command Usage Guidelines**: Tanka now uses the [command description
  syntax](https://en.wikipedia.org/wiki/Command-line_interface#Command_description_syntax)
  ([#94](https://github.com/grafana/tanka/commit/13238e5941bd6e68f410d3938d1a285224c2f91d))
- **cli/env** resolved panic on missing `spec.json`
  ([#108](https://github.com/grafana/tanka/commit/9bd15e6b4226164efe45f50c9ed41c4a5673ea2d))

---

## 0.5.0 (2019-09-20)

This version adds a set of commands to manipulate environments (`tk env add, rm, set, list`) ([#73](https://github.com/grafana/tanka/pull/73)). The commands are
mostly `ks env` compatible, allowing `tk env` be used as a drop-in replacement
in scripts.

Furthermore, an error message has been improved, to make sure users can
differentiate between parse issues in `.jsonnet` and `spec.json`
([#71](https://github.com/grafana/tanka/pull/71)).

---

## 0.4.0 (2019-09-06)

After nearly a month, the next feature packed release of Tanka is ready!
Highlights include the new documentation website https://tanka.dev, regular
expression support for targets, diff histograms and several bug-fixes.

### Features

- **cli**: `tk show` now aborts by default, when invoked in a non-interactive
  session. Use `--dangerous-allow-redirect` to disable this safe-guard
  ([#47](https://github.com/grafana/tanka/issues/47)).
- **kubernetes**: Regexp Targets: It is now possible to use regular expressions
  when specifying the targets using `--target` / `-t`. Use it to easily select
  multiple objects at once: https://tanka.dev/targets/#regular-expressions
  ([#64](https://github.com/grafana/tanka/issues/64)).
- **kubernetes**: Diff histogram: Tanka now allows to summarize the differences
  between the live configuration and the local one, by using the unix
  `diffstat(1)` utility. Gain a sneek peek at a change using `tk diff -s .`!
  ([#67](https://github.com/grafana/tanka/issues/67))

### Bug Fixes

- **kubernetes**: Tanka does not fail anymore, when the configuration file
  `spec.json` is missing from an Environment. While you cannot apply or diff,
  the show operation works totally fine
  ([#56](https://github.com/grafana/tanka/issues/56),
  [#63](https://github.com/grafana/tanka/issues/63)).
- **kubernetes**: Errors from `kubectl` are now correctly passed to the user
  ([#61](https://github.com/grafana/tanka/issues/61)).
- **cli**: `tk diff` does not output useless empty lines (`\n`) anymore
  ([#62](https://github.com/grafana/tanka/issues/62)).

---

## 0.3.0 (2019-08-13)

Tanka v0.3.0 is here!

This version includes lots of tiny fixes and detail improvements, to make it easier for everyone to configure their Kubernetes clusters.

Enjoy target support, enhancements to the diff UX and an improved CLI experience.

### Features

The most important feature is **target support** ([#30](https://github.com/tbraack/tanka/issues/30)) ([caf205a](https://github.com/tbraack/tanka/commit/caf205a)): Using `--target=kind/name`, you can limit your working set to a subset of the objects, e.g. to do a staged rollout.

There where some other features added:

- **cli:** autoApprove, forceApply ([#35](https://github.com/tbraack/tanka/issues/35)) ([626b097](https://github.com/tbraack/tanka/commit/626b097)): allows to skip the interactive verification. Furthermore, `kubectl` can now be invoked with `--force`.
- **cli:** print deprecated warnings in verbose mode. ([#39](https://github.com/tbraack/tanka/issues/39)) ([6de170d](https://github.com/tbraack/tanka/commit/6de170d)): Warnings about the deprecated configs are only printed in verbose mode
- **kubernetes:** add namespace to apply preamble ([#23](https://github.com/tbraack/tanka/issues/23)) ([9e2d927](https://github.com/tbraack/tanka/commit/9e2d927)): The interactive verification now shows the `metadata.namespace` as well.
- **cli:** diff UX enhancements ([#34](https://github.com/tbraack/tanka/issues/34)) ([7602a19](https://github.com/tbraack/tanka/commit/7602a19)): The user experience of the `tk diff` subcommand has been improved:
  - if the output is too long to fit on a single screen, the systems `PAGER` is invoked
  - if differences are found, the exit status is set to `16`.
  - When `tk apply` is invoked, the diff is shown again, to make sure you apply what you want

### Bug Fixes

- **cli:** invalid command being executed twice ([#42](https://github.com/tbraack/tanka/issues/42)) ([28c6898](https://github.com/tbraack/tanka/commit/28c6898)): When the command failed, it was executed twice, due to an error in the error handling of the CLI.
- **cli**: config miss ([#22](https://github.com/tbraack/tanka/issues/22)) ([32bc8a4](https://github.com/tbraack/tanka/commit/32bc8a4)): It was not possible to use the new configuration format, due to an error in the config parsing.
- **cli:** remove datetime from log ([#24](https://github.com/tbraack/tanka/issues/24)) ([1e37b20](https://github.com/tbraack/tanka/commit/1e37b20))
- **kubernetes:** correct diff type on 1.13 ([#31](https://github.com/tbraack/tanka/issues/31)) ([574f946](https://github.com/tbraack/tanka/commit/574f946)): On kubernetes 1.13.0, `subset` was used, although `native` is already supported.
- **kubernetes:** Nil pointer deference in subset diff. ([#36](https://github.com/tbraack/tanka/issues/36)) ([f53c2b5](https://github.com/tbraack/tanka/commit/f53c2b5))
- **kubernetes:** sort during reconcile ([#33](https://github.com/tbraack/tanka/issues/33)) ([ab9c43a](https://github.com/tbraack/tanka/commit/ab9c43a)): The output of the reconcilation phase is now stable in ordering

---

## [0.2.0](https://github.com/tbraack/tanka/compare/v0.1.0...v0.2.0) (2019-08-07)

### Features

- **cli:** Completions ([#7](https://github.com/tbraack/tanka/issues/7)) ([aea3bdf](https://github.com/tbraack/tanka/commit/aea3bdf)): Tanka is now able auto-complete most of the command line arguments and flags. Supported shells are `bash`, `zsh` and `fish`.
- **cmd:** allow the baseDir to be passed as an argument ([#6](https://github.com/tbraack/tanka/issues/6)) ([55adf80](https://github.com/tbraack/tanka/commit/55adf80)), ([#12](https://github.com/tbraack/tanka/issues/12)) ([3248bb9](https://github.com/tbraack/tanka/commit/3248bb9)): `tk` **breaks** with the current behaviour and requires the baseDir / environment to be passed explicitely on the command line, instead of assuming it as `pwd`. This is because it allows more `go`-like UX. It is also very handy for scripts not needing to switch the directory.
- **kubernetes:** subset-diff ([#11](https://github.com/tbraack/tanka/issues/11)) ([13f6fdd](https://github.com/tbraack/tanka/commit/13f6fdd)): `tk diff` support for version below Kubernetes `1.13` is here :tada:! The strategy is called _subset diff_ and effectively compares only the fields already present in the config. This allows the (hopefully) most bloat-free experience possible without server side diff.
- **tooling:** import analysis ([#10](https://github.com/tbraack/tanka/issues/10)) ([ce2b0d3](https://github.com/tbraack/tanka/commit/ce2b0d3)): Adds `tk tool imports`, which allows to list all imports of a single file (even transitive ones). Optionally pass a git commit hash, to check whether any of the changed files is imported, to figure out which environments need to be re-applied.

---

## 0.1.0 (2019-07-31)

This release marks the begin of tanka's history :tada:!

As of now, tanka aims to nearly seemlessly connect to the point where [ksonnet](https://github.com/ksonnet/ksonnet) left.
The current feature-set is basic, but usable: The three main workflow commands are available (`show`, `diff`, `apply`), environments are supported, code-sharing is done using [`jb`](https://github.com/jsonnet-bundler/jsonnet-bundler).

Stay tuned!

### Features

- **kubernetes:** Show ([7c4bee8](https://github.com/tbraack/tanka/commit/7c4bee8)): Equivalent to `ks show`, allows previewing the generated yaml.
- **kubernetes:** Diff ([a959f38](https://github.com/tbraack/tanka/commit/a959f38)): Uses the `kubectl diff` to obtain a sanitized difference betweent the current and the desired state. Requires Kubernetes 1.13+
- **kubernetes:** Apply ([8fcb4c1](https://github.com/tbraack/tanka/commit/8fcb4c1)): Applies the changes to the cluster (like `ks apply`)
- **kubernetes:** Apply approval ([4c6414f](https://github.com/tbraack/tanka/commit/4c6414f)): Requires a typed `yes` to apply, gives the user the chance to verify cluster and context.
- **kubernetes:** Smart context ([2b3fd3c](https://github.com/tbraack/tanka/commit/2b3fd3c)): Infers the correct context from the `spec.json`. Prevents applying the correct config to the wrong cluster.
- Init Command ([ff8857c](https://github.com/tbraack/tanka/commit/ff8857c)): Initializes a new repository with the suggested directory structure.
