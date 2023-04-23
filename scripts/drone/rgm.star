"""
rgm uses 'github.com/grafana/grafana-build' to build Grafana on the following events:
* A merge to main
* A tag that begins with a 'v'
"""

load(
    "scripts/drone/utils/utils.star",
    "pipeline",
)

load("scripts/drone/vault.star",
  "from_secret",
  "vault_secret",
  "rgm_gcp_key_base64",
  "rgm_destination",
  "rgm_github_token",
)

rgm_env_secrets = {
  "GCP_KEY_BASE64": from_secret(rgm_gcp_key_base64),
  "DESTINATION": from_secret(rgm_destination),
  "GITHUB_TOKEN": from_secret(rgm_github_token),
}

def rgm_build(distros=["linux/amd64"], script="drone_publish_main.sh"):
  clone_step = {
    "name": "clone-rgm",
    "image": "alpine/git",
    "commands": [
      "git clone https://github.com/grafana/grafana-build.git rgm",
    ],
    "failure": "ignore",
  }

  distroStr = ",".join(distros)
  env = {
    "DISTROS": distroStr,
  }

  for key in rgm_env_secrets:
    env[key] = rgm_env_secrets[key]

  rgm_build_step = {
    "name": "rgm-build",
    "image": "golang:1.20.3-alpine",
    "commands": [
      # the docker program is a requirement for running dagger programs
      "apk update && apk add docker",
      "export GRAFANA_DIR=$$(pwd)",
      "cd rgm && ./scripts/{}".format(script)
    ],
    "environment": env,
    # The docker socket is a requirement for running dagger programs
    # In the future we should find a way to use dagger without mounting the docker socket.
    "volumes": [{"name": "docker", "path": "/var/run/docker.sock"}],
    "failure": "ignore",
  }

  return [
    clone_step,
    rgm_build_step,
  ]


def rgm_main():
  trigger = {
      "event": [
          "push",
      ],
      "branch": "add-rgm-to-drone",
  }

  return pipeline(
    name="rgm-main-build",
    edition="all",
    trigger=trigger,
    steps=rgm_build(script="drone_publish_tag.sh"),
  )

def rgm_tag():
  trigger = {
      "event": {
          "exclude": [
              "promote",
          ],
      },
      "ref": {
          "include": [
              "refs/tags/v*",
          ],
          "exclude": [
              "refs/tags/*-cloud*",
          ],
      },
  }

  return pipeline(
    name="rgm-tag-build",
    edition="all",
    trigger=trigger,
    steps=rgm_build(),
  )

def rgm():
  return [
    rgm_main(),
    rgm_tag(),
  ]
