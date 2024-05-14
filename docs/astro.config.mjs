import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLinksValidator from 'starlight-links-validator';

import tailwind from '@astrojs/tailwind';

// https://astro.build/config
export default defineConfig({
  site: 'https://tanka.dev',
  integrations: [
    starlight({
      plugins: [starlightLinksValidator()],
      title: 'Grafana Tanka',
      description:
        'Grafana Tanka is the robust configuration utility for your Kubernetes cluster, powered by the Jsonnet language.',
      social: {
        github: 'https://github.com/grafana/tanka',
      },
      logo: {
        src: './img/logo.svg',
        alt: 'Grafana Tanka logo',
      },
      editLink: {
        baseUrl: 'https://github.com/grafana/tanka/edit/main/docs/',
      },
      components: {
        Hero: './src/components/Hero.astro',
      },
      customCss: ['./src/tailwind.css', '@fontsource-variable/inter'],
      sidebar: [
        {
          label: 'Installation',
          link: '/install',
        },
        {
          label: 'Tutorial',
          collapsed: true,
          autogenerate: {
            directory: 'tutorial',
          },
        },
        {
          label: 'Writing Jsonnet',
          collapsed: true,
          autogenerate: {
            directory: 'jsonnet',
          },
        },
        {
          label: 'Libraries',
          collapsed: true,
          autogenerate: {
            directory: 'libraries',
          },
        },
        {
          label: 'Advanced features',
          collapsed: true,
          items: [
            {
              label: 'Garbage collection',
              link: '/garbage-collection',
            },
            {
              label: 'Helm support',
              link: '/helm',
            },
            {
              label: 'Kustomize support',
              link: '/kustomize',
            },
            {
              label: 'Output filtering',
              link: '/output-filtering',
            },
            {
              label: 'Exporting as YAML',
              link: '/exporting',
            },
            {
              label: 'Inline environments',
              link: '/inline-environments',
            },
            {
              label: 'Server-Side Apply',
              link: '/server-side-apply',
            },
          ],
        },
        {
          label: 'References',
          collapsed: true,
          items: [
            {
              label: 'Configuration Reference',
              link: '/config',
            },
            {
              label: 'Directory structure',
              link: '/directory-structure',
            },
            {
              label: 'Environment variables',
              link: '/env-vars',
            },
            {
              label: 'Command-line completion',
              link: '/completion',
            },
            {
              label: 'Diff strategies',
              link: '/diff-strategy',
            },
            {
              label: 'Namespaces',
              link: '/namespaces',
            },
            {
              label: 'Formatting',
              link: '/formatting',
            },
          ],
        },
        {
          label: 'Frequently asked questions',
          link: '/faq',
        },
        {
          label: 'Known issues',
          link: '/known-issues',
        },
      ],
    }),
    tailwind({
      applyBaseStyles: false,
    }),
  ],
});
