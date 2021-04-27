local secret(name, vault_path, key) = {
  kind: 'secret',
  name: name,
  get: {
    path: vault_path,
    name: key,
  },
};

{
  dockerhub_username: 'dockerhub_username',
  dockerhub_password: 'dockerhub_password',
  grafanabot_public_account_token: 'grafanabot_pat',

  secrets: [
    secret($.grafanabot_public_account_token, 'infra/data/ci/github/grafanabot', 'pat'),
    secret($.dockerhub_username, 'infra/data/ci/docker_hub', 'username'),
    secret($.dockerhub_password, 'infra/data/ci/docker_hub', 'password'),
  ],
}
