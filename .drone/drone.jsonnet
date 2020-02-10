local golang = 'golang:1.13';

local volumes = [{ name: 'gopath', temp: {} }];
local mounts = [{ name: 'gopath', path: '/go' }];

local constraints = {
  onlyTagOrMaster: { trigger: {
    ref: [
      'refs/heads/master',
      'refs/heads/docker',
      'refs/tags/v*',
    ],
  } },
  onlyTags: { trigger: {
    event: ['tag'],
  } },
  always: {},
};

local go(name, commands) = {
  name: name,
  image: golang,
  volumes: mounts,
  commands: commands,
};

local make(target) = go(target, ['make ' + target]);

local pipeline(name) = {
  kind: 'pipeline',
  name: name,
  volumes: volumes,
  steps: [],
};

local docker(arch) = pipeline('docker-' + arch) {
  platform: {
    os: 'linux',
    arch: arch,
  },
  steps: [
    make('static'),
    {
      name: 'container',
      image: 'plugins/docker',
      settings: {
        repo: 'grafana/tanka',
        auto_tag: true,
        auto_tag_suffix: arch,
        username: { from_secret: 'docker_username' },
        password: { from_secret: 'docker_password' },
      },
    },
  ],
};

[
  pipeline('check') {
    steps: [
      go('download', ['go mod download']),
      make('lint') { depends_on: ['download'] },
      make('test') { depends_on: ['download'] },
    ],
  } + constraints.always,

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
          api_key: { from_secret: 'GITHUB_TOKEN' },
          files: 'dist/*',
        },
      },
    ],
  } + { depends_on: ['check'] } + constraints.onlyTags,

  docker('amd64') { depends_on: ['check'] } + constraints.onlyTagOrMaster,
  docker('arm') { depends_on: ['check'] } + constraints.onlyTagOrMaster,
  docker('arm64') { depends_on: ['check'] } + constraints.onlyTagOrMaster,

  pipeline('manifest') {
    steps: [{
      name: 'manifest',
      image: 'plugins/manifest',
      settings: {
        auto_tag: true,
        ignore_missing: true,
        spec: '.drone/docker-manifest.tmpl',
        username: { from_secret: 'docker_username' },
        password: { from_secret: 'docker_password' },
      },
    }],
  } + {
    depends_on: [
      'docker-amd64',
      'docker-arm',
      'docker-arm64',
    ],
  } + constraints.onlyTagOrMaster,
]
