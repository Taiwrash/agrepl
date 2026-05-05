import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';

import cloudflare from '@astrojs/cloudflare';

// https://astro.build/config
export default defineConfig({
  site: 'https://taiwrash.github.io',
  base: '/agrepl',
  integrations: [mdx()],
  adapter: cloudflare(),
});