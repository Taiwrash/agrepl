import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';

// https://astro.build/config
export default defineConfig({
  site: 'https://taiwrash.github.io',
  base: '/agrepl',
  integrations: [mdx()],
});
