load(
    'scripts/lib.star',
    'test_release_ver',
    'build_image',
    'publish_image',
    'pipeline',
    'lint_backend_step',
    'codespell_step',
    'shellcheck_step',
    'dashboard_schemas_check',
    'test_backend_step',
    'test_frontend_step',
    'build_backend_step',
    'build_frontend_step',
    'build_plugins_step',
    'package_step',
    'e2e_tests_server_step',
    'e2e_tests_step',
    'build_storybook_step',
    'copy_packages_for_docker_step',
    'build_docker_images_step',
    'postgres_integration_tests_step',
    'mysql_integration_tests_step',
    'get_windows_steps',
    'benchmark_ldap_step',
    'ldap_service',
    'frontend_metrics_step',
    'publish_storybook_step',
    'upload_packages_step',
    'notify_pipeline',
    'integration_test_services',
    'publish_packages_step',
)

def release_npm_packages_step(edition, ver_mode):
    if edition == 'enterprise':
        return None

    if ver_mode == 'release':
        commands = ['./scripts/build/release-packages.sh ${DRONE_TAG}']
    else:
        commands = []

    return {
        'name': 'release-npm-packages',
        'image': build_image,
        'depends_on': [
            # Has to run after publish-storybook since this step cleans the files publish-storybook depends on
            'publish-storybook',
        ],
        'environment': {
            'NPM_TOKEN': {
                'from_secret': 'npm_token',
            },
        },
        'commands': commands,
    }

def get_steps(edition, ver_mode):
    should_publish = ver_mode in ('release', 'test-release',)
    should_upload = should_publish or ver_mode in ('release-branch',)

    steps = [
        lint_backend_step(edition=edition),
        codespell_step(),
        shellcheck_step(),
        dashboard_schemas_check(),
        test_backend_step(),
        test_frontend_step(),
        build_backend_step(edition=edition, ver_mode=ver_mode),
        build_frontend_step(edition=edition, ver_mode=ver_mode),
        build_plugins_step(edition=edition, sign=True),
        package_step(edition=edition, ver_mode=ver_mode),
        e2e_tests_server_step(edition=edition),
        e2e_tests_step(),
        build_storybook_step(edition=edition, ver_mode=ver_mode),
        copy_packages_for_docker_step(),
        build_docker_images_step(edition=edition, ver_mode=ver_mode, publish=should_publish),
        build_docker_images_step(edition=edition, ver_mode=ver_mode, ubuntu=True, publish=should_publish),
        postgres_integration_tests_step(),
        mysql_integration_tests_step(),
    ]
    if should_upload:
        steps.append(upload_packages_step(edition=edition, ver_mode=ver_mode))
    if should_publish:
        steps.extend([
            publish_storybook_step(edition=edition, ver_mode=ver_mode),
            release_npm_packages_step(edition=edition, ver_mode=ver_mode),
        ])
    windows_steps = get_windows_steps(edition=edition, ver_mode=ver_mode)

    if edition == 'enterprise':
        build_tags = ['enterprise2']
        steps.extend([
            lint_backend_step(edition=edition, build_tags=build_tags),
            test_backend_step(build_tags=build_tags),
            build_backend_step(edition=edition, ver_mode=ver_mode, variants=['linux-x64'], build_tags=build_tags),
            package_step(edition=edition, ver_mode=ver_mode, variants=['linux-x64'], build_tags=build_tags),
            e2e_tests_server_step(edition=edition, build_tags=build_tags),
            e2e_tests_step(build_tags=build_tags),
        ])
        if should_upload:
            steps.append(upload_packages_step(edition=edition, ver_mode=ver_mode, build_tags=build_tags))

    return steps, windows_steps

def get_oss_pipelines(trigger, ver_mode):
    services = integration_test_services()
    steps, windows_steps = get_steps(edition='oss', ver_mode=ver_mode)
    return [
        pipeline(
            name='oss-build-{}'.format(ver_mode), edition='oss', trigger=trigger, services=services, steps=steps,
            ver_mode=ver_mode,
        ),
        pipeline(
            name='oss-windows-{}'.format(ver_mode), edition='oss', trigger=trigger, steps=windows_steps,
            platform='windows', depends_on=['oss-build-{}'.format(ver_mode)], ver_mode=ver_mode,
        ),
    ]

def get_enterprise_pipelines(trigger, ver_mode):
    services = integration_test_services()
    steps, windows_steps = get_steps(edition='enterprise', ver_mode=ver_mode)
    return [
        pipeline(
            name='enterprise-build-{}'.format(ver_mode), edition='enterprise', trigger=trigger, services=services,
            steps=steps, ver_mode=ver_mode,
        ),
        pipeline(
            name='enterprise-windows-{}'.format(ver_mode), edition='enterprise', trigger=trigger, steps=windows_steps,
            platform='windows', depends_on=['enterprise-build-{}'.format(ver_mode)], ver_mode=ver_mode,
        ),
    ]

def release_pipelines(ver_mode='release', trigger=None):
    services = integration_test_services()
    if not trigger:
        trigger = {
            'ref': ['refs/tags/v*',],
        }

    should_publish = ver_mode in ('release', 'test-release',)

    # The release pipelines include also enterprise ones, so both editions are built for a release.
    # We could also solve this by triggering a downstream build for the enterprise repo, but by including enterprise
    # in OSS release builds, we simplify the UX for the release engineer.
    oss_pipelines = get_oss_pipelines(ver_mode=ver_mode, trigger=trigger)
    enterprise_pipelines = get_enterprise_pipelines(ver_mode=ver_mode, trigger=trigger)

    pipelines = oss_pipelines + enterprise_pipelines
    if should_publish:
        publish_pipeline = pipeline(
            name='publish-{}'.format(ver_mode), trigger=trigger, edition='oss', steps=[
                publish_packages_step(edition='oss', ver_mode=ver_mode),
                publish_packages_step(edition='enterprise', ver_mode=ver_mode),
            ], depends_on=[p['name'] for p in oss_pipelines + enterprise_pipelines], install_deps=False,
            ver_mode=ver_mode,
        )
        pipelines.append(publish_pipeline)

    pipelines.append(notify_pipeline(
        name='notify-{}'.format(ver_mode), slack_channel='grafana-ci-notifications', trigger=trigger,
        depends_on=[p['name'] for p in pipelines],
    ))

    return pipelines

def test_release_pipelines():
    ver_mode = 'test-release'

    services = integration_test_services()
    trigger = {
        'event': ['custom',],
    }

    oss_pipelines = get_oss_pipelines(ver_mode=ver_mode, trigger=trigger)
    enterprise_pipelines = get_enterprise_pipelines(ver_mode=ver_mode, trigger=trigger)

    publish_cmd = './bin/grabpl publish-packages --edition {{}} --dry-run {}'.format(test_release_ver)

    publish_pipeline = pipeline(
        name='publish-{}'.format(ver_mode), trigger=trigger, edition='oss', steps=[
            publish_packages_step(edition='oss', ver_mode=ver_mode),
            publish_packages_step(edition='enterprise', ver_mode=ver_mode),
        ], depends_on=[p['name'] for p in oss_pipelines + enterprise_pipelines], install_deps=False,
        ver_mode=ver_mode,
    )

    pipelines = oss_pipelines + enterprise_pipelines + [publish_pipeline,]

    pipelines.append(notify_pipeline(
        name='notify-{}'.format(ver_mode), slack_channel='grafana-ci-notifications', trigger=trigger,
        depends_on=[p['name'] for p in pipelines],
    ))

    return pipelines
