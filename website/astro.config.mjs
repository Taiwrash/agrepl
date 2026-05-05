import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';



const isGitHub = process.env.GITHUB_ACTIONS === 'true';

export default defineConfig({
  site: isGitHub ? 'https://taiwrash.github.io' : 'https://agrepl.pages.dev',
  base: isGitHub ? '/agrepl' : '/',
  integrations: [mdx()]
});