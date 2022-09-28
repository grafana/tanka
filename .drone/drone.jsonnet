local vault = import 'vault.libsonnet';

local golang = 'golang:1.19';

local volumes = [{ name: 'gopath', temp: {} }];
local mounts = [{ name: 'gopath', path: '/go' }];

local constraints = {
  local withRef(ref) = {
    trigger+: {
      ref+: [ref],
    },
  },

  tags: withRef('refs/tags/v*'),
  mainPush: withRef('refs/heads/main'),
  pullRequest: withRef('refs/pull/*/head'),
};

local go(name, commands) = {
  name: name,
  image: golang,
  volumes: mounts,
  commands: commands,
};

local make(target) = go(target, [
  // Only download it once, then for every step, copy it to the right place.
  'if [ ! -f linux-amd64/helm ]; then',
  '  wget -q https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz',
  '  tar -zxvf helm-v3.9.0-linux-amd64.tar.gz',
  '  rm -f helm-v3.9.0-linux-amd64.tar.gz',
  'fi',
  'cp linux-amd64/helm /usr/local/bin/helm',
  'make ' + target,
]);

local pipeline(name) = {
  kind: 'pipeline',
  name: name,
  volumes: volumes,
  steps: [],
};

local docker(name, arch, tags=null, depends_on=[]) =
  pipeline(name) {
    platform: {
      os: 'linux',
      arch: arch,
    },
    steps: [
      go('fetch-tags', ['git fetch origin --tags']),
      make('static'),
      {
        name: 'container',
        image: 'plugins/docker',
        settings: {
          repo: 'grafana/tanka',

          username: { from_secret: vault.dockerhub_username },
          password: { from_secret: vault.dockerhub_password },
        } + if tags == null then {
          auto_tag: true,
          auto_tag_suffix: arch,
        } else {
          tags: tags,
        },
      },
    ],
    depends_on: depends_on,
  };

[
  pipeline('check') {
    steps: [
      go('download', ['go mod download']),
      make('lint'),
      make('test'),
      make('cross') { name: 'build' },
    ],
  } + constraints.pullRequest + constraints.mainPush,

  pipeline('release') {
    steps: [
      go('fetch-tags', ['git fetch origin --tags']),
      make('cross'),
      {
        name: 'publish',
        image: 'plugins/github-release',
        settings: {
          title: '${DRONE_TAG}',
          note: importstr 'release-note.md',
          api_key: { from_secret: vault.grafanabot_public_account_token },
          files: 'dist/*',
          draft: true,
        },
      },
    ],
  } + { depends_on: ['check'] } + constraints.tags,

  docker('docker-main-commit', 'amd64', tags=['${DRONE_COMMIT}'], depends_on=['check']) + constraints.mainPush,
  docker('docker-amd64', 'amd64', depends_on=['check']) + constraints.tags + constraints.mainPush,
  docker('docker-arm64', 'arm64', depends_on=['check']) + constraints.tags + constraints.mainPush,

  pipeline('manifest') {
    steps: [{
      name: 'manifest',
      image: 'plugins/manifest',
      settings: {
        auto_tag: true,
        ignore_missing: true,
        spec: '.drone/docker-manifest.tmpl',
        username: { from_secret: vault.dockerhub_username },
        password: { from_secret: vault.dockerhub_password },
      },
    }],
  } + {
    depends_on: [
      'docker-amd64',
      'docker-arm64',
    ],
  } + constraints.tags + constraints.mainPush,
] + vault.secrets
