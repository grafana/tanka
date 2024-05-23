# Tanka documentation & website

This folder contains the documentation for Tanka which is also published to
<https://tanka.dev>. Under the hood we are using [Starlight][Starlight] and
[Astro][Astro] for this, which allows for a lot of flexiblity regarding markup
and an easy preview of any changes.

## 🚀 Getting started

Before making your first changes, make sure you have a recent version of NodeJS
(\>= 20.x) installed. Once you have that, you can build a local previous of the
docs using the following commands:

```bash
# Install all dependencies
pnpm install

# Start preview server
pnpm run dev
```

This will prompt you an URL where you can see the preview.

You can find the source code for the documentation pages inside the
`src/content/docs` folder. These are either Markdown files or Markdown + JSX
(MDX) files.

If you want to make some changes, go, for instance, to
<http://localhost:4321/tutorial/overview> while having your development server
running. Now open an editor and make some changes to
`src/content/tutorial/overview.md`. The preview will reload upon saving that
file.

If you want to add images, place them in `src/assets` and embed them in your
Markdown files with relative links.

## 👀 Want to learn more?

Check out [Starlight’s docs](https://starlight.astro.build/) and the
[the Astro documentation](https://docs.astro.build)

[astro]: https://astro.build
[starlight]: https://starlight.astro.build
