local golang = 'golang:1.12';

local volumes = [{ name: 'gopath', temp: {} }];
local mounts = [{ name: 'gopath', path: '/go' }];

local constraints = {
  onlyTagOrMaster: {
    trigger: {
      branch: ['master', 'docker'],
      event: ['push', 'tag'],
    },
  },
  always: {}
};

local go(name, commands) = {
  name: name,
  image: golang,
  volumes: mounts,
  commands: commands,
};

local make(target) = go(target, ['make ' + target]);

local docker(arch) = {
  kind: 'pipeline',
  name: 'docker-' + arch,
  volumes: volumes,
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
        repo: 'shorez/tanka',
        auto_tag: true,
        auto_tag_suffix: arch,
        username: { from_secret: 'docker_username' },
        password: { from_secret: 'docker_password' },
      },
    },
  ],
};

local pipeline(name) = {
  kind: 'pipeline',
  name: name,
  volumes: volumes,
  steps: [],
};

local drone = [
  pipeline('check') {
    steps: [
      go('download', ['go mod download']),
      make('lint') { depends_on: ['download'] } + constraints.always,
      make('test') { depends_on: ['download'] } + constraints.always,
    ],
  },

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
        spec: '.docker-manifest.tmpl',
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
];

{
  drone: std.manifestYamlStream(drone),
}
