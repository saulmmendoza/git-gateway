// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLlmsTxt from 'starlight-llms-txt';
import pagefind from 'astro-pagefind';
import compress from '@playform/compress';
import swup from '@swup/astro';

// https://astro.build/config
export default defineConfig({
	site: 'https://netlify.github.io/git-gateway',
	integrations: [
		starlight({
			title: 'Git Gateway Docs',
			customCss: [],
			plugins: [
				starlightLlmsTxt(),
			],
			sidebar: [
				{
					label: 'Quickstart',
					items: [{ autogenerate: { directory: 'quickstart' } }],
				},
				{
					label: 'Overview',
					items: [{ autogenerate: { directory: 'overview' } }],
				},
				{
					label: 'Guides',
					items: [{ autogenerate: { directory: 'guides' } }],
				},
				{
					label: 'Integrations',
					items: [{ autogenerate: { directory: 'integrations' } }],
				},
				{
					label: 'API Reference',
					items: [{ autogenerate: { directory: 'reference' } }],
				},
			],
		}),
		pagefind(),
		swup(),
		compress(),
	],
});
