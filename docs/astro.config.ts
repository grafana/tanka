import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  site: 'https://tanka.dev',
  base: process.env.PATH_PREFIX ?? '/',
  trailingSlash: 'always',
  integrations: [
    starlight({
      title: 'Grafana Tanka',
      description:
        'Grafana Tanka is the robust configuration utility for your Kubernetes cluster, powered by the Jsonnet language.',
      social: [{
        icon: "github",
        label: "github",
        href: 'https://github.com/grafana/tanka',
      }],
      logo: {
        src: './img/logo.svg',
        alt: 'Grafana Tanka logo',
      },
      favicon: '/favicon.svg',
      editLink: {
        baseUrl: 'https://github.com/grafana/tanka/edit/main/docs/',
      },
      components: {
        Hero: './src/components/Hero.astro',
        TableOfContents: './src/components/TableOfContents.astro',
        MobileTableOfContents: './src/components/MobileTableOfContents.astro',
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
            {
              label: 'Telemetry',
              link: '/telemetry',
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
        {
          label: 'Project internals',
          collapsed: true,
          items: [
            {
              label: 'Releasing a new version',
              link: '/internal/releasing',
            },
          ],
        },
      ],
    }),
  ],
  vite: {
    plugins: [tailwindcss()],
  },
});
