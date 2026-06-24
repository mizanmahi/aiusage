# Changelog

## [0.4.2](https://github.com/mizanmahi/aiusage/compare/v0.4.1...v0.4.2) (2026-06-24)


### Bug Fixes

* **cli:** send usage events in batches ([b5d9915](https://github.com/mizanmahi/aiusage/commit/b5d9915ef18c40653a761bc00e0eab6b8f48cdd9))

## [0.4.1](https://github.com/mizanmahi/aiusage/compare/v0.4.0...v0.4.1) (2026-06-24)


### Bug Fixes

* add docker-compose configuration for PostgreSQL and update HTML title ([7090d36](https://github.com/mizanmahi/aiusage/commit/7090d36806a36c447b7a6448346d03894845d41d))

## [0.4.0](https://github.com/mizanmahi/aiusage/compare/v0.3.0...v0.4.0) (2026-06-24)


### Features

* add GitHub Actions workflow for CLI publishing and set CLI version ([5548657](https://github.com/mizanmahi/aiusage/commit/554865785190d9b5862aa51972fddf7492e9e774))

## [0.3.0](https://github.com/mizanmahi/aiusage/compare/v0.2.0...v0.3.0) (2026-06-24)


### Features

* install goose command in Dockerfile and copy migrations ([ce83305](https://github.com/mizanmahi/aiusage/commit/ce833057a367692ea2aea9cf8960ec9d4b9ca6db))

## [0.2.0](https://github.com/mizanmahi/aiusage/compare/v0.1.0...v0.2.0) (2026-06-24)


### Features

* add admin dashboard ui ([e695aa2](https://github.com/mizanmahi/aiusage/commit/e695aa2e05b438af6b41c054e75f1855b7984079))
* add component composition and form rules documentation ([c1fedd4](https://github.com/mizanmahi/aiusage/commit/c1fedd470699ac7534c71cbfa6e646b50d30f24f))
* add cors and static ui serving ([c259bf1](https://github.com/mizanmahi/aiusage/commit/c259bf13f76cbacb951d8112edb10cf6a73fadf9))
* add dashboard analytics ([ccdc772](https://github.com/mizanmahi/aiusage/commit/ccdc772bcbeda6a0accec0df8f1d9f4dfe2b2816))
* add dashboard components for user analytics, including breakdown and summary tabs ([b558c95](https://github.com/mizanmahi/aiusage/commit/b558c95a517311471ca429a59b0d0d0e07623365))
* add Dockerfile and .dockerignore for containerization setup ([e7cd5ea](https://github.com/mizanmahi/aiusage/commit/e7cd5ea475590b46c033d020dfb955cc6b9c722f))
* add layered admin endpoints ([4ab6363](https://github.com/mizanmahi/aiusage/commit/4ab636399d60aaef63b482801c03b562827d3c0c))
* add litellm pricing resolver ([1384c53](https://github.com/mizanmahi/aiusage/commit/1384c53546806b527b9d577e55cf9f6ff8f95b8d))
* add push dry-run command ([32c9ecb](https://github.com/mizanmahi/aiusage/commit/32c9ecbdeff43c6219924c34bfa4c4b9ed8445c0))
* add real push command behavior ([e67ad0a](https://github.com/mizanmahi/aiusage/commit/e67ad0a9370f5ea960b1528e3b1c7a85b06a02b7))
* add release automation configuration with release-please ([b1868a9](https://github.com/mizanmahi/aiusage/commit/b1868a945487bb2661ae3e90b64410f8ba45d10c))
* add selected user analytics endpoints ([5f97458](https://github.com/mizanmahi/aiusage/commit/5f97458b9e507c5612579734c6eb48c5eaca2f68))
* add selected user analytics tabs ([b1e53fc](https://github.com/mizanmahi/aiusage/commit/b1e53fc3ef09fe8b31834c9d1fd747dd7bd2c602))
* add server routing and goose migrations ([e23cbba](https://github.com/mizanmahi/aiusage/commit/e23cbba8c3ed1a1ad38a8b0ec3b67bb2f9b24d5c))
* **cli:** add Claude session reader ([2f278b0](https://github.com/mizanmahi/aiusage/commit/2f278b0d7776706f81cb97faf4133786d9a29b2f))
* **cli:** add Codex session reader ([0193e9c](https://github.com/mizanmahi/aiusage/commit/0193e9cb68e3371e5b22406fb9740d225548409b))
* **cli:** add HTTP push client ([daacccc](https://github.com/mizanmahi/aiusage/commit/daaccccb1f631b5b9791ab43d871642d90ab6fd6))
* **cli:** add init and status commands ([eb9cf3a](https://github.com/mizanmahi/aiusage/commit/eb9cf3a483c277a997a0203cd84da0c7022f059f))
* **cli:** add local push state storage ([558f7c1](https://github.com/mizanmahi/aiusage/commit/558f7c14443cf99f9f0a8c542800e4accb0ea24b))
* **cli:** add local push state storage ([97a49d5](https://github.com/mizanmahi/aiusage/commit/97a49d504301767b1451664244e067a47d2f7cc2))
* dashboard analytics ([65d8a0d](https://github.com/mizanmahi/aiusage/commit/65d8a0dc9d63538da66f18f8ed5bba20109d0289))
* enhance dashboard components with new admin token dialog and developer selection, improve button and checkbox interactivity ([6f687bb](https://github.com/mizanmahi/aiusage/commit/6f687bbe27738275d3c22c1725a9870ffdfc71e3))
* implement authenticated ingest endpoint ([57ed841](https://github.com/mizanmahi/aiusage/commit/57ed84138426a2a8f1344f2f0f1790ede01d8c9d))
* **types:** add admin API response types ([eb64045](https://github.com/mizanmahi/aiusage/commit/eb64045e2ce0fa7f3a490973eea76298634f6971))
* update database migration commands and add .env to gitignore ([ae2d8ed](https://github.com/mizanmahi/aiusage/commit/ae2d8ed8ba48482322a98ebcb90f1408f7166867))
* update release-please action to include token configuration ([31e0555](https://github.com/mizanmahi/aiusage/commit/31e055550af197611ffb2014027603e29dcc3209))


### Bug Fixes

* preserve grouped breakdown row order ([72ba6fe](https://github.com/mizanmahi/aiusage/commit/72ba6fee26b18dc09403cbac8c7898924296d02e))
* read codex model from turn context ([3847e54](https://github.com/mizanmahi/aiusage/commit/3847e54d758861b700d80887461ddb52eab7a3e9))
* remove duplicate token configuration in release-please action ([4843492](https://github.com/mizanmahi/aiusage/commit/4843492b3db53063c04c668a229fafc9994a2192))
* simplify breakdown sorting ([25d67c2](https://github.com/mizanmahi/aiusage/commit/25d67c2e3dd395a70a0cf11010ef72bcd81c96c0))
